package user

import (
	"log"
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
	/*
	token := c.DefaultQuery("access_token", "")
	if token == nil {
		c.JSON(error.Bad(error.DomainUserRegister, error.ReasonDuplicatedEntry, "User already exists"))
		return
	}
	*/

	// We could also use object instead
	var model facebook.ConnectModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(error.Bad(error.DomainFacebookConnect, error.ReasonInvalidJSON, err.Error()))
		return
	}

	if model.AccessToken == "" {
		c.JSON(error.Bad(error.DomainFacebookConnect, error.ReasonInvalidJSON, "Invalid access_token"))
		return
	}

	log.Println(model)

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
	c.JSON(http.StatusOK, nil)
}
