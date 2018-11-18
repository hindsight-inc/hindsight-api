package user

import (
	//"log"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/facebook"
	//"hindsight/herror"
)

type User struct {
	gorm.Model
	Username string `binding:"required"`
	Password string
	// facebook
	FacebookUser	facebook.User
	FacebookUserID	uint
}

/* Request */

type UpdateRequest struct {
	Username	string
}

/*
func (UpdateRequest) TableName() string {
    return "user"
}
*/

/* Response */

func (self *User) Response() gin.H {
	return gin.H{
		"id": self.ID,
		"username": self.Username,
	}
}

func (self *User) DetailResponse() (int, gin.H) {
	// TODO: this is just an example, it will be changed base on actual business logic.
	db := database.Shared()
	db.Model(self).Related(&self.FacebookUser, "FacebookUser")
	/*
	if err := db.Model(self).Related(&self.FacebookUser, "FacebookUser").Error; err != nil && err != db.UserNotFound() {
		return herror.Bad(herror.DomainUserResponse, herror.ReasonDatabaseError, err.Error())
	}
	*/
	return http.StatusOK, gin.H{
		"id": self.ID,
		"username": self.Username,
		"facebook_user": self.FacebookUser.Response(),
	}
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
	db := database.Shared()
	db.Where(User{Username: username}).First(&u)
	if u.ID == 0 {
		return nil
	}
	return &u
}