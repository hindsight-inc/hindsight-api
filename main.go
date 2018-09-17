package main

import (
	"log"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/appleboy/gin-jwt"
	"hindsight/database"
	"hindsight/topic"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
}

func (self *User) Response() gin.H {
	return gin.H{"username": self.Username}
}

///	main

/*
curl -v -X POST \
  http://localhost:8080/user/register \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserRegister(context *gin.Context) {
	db := database.GetDB()
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&User{Username: user.Username, Password: user.Password})
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

/*
curl -v -X POST \
  http://localhost:8080/user/login \
  -H 'content-type: application/json' \
  -d '{ "username": "username001", "password": "password001" }'
*/
func UserLogin(context *gin.Context) {
	var user User
	if err := context.ShouldBindJSON(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.Username != "username001" || user.Password != "password001" {
		context.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"status": "success"})
}

func UserInfo(context *gin.Context) {
	user := User{Username: "username002"}
	context.JSON(200, user.Response())
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
func helloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get("uid")
	c.JSON(200, gin.H{
		"userID":   claims["uid"],
		"username": user.(*User).Username,
		"text":     "Hello World.",
	})
}

func main() {
	db := database.Init()
	db.AutoMigrate(&User{})
	db.AutoMigrate(&topic.Topic{})
	defer db.Close()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key is required"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "uid",

		//	get identity from json, i.e. Username
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"uid": v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		//	get user from identity
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Username: claims["uid"].(string),
			}
		},
		//	login: `admin` and `test` can login
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var user User
			if err := c.ShouldBind(&user); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := user.Username
			password := user.Password

			if (username == "admin" && password == "admin") || (username == "test" && password == "test") {
				return &User{
					Username: username,
					Password: password,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		//	access control: `admin` is authorized
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.Username == "admin" {
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

	//	route
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//	auth
	r.POST("/login", authMiddleware.LoginHandler)
	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := r.Group("/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", helloHandler)
		auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	}

	//	TODO: FB user study
	r.POST("/user/register", UserRegister)
	r.POST("/user/login", UserLogin)
	r.GET("/user", UserInfo)

	//	TODO: dummy code here for now
	r.GET("/users", DummyUsersList)
	r.GET("/users/:uid", DummyUsersInfo)

	r.GET("/topics", topic.List)
	r.GET("/topics/:id", topic.Detail)
	r.POST("/topics", topic.Create)

	r.Run()
}