package error

import (
	//"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

const DomainUserRegister = "user.register"
const DomainUserLogin = "user.login"

const ReasonDuplicatedEntry = "Duplicated entry"
const ReasonInvalidJSON = "Invalid JSON"

func Bad(domain string, reason string, message string) (int, gin.H) {
	return http.StatusBadRequest, gin.H{"domain": domain, "reason": reason, "message": message}
}

/*
type ApiError struct {
	code int
	domain string
	reason string
}

var apiErrors = map[string]ApiError {
	"user.register.existing": ApiError{http.StatusBadRequest, "User already exists"},
}

func New(eid string) (int, gin.H) {
	e := apiErrors[eid]
	return e.code, gin.H{"id": eid, "reason": e.reason}
}
*/