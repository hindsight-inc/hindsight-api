package user

import (
	//"log"
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

/* Response */

func (self *User) Response() gin.H {
	return gin.H{"id": self.ID, "username": self.Username}
}

func (self *User) DetailResponse() gin.H {
	// TODO: this is just an example, it will be changed base on actual business logic.
	db := database.GetDB()
	if err := db.Model(self).Related(&self.FacebookUser, "FacebookUser").Error; err != nil {
		return gin.H{"id": self.ID, "username": self.Username, "facebook_name": "N/A"}
	}
	return gin.H{"id": self.ID, "username": self.Username, "facebook_name": self.FacebookUser.Name}
}

/* Auth */

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