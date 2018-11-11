package error

import (
	//"errors"
	"net/http"
	"github.com/gin-gonic/gin"
)

type APIError struct {
	Domain string
	Reason string
	Message string
}

const DomainAuthJWT = "auth.jwt"
const DomainUserRegister = "user.register"
const DomainUserLogin = "user.login"
const DomainUserInfo = "user.info"
const DomainFacebookConnect = "facebook.connect"
const DomainTopicCreate = "topic.create"

const ReasonUnauthorized = "Unauthorized"
const ReasonDuplicatedEntry = "Duplicated Entry"
const ReasonMismatchedEntry = "Mismatched Entry"
const ReasonNonexistentEntry = "Nonexistent Entry"
const ReasonInvalidEntry = "Invalid Entry"
const ReasonEmptyEntry = "Empty Entry"
const ReasonInvalidJSON = "Invalid JSON"
const ReasonDatabaseError = "Database Error"

func Bad(domain string, reason string, message string) (int, gin.H) {
	return http.StatusBadRequest, gin.H{"domain": domain, "reason": reason, "message": message}
}

func Unauthorized(domain string, reason string, message string) (int, gin.H) {
	return http.StatusUnauthorized, gin.H{"domain": domain, "reason": reason, "message": message}
}