package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
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
	if notFound := db.Where(User{Username: user.Username}).First(&user).RecordNotFound(); !notFound {
		c.JSON(error.Bad(error.DomainUserRegister, error.ReasonDuplicatedEntry, "User already exists"))
		return
	}

	db.Create(&user)
	c.JSON(http.StatusOK, user.Response())
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
	//	Current behavior: user has access to everything as long as s/he's logged in
	var user User
	db := database.GetDB()
	return !db.Where(User{Username: username}).First(&user).RecordNotFound()
}

func UserLogin(c *gin.Context) {
	code, response, _ := Authenticate(c)
	c.JSON(code, response)
}

/*
curl -v GET \
  http://localhost:8080/user \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'
*/
func UserInfo(c *gin.Context) {
	var user User
	claim := jwt.ExtractClaims(c)[IdentityKey]
	if claim == nil {
		c.JSON(error.Unauthorized(error.DomainUserInfo, error.ReasonEmptyEntry, "Missing authorization info"))
		return
	}
	username := claim.(string)
	//user, _ := c.Get(IdentityKey)
	//username := user.(*User).Username
	db := database.GetDB()
	db.Where(User{Username: username}).First(&user)
	if user.ID == 0 {
		//	Shouldn't reach here unless user has been deleted but active token is not
		c.JSON(error.Bad(error.DomainUserInfo, error.ReasonNonexistentEntry, "User not found"))
		return
	}
	c.JSON(http.StatusOK, user)
}