package user

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	//"github.com/appleboy/gin-jwt"
	"hindsight/database"
	apiError "hindsight/error"
	"hindsight/facebook"
)

//	Use this function if we are going to implement our own auth middleware
/*
func FacebookConnect(c *gin.Context) {
	code, response, user := FacebookAuthenticate(c)
	if user != nil {
		c.JSON(code, user.Response())
	} else {
		c.JSON(code, response)
	}
}
*/

func FacebookAuthenticate(c *gin.Context) (int, gin.H, *User) {
	var request facebook.ConnectRequest
	if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		code, response := apiError.Bad(apiError.DomainFacebookConnect, apiError.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	if request.AccessToken == "" {
		code, response := apiError.Bad(apiError.DomainFacebookConnect, apiError.ReasonInvalidJSON, "Invalid access_token")
		return code, response, nil
	}

	var fbUser facebook.User
	var err error
	//log.Println(request)
	if fbUser, err = facebook.Connect(request.AccessToken); err != nil {
		code, response := apiError.Bad(apiError.DomainFacebookConnect, apiError.ReasonInvalidJSON, err.Error())
		return code, response, nil
	}

	//c.JSON(http.StatusOK, fbUser)
	log.Println(fbUser)

	var user User
	db := database.GetDB()

	if notFound := db.Where(User{FacebookUserID: fbUser.ID}).First(&user).RecordNotFound(); !notFound {
		// user already exists
		log.Println(user)
		return http.StatusOK, user.Response(), &user
	}

	// create new user with a randomized unique username
	user = User{Username: fbUser.UniqueUsername(), FacebookUserID: fbUser.ID}
	if err := db.Create(&user).Error; err != nil {
		code, response := apiError.Bad(apiError.DomainUserRegister, apiError.ReasonDuplicatedEntry, err.Error())
		return code, response, nil
	}

	return http.StatusOK, user.Response(), &user
}
