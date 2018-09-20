package error

import (
	//"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Domain string `json:"domain"`
	Reason string `json:"reason"`
	Message string `json:"message"`
}

const DomainUserRegister = "user.register"
const DomainUserLogin = "user.login"

const ReasonDuplicatedEntry = "Duplicated entry"
const ReasonMismatchedEntry = "Mismatched entry"
const ReasonNonexistentEntry = "Nonexistent entry"
const ReasonInvalidJSON = "Invalid JSON"

func Bad(domain string, reason string, message string) (int, gin.H) {
	return http.StatusBadRequest, gin.H{"domain": domain, "reason": reason, "message": message}
}

func Unauthorized(domain string, reason string, message string) (int, gin.H) {
	return http.StatusUnauthorized, gin.H{"domain": domain, "reason": reason, "message": message}
}