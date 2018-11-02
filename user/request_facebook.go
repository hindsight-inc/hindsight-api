package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	//"github.com/appleboy/gin-jwt"
	//"hindsight/database"
	"hindsight/error"
	"hindsight/facebook"
)

// POST {{host}}/user/facebook/connect
// access_token: FB_TOKEN
func FacebookConnect(c *gin.Context) {
	var request facebook.ConnectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(error.Bad(error.DomainFacebookConnect, error.ReasonInvalidJSON, err.Error()))
		return
	}

	if request.AccessToken == "" {
		c.JSON(error.Bad(error.DomainFacebookConnect, error.ReasonInvalidJSON, "Invalid access_token"))
		return
	}

	//log.Println(request)
	if fbUser, err := facebook.Connect(request.AccessToken); err == nil {
		c.JSON(http.StatusOK, fbUser)
	} else {
		c.JSON(error.Bad(error.DomainFacebookConnect, error.ReasonInvalidJSON, err.Error()))
	}

	/*
	db := database.GetDB()
	db.Where(User{Username: user.Username}).First(&user)
	if notFound := db.Where(User{Username: user.Username}).First(&user).RecordNotFound(); !notFound {
		c.JSON(error.Bad(error.DomainUserRegister, error.ReasonDuplicatedEntry, "User already exists"))
		return
	}

	db.Create(&user)
	*/
	//c.JSON(http.StatusOK, model.Response())
}
