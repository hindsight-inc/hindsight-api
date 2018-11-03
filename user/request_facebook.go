package user

import (
	"log"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
	//"github.com/appleboy/gin-jwt"
	"hindsight/database"
	herror "hindsight/error"
	"hindsight/facebook"
)

// POST {{host}}/user/facebook/connect
// access_token: FB_TOKEN
func FacebookConnect(c *gin.Context) {
	var request facebook.ConnectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, err.Error()))
		return
	}

	if request.AccessToken == "" {
		c.JSON(herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, "Invalid access_token"))
		return
	}

	var fbUser facebook.User
	var err error
	//log.Println(request)
	if fbUser, err = facebook.Connect(request.AccessToken); err != nil {
		c.JSON(herror.Bad(herror.DomainFacebookConnect, herror.ReasonInvalidJSON, err.Error()))
		return
	}

	//c.JSON(http.StatusOK, fbUser)
	log.Println(fbUser)

	var user User
	db := database.GetDB()

	if notFound := db.Where(User{FacebookUserID: fbUser.FacebookID}).First(&user).RecordNotFound(); !notFound {
		// user already exists
		c.JSON(http.StatusOK, user)
		return
	}

	// create new user
	//user = User{Username: "fb" + strconv.FormatInt(fbUser.FacebookID, 10), FacebookUserID: fbUser.id}
	user = User{Username: "fb" + strconv.FormatInt(fbUser.FacebookID, 10)}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(herror.Bad(herror.DomainUserRegister, herror.ReasonDuplicatedEntry, err.Error()))
		return
	}

	c.JSON(http.StatusOK, user)

	/*
	db.Create(&user)
	*/
	//c.JSON(http.StatusOK, model.Response())
}
