package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

///	DB

type Database struct {
	*gorm.DB
}

var DB *gorm.DB

// Opening a database and save the reference to `Database` struct.
func Init() *gorm.DB {
	//db, err := gorm.Open("sqlite3", "./../gorm.db")
	db, err := gorm.Open("mysql", "golang:password@/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println("db err: ", err)
	}
	db.DB().SetMaxIdleConns(10)
	//db.LogMode(true)
	DB = db
	return DB
}

//var db gorm.DB

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (self *User) Response() gin.H {
	return gin.H{"username": self.Username}
}

///	main

/*
curl -v -X POST \
  http://localhost:8080/user/register \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserRegister(context *gin.Context) {
	var json User
	if err := context.ShouldBindJSON(&json); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(DB)
	log.Println(json)
	DB.Create(&User{Username: json.Username, Password: json.Password})
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

/*
curl -v -X POST \
  http://localhost:8080/user/login \
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
	/*
	var err error
	db, err := gorm.Open("mysql", "golang:password@/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	*/
	//db.DB().SetMaxIdleConns(10)
	//db.LogMode(true)
	//defer db.Close()
	//log.Println(db)
	//db.Create(&User{Username: "admin", Password: "password"})

	DB = Init()
	DB.AutoMigrate(&User{})
	defer DB.Close()

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/user/register", UserRegister)
	router.POST("/user/login", UserLogin)
	router.GET("/user", UserInfo)

	router.Run()
}