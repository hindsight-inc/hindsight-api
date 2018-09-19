package user

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
	"hindsight/error"
)

func (self *User) Response() gin.H {
	return gin.H{"id": self.ID, "username": self.Username}
}

/*
curl -v -X POST \
  http://localhost:8080/user/register \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserRegister(c *gin.Context) {
	db := database.GetDB()
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Where(User{Username: user.Username}).First(&user)
	if user.ID == 0 {
		db.Create(&user)
		c.JSON(http.StatusOK, user.Response())
	} else {
		//c.JSON(http.StatusBadRequest, error.New("User already exists"))
		//c.JSON(error.New("user.register.existing"))
		c.JSON(error.H(error.UserRegisterExisting))
	}
}

/*
curl -v -X POST \
  http://localhost:8080/user/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserLogin(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username != "username001" || user.Password != "password001" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UserInfo(c *gin.Context) {
	user := User{Username: "username002"}
	c.JSON(200, user.Response())
}