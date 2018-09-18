package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

//	Don't use old token after changing this, see: https://github.com/appleboy/gin-jwt/issues/170
const IdentityKey = "user.id"

/*
func New(username string) *User {
	return &User{Username: username}
}
*/