package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/herror"
)

func UserRegister(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(herror.Bad(herror.DomainUserRegister, herror.ReasonInvalidJSON, err.Error()))
		return
	}

	db := database.GetDB()
	if notFound := db.Where(User{Username: user.Username}).First(&user).RecordNotFound(); !notFound {
		c.JSON(herror.Bad(herror.DomainUserRegister, herror.ReasonDuplicatedEntry, "User already exists"))
		return
	}

	db.Create(&user)
	c.JSON(http.StatusOK, user.Response())
}

func Authenticate(c *gin.Context) (int, gin.H, *User) {
	var json User
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		code, response := herror.Bad(herror.DomainUserLogin, herror.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	var user User
	db := database.GetDB()
	db.Where(User{Username: json.Username}).First(&user)
	if user.ID == 0 {
		code, response := herror.Unauthorized(herror.DomainUserLogin, herror.ReasonNonexistentEntry, "User not found")
		return code, response, nil
	}
	if user.Password != json.Password {
		code, response := herror.Unauthorized(herror.DomainUserLogin, herror.ReasonMismatchedEntry, "Wrong password")
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

//	Use this function if we are going to implement our own auth middleware
/*
func UserLogin(c *gin.Context) {
	code, response, _ := Authenticate(c)
	c.JSON(code, response)
}
*/

func UserInfo(c *gin.Context) {
	var user User
	claim := jwt.ExtractClaims(c)[IdentityKey]
	if claim == nil {
		c.JSON(herror.Unauthorized(herror.DomainUserInfo, herror.ReasonEmptyEntry, "Missing authorization info"))
		return
	}
	username := claim.(string)
	//user, _ := c.Get(IdentityKey)
	//username := user.(*User).Username
	db := database.GetDB()
	db.Where(User{Username: username}).First(&user)
	if user.ID == 0 {
		//	Shouldn't reach here unless user has been deleted but active token is not
		c.JSON(herror.Bad(herror.DomainUserInfo, herror.ReasonNonexistentEntry, "User not found"))
		return
	}
	c.JSON(user.DetailResponse())
}