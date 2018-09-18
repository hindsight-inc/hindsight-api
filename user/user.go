package user

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

/*
func New(username string) *User {
	return &User{Username: username}
}
*/