package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/user"
	"hindsight/topic"
	"hindsight/auth"
)

var authMiddleware = auth.GetMiddleware()

func setupRouter() *gin.Engine {
	//	route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//	auth
	//authMiddleware := auth.GetMiddleware()
	r.POST("/login", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", user.HelloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	//	TODO: FB user study
	r.POST("/user/register", user.UserRegister)
	r.POST("/user/login", user.UserLogin)
	r.GET("/user", user.UserInfo)

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
	//	database
	db := database.Init()
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&topic.Topic{})
	defer db.Close()

	r := setupRouter()
	r.Run(":8080")
}