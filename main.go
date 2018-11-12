package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"github.com/jinzhu/gorm"

	"hindsight/database"
	"hindsight/user"
	"hindsight/topic"
	"hindsight/auth"
	"hindsight/config"
	"hindsight/facebook"
)

var authMiddleware *jwt.GinJWTMiddleware

func setupAuth() *jwt.GinJWTMiddleware {
	authMiddleware = auth.GetMiddleware()
	return authMiddleware
}

func setupDB() *gorm.DB {
	db := database.Init()
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&topic.Topic{})
	db.AutoMigrate(&topic.Opinion{})
	db.AutoMigrate(&topic.Vote{})
	db.AutoMigrate(&facebook.User{})
	return db
}

func internalTest(c *gin.Context) {
	//facebook.Connect("access_token")
}

func setupRouter() *gin.Engine {
	//	route
	r := gin.Default()
	r.GET("/test", internalTest)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//	public endpoint
	r.POST("/user/register", user.UserRegister)

	//	auth
	r.POST("/user/login", authMiddleware.LoginHandler)
	r.POST("/user/connect", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/token")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/ping", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			u, _ := c.Get(user.IdentityKey)
			c.JSON(200, gin.H{
				"message": "pong",
				"claim_id": claims[user.IdentityKey],
				"username": u.(*user.User).Username,	//	TODO: replace Username with ID, as Username will be nullable
			})
		})
		auth.GET("/refresh", authMiddleware.RefreshHandler)
	}

	authRoot := r.Group("/")
	authRoot.Use(authMiddleware.MiddlewareFunc())
	{
		authRoot.GET("/user", user.UserInfo)

		authRoot.GET("/topics", topic.List)
		authRoot.GET("/topics/:id", topic.Detail)
		authRoot.POST("/topics", topic.Create)
	}

	return r
}

func setupConfig() {
	if _, err := config.Init(); err != nil {
		panic(err)
	}
}

func setupFacebook() {
	if err := facebook.Init(); err != nil {
		panic(err)
	}
}

func main() {
	setupConfig()

	db := setupDB()
	defer db.Close()

	setupFacebook()
	setupAuth()

	r := setupRouter()
	r.Run(":8080")
}