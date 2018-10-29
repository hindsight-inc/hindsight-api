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
	//"hindsight/facebook"
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
	r.POST("/user/login", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/token")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/ping", user.PingHandler)
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

	//r.POST("/user/login", user.UserLogin)
	//r.GET("/users", user.DummyUsersList)
	//r.GET("/users/:uid", user.DummyUsersInfo)

	return r
}

func main() {
	if _, err := config.Init(); err != nil {
		panic(err)
	}

	db := setupDB()
	defer db.Close()
	
	//facebook.Test()
	r := setupRouter()
	r.Run(":8080")
}
