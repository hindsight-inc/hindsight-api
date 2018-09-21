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

func Authenticate(c *gin.Context) (int, gin.H, *User) {
	var json User
	if err := c.ShouldBindJSON(&json); err != nil {
		code, response := error.Bad(error.DomainUserLogin, error.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	var user User
	db := database.GetDB()
	db.Where(User{Username: json.Username}).First(&user)
	if user.ID == 0 {
		code, response := error.Unauthorized(error.DomainUserLogin, error.ReasonNonexistentEntry, "User not found")
		return code, response, nil
	}
	if user.Password != json.Password {
		code, response := error.Unauthorized(error.DomainUserLogin, error.ReasonMismatchedEntry, "Wrong password")
		return code, response, nil
	}
	return http.StatusOK, user.Response(), &user
}

func Authorize(c *gin.Context, username string) bool {
	var user User
	db := database.GetDB()
	db.Where(User{Username: username}).First(&user)
	if user.ID == 0 {
		return false
	}
	return true
}

func UserLogin(c *gin.Context) {
	code, response, _ := Authenticate(c)
	c.JSON(code, response)
}

/*
http://localhost:8080/user
*/
func UserInfo(c *gin.Context) {
	user := User{Username: "username002"}
	c.JSON(200, user.Response())
}