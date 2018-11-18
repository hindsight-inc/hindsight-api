package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"hindsight/database"
	"hindsight/herror"
)

func UserRegister(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(herror.Bad(herror.DomainUserRegister, herror.ReasonInvalidJSON, err.Error()))
		return
	}

	//	TODO: validate username

	db := database.Shared()
	if notFound := db.Where(User{Username: user.Username}).First(&user).RecordNotFound(); !notFound {
		c.JSON(herror.Bad(herror.DomainUserRegister, herror.ReasonDuplicatedEntry, "User already exists"))
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainUserUpdate, herror.ReasonDatabaseError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, user.Response())
}

func UserUpdate(c *gin.Context) {
	var request UpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(herror.Bad(herror.DomainUserUpdate, herror.ReasonInvalidJSON, err.Error()))
		return
	}

	//	TODO: validate username

	u := Current(c)
	if u == nil {
		c.JSON(herror.Bad(herror.DomainUserUpdate, herror.ReasonNonexistentEntry, "Current user not found"))
		return
	}

	db := database.Shared()
	u.Username = request.Username
	//	TODO: why cannot use Update?
	if err := db.Save(&u).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainUserUpdate, herror.ReasonDatabaseError, err.Error()))
		return
	}

	//	TODO: token will be invalidated at this state, should return new token
	c.JSON(http.StatusOK, u.Response())
}

func Authenticate(c *gin.Context) (int, gin.H, *User) {
	var json User
	if err := c.ShouldBindBodyWith(&json, binding.JSON); err != nil {
		code, response := herror.Bad(herror.DomainUserLogin, herror.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	var user User
	db := database.Shared()
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
	db := database.Shared()
	return !db.Where(User{Username: username}).First(&user).RecordNotFound()
}

func UserDetail(c *gin.Context) {
	u := Current(c)
	c.JSON(u.DetailResponse())
}