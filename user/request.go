package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
)

func (self *User) Response() gin.H {
	return gin.H{"username": self.Username}
}

/*
curl -v -X POST \
  http://localhost:8080/user/register \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserRegister(context *gin.Context) {
	db := database.GetDB()
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Where(User{Username: user.Username}).First(&user)
	if user.ID == 0 {
		db.Create(&User{Username: user.Username, Password: user.Password})
		context.JSON(http.StatusOK, user.Response())
	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
	}
}

/*
curl -v -X POST \
  http://localhost:8080/user/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username != "username001" || user.Password != "password001" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UserInfo(context *gin.Context) {
	user := User{Username: "username002"}
	context.JSON(200, user.Response())
}
