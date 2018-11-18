package main

import (
	//"log"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"hindsight/auth"
	"hindsight/config"
	"hindsight/database"
	"hindsight/facebook"
	"hindsight/file"
	"hindsight/topic"
	"hindsight/user"
)

var authMiddleware *jwt.GinJWTMiddleware

func setupAuth() *jwt.GinJWTMiddleware {
	authMiddleware = auth.GetMiddleware()
	return authMiddleware
}

func setupDB() *gorm.DB {
	db := database.Init()
	db.AutoMigrate(&user.User{})
	db.AutoMigrate(&facebook.User{})

	db.AutoMigrate(&topic.Topic{})
	db.AutoMigrate(&topic.Opinion{})
	db.AutoMigrate(&topic.Vote{})

	db.AutoMigrate(&file.Image{})
	return db
}

func internalTest(c *gin.Context) {
	//facebook.Connect("access_token")
}

func setupRouter() *gin.Engine {
	//	route
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.Static("/image", "./public/upload/image")
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
		//claims := jwt.ExtractClaims(c)
		//log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"message": "Invalid API"})
	})

	auth := r.Group("/token")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/ping", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			u, _ := c.Get(user.IdentityKey)
			c.JSON(200, gin.H{
				"message":  "pong",
				"claim_id": claims[user.IdentityKey],
				"username": u.(*user.User).Username, //	TODO: replace Username with ID, as Username will be nullable
				"id":       u.(*user.User).ID,       //	TODO: why cannot obtain u.ID?
			})
		})
		auth.GET("/refresh", authMiddleware.RefreshHandler)
	}

	authRoot := r.Group("/")
	authRoot.Use(authMiddleware.MiddlewareFunc())
	{
		authRoot.GET("/user", user.UserInfo)
		authRoot.PATCH("/user", user.UserUpdate)

		authRoot.GET("/topics", topic.List)
		authRoot.GET("/topics/:id", topic.Detail)
		authRoot.POST("/topics", topic.Create)
		authRoot.POST("/topics/:id/vote/:oid", topic.VoteOpinion)

		authRoot.POST("/file/image", file.UploadImage)
	}

	return r
}

func setupConfig() {
	provider := new(config.ViperProvider)
	if _, err := config.Init(provider); err != nil {
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