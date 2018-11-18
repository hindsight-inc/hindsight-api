package herror

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
const DomainUserUpdate = "user.update"
const DomainFacebookConnect = "facebook.connect"
const DomainTopicCreate = "topic.create"
const DomainUserResponse = "topic.response"
const DomainTopicResponse = "topic.response"
const DomainUploadImage = "upload.image"

const ReasonUnauthorized = "Unauthorized"
const ReasonDuplicatedEntry = "Duplicated Entry"
const ReasonMismatchedEntry = "Mismatched Entry"
const ReasonNonexistentEntry = "Nonexistent Entry"
const ReasonInvalidEntry = "Invalid Entry"
const ReasonEmptyEntry = "Empty Entry"
const ReasonInvalidJSON = "Invalid JSON"
const ReasonDatabaseError = "Database Error"		//	TODO: wrap up db errors?
const ReasonOperationFailure = "Operation Failure"

func Bad(domain string, reason string, message string) (int, gin.H) {
	return http.StatusBadRequest, gin.H{"domain": domain, "reason": reason, "message": message}
}

func Unauthorized(domain string, reason string, message string) (int, gin.H) {
	return http.StatusUnauthorized, gin.H{"domain": domain, "reason": reason, "message": message}
}