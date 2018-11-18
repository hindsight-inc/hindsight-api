package user

import (
	//"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	//"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/herror"
	"hindsight/facebook"
)

func FacebookAuthenticate(c *gin.Context) (int, gin.H, *User) {
	var request facebook.ConnectRequest
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		code, response := herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	if request.AccessToken == "" {
		code, response := herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, "Invalid access_token")
		return code, response, nil
	}

	var fbUser facebook.User
	var err error
	if fbUser, err = facebook.Connect(request.AccessToken); err != nil {
		code, response := herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	var user User
	db := database.Shared()

	if notFound := db.Where(User{FacebookUserID: fbUser.ID}).First(&user).RecordNotFound(); !notFound {
		//log.Println("User already exists: " + user)
		return http.StatusOK, user.Response(), &user
	}

	// Create new user with a randomized unique username
	user = User{Username: fbUser.UniqueUsername(), FacebookUserID: fbUser.ID}
	if err := db.Create(&user).Error; err != nil {
		code, response := herror.Bad(herror.DomainUserRegister, herror.ReasonDuplicatedEntry, err.Error())
		return code, response, nil
	}

	return http.StatusOK, user.Response(), &user
}