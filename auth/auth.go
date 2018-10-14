package auth

import (
	"log"
	"time"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"hindsight/user"
	apiError "hindsight/error"
)

type Token struct {
	Expire string `json:"expire"`
	Token string `json:"token"`
}

func GetMiddleware() *jwt.GinJWTMiddleware {
	//https://godoc.org/gopkg.in/appleboy/gin-jwt.v2
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "hindsight-inc",
		Key:         []byte("secret key is required"),
		Timeout:     time.Hour * 99999,
		MaxRefresh:  time.Hour * 99999,
		IdentityKey: user.IdentityKey,

		//	TODO: replace user.Username with user.ID? Need to understand more about Claim
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					user.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		//	see above
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[user.IdentityKey].(string)
			return &user.User{
				Username: username,
			}
		},
		//	login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			_, response, user := user.Authenticate(c)
			json, _ := json.Marshal(response)
			if user == nil {
				return nil, errors.New(string(json))
			}
			return user, nil
		},
		//	access control
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if u, ok := data.(*user.User); ok {
				return user.Authorize(c, u.Username)
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"domain": apiError.DomainAuthJWT,
				"reason": apiError.ReasonUnauthorized,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	return middleware
}