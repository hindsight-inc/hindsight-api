package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (self *User) Response() gin.H {
	return gin.H{"username": self.Username}
}

/*
curl -v -X POST \
  http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserLogin(context *gin.Context) {
	var json User
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.Username != "username001" || json.Password != "password001" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UserInfo(context *gin.Context) {
	user := User{Username: "username002"}
	context.JSON(200, user.Response())
}

func main() {
	db, err := gorm.Open("mysql", "golang:password@/golang?charset=utf8&parseTime=True&loc=Local")
	log.Println(err)
	defer db.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/user/register", UserInfo)
	router.POST("/user/login", UserLogin)
	router.GET("/user", UserInfo)

	router.Run()
}