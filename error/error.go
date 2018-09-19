package error

import (
	//"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

type ApiError struct {
	code int
	id string
	message string
}

var UserRegisterExisting = ApiError{http.StatusBadRequest, "error.user.register.existing", "User already exists"}

func H(e ApiError) (int, gin.H) {
	return e.code, gin.H{"id": e.id, "message": e.message}
}

/*
var apiErrors = map[string]ApiError {
	"user.register.existing": ApiError{http.StatusBadRequest, "User already exists"},
}

func New(eid string) (int, gin.H) {
	e := apiErrors[eid]
	return e.code, gin.H{"id": eid, "message": e.message}
}
*/