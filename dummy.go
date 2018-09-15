package main

import (
	"github.com/gin-gonic/gin"
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