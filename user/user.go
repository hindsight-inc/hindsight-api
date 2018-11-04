package user

import (
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/facebook"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	// facebook
	FacebookUser	facebook.User
	FacebookUserID	uint
}

//	Don't use old token after changing this, see: https://github.com/appleboy/gin-jwt/issues/170
const IdentityKey = "user.id"

func Current(c *gin.Context) *User {
	// TODO: performance issue - when topic.Create is called, firstly Authenticate checks if user is valid, then Current is called to obtain the user. Can we combine these 2 queries by getting user.ID from JWT?
	var u User
	claim := jwt.ExtractClaims(c)[IdentityKey]
	if claim == nil {
		return nil
	}
	username := claim.(string)
	db := database.GetDB()
	db.Where(User{Username: username}).First(&u)
	if u.ID == 0 {
		return nil
	}
	return &u
}