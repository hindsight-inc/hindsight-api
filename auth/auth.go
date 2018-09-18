package auth

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/appleboy/gin-jwt"
	"hindsight/user"
)

func GetMiddleware() *jwt.GinJWTMiddleware {
	middleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key is required"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: user.IdentityKey,

		//	get identity from json, i.e. Username
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					user.IdentityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		//	get user from identity
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			username := claims[user.IdentityKey].(string)
			return &user.User{
				Username: username,
			}
		},
		//	login: `admin` and `test` can login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var u user.User
			if err := c.ShouldBind(&u); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := u.Username
			password := u.Password

			if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
				return &user.User{
					Username: username,
					Password: password,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		//	access control: `admin` is authorized
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*user.User); ok && v.Username == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
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
	log.Println("xxx")
	return middleware
}

/*
curl -v POST \
  http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{ "username": "admin", "password": "admin" }'

curl -v GET \
  http://localhost:8080/auth/refresh_token \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'

curl -v GET \
  http://localhost:8080/auth/hello \
  -H 'content-type: application/json' \
  -H 'Authorization:Bearer xxx'
*/