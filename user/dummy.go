package user

import (
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
)

//	Test only dummy code

func DummyUsersList(context *gin.Context) {
	user := User{Username: "username002"}
	context.JSON(200, user.Response())
}

func DummyUsersInfo(context *gin.Context) {
	user := User{Username: "username002"}
	context.JSON(200, user.Response())
}

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(IdentityKey)
	c.JSON(200, gin.H{
		"userID": claims[IdentityKey],
		"user": user,
		"username": user.(*User).Username,
		"text": "Hello World.",
	})
}