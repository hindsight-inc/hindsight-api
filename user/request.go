package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"hindsight/database"
	"hindsight/error"
)

func (self *User) Response() gin.H {
	return gin.H{"id": self.ID, "username": self.Username}
}

func UserRegister(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(error.Bad(error.DomainUserRegister, error.ReasonInvalidJSON, err.Error()))
		return
	}
	db := database.GetDB()
	db.Where(User{Username: user.Username}).First(&user)
	if user.ID == 0 {
		db.Create(&user)
		c.JSON(http.StatusOK, user.Response())
	} else {
		c.JSON(error.Bad(error.DomainUserRegister, error.ReasonDuplicatedEntry, "User already exists"))
	}
}

func UserLogin(c *gin.Context) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(error.Bad(error.DomainUserLogin, error.ReasonInvalidJSON, err.Error()))
		return
	}
	var user User
	db := database.GetDB()
	db.Where(User{Username: json.Username}).First(&user)
	if user.ID == 0 {
		c.JSON(error.Unauthorized(error.DomainUserLogin, error.ReasonNonexistentEntry, "User not found"))
		return
	}
	if user.Password != json.Password {
		c.JSON(error.Unauthorized(error.DomainUserLogin, error.ReasonMismatchedEntry, "Wrong password"))
		return
	}
	c.JSON(http.StatusOK, user.Response())
}

func UserInfo(c *gin.Context) {
	user := User{Username: "username002"}
	c.JSON(200, user.Response())
}