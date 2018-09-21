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
)

var authMiddleware = auth.GetMiddleware()

func setupDB() *gorm.DB {
	db := database.Init()
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&topic.Topic{})
	return db
}

func setupRouter() *gin.Engine {
	//	route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//	public endpoint
	//	TODO: FB user study
	r.POST("/user/register", user.UserRegister)

	//	auth
	r.POST("/login", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/ping", user.PingHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}
	authRoot := r.Group("/")
	authRoot.Use(authMiddleware.MiddlewareFunc())
	{
		authRoot.GET("/user", user.UserInfo)
	}

	r.POST("/user/login", user.UserLogin)

	//	TODO: dummy code here for now
	r.GET("/users", user.DummyUsersList)
	r.GET("/users/:uid", user.DummyUsersInfo)

	//	user
	r.GET("/topics", topic.List)
	r.GET("/topics/:id", topic.Detail)
	r.POST("/topics", topic.Create)

	return r
}

func main() {
	db := setupDB()
	defer db.Close()

	r := setupRouter()
	r.Run(":8080")
}