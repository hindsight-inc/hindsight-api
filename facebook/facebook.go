package facebook

import (
	"strconv"
	"errors"
    "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"hindsight/config"
	"hindsight/database"
)

type User struct {
	gorm.Model
	FacebookID int64
	Name string
	ShortName string
	FirstName string
	LastName string
	MiddleName string
	NameFormat string
	AvatarURL string
}

func (User) TableName() string {
    return "facebook_user"
}

func (self *User) UniqueUsername() string {
	guid := xid.New()
	return "fb_" + strconv.FormatInt(self.FacebookID, 10) + "_" + guid.String()
}

/* Request */

type ConnectRequest struct {
	AccessToken string `json:"access_token"`
}

/* Response */

func (self *User) Response() gin.H {
	return gin.H{
		"id": self.ID,
		"name": self.Name,
		"avatar_url": self.AvatarURL,
	}
}

// local

var db *gorm.DB
var cfg *config.Configuration
var app *facebook.App
var session *facebook.Session
var me User

/** 
 Initialize facebook app. 
 Requires database initialization.
 */
func Init() error {
	db = database.Shared()
	/*
	if _, err := config.Init(); err != nil {
		return err
	}
	*/
	cfg = config.Shared()

	//facebook.Debug = facebook.DEBUG_ALL
	//facebook.Version = "v3.0"

	app = facebook.New(cfg.FacebookAppID, cfg.FacebookAppSecret)
	return nil
}

func updateSession(token string) error {
	session = app.Session(token)
	return session.Validate()
}

func updateMe() error {
	var res facebook.Result
	var err error
	if res, err = session.Get("/me", facebook.Params{
		//"fields": "id,first_name,last_name,middle_name,name,name_format,picture,short_name",
		"fields": "id,first_name,last_name,middle_name,name,name_format,picture",
	}); err != nil {
		return err
	}
	me = User{}

	// Default permissions: https://developers.facebook.com/docs/facebook-login/permissions/#reference-default
	if s, ok := res["id"].(string); ok {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			me.FacebookID = i
		}
	}
	if me.FacebookID <= 0 {
		return errors.New("Invalid Facebook ID")
	}

	// The following fields are all optional
	if s, ok := res["first_name"].(string); ok {
		me.FirstName = s
	}
	if s, ok := res["middle_name"].(string); ok {
		me.MiddleName = s
	}
	if s, ok := res["last_name"].(string); ok {
		me.LastName = s
	}
	if s, ok := res["name"].(string); ok {
		me.Name = s
	}
	if s, ok := res["short_name"].(string); ok {
		me.ShortName = s
	}
	if s, ok := res["name_format"].(string); ok {
		me.NameFormat = s
	}

	// Picture available fields: width, height, url, is_silhouette
	if pic, ok := res["picture"].(map[string] interface {}); ok {
		if data, ok := pic["data"].(map[string] interface {}); ok {
			if s, ok := data["url"].(string); ok {
				me.AvatarURL = s
			}
		}
	}
	return nil
}

/**
 Create user if it's not existing.
 This function is not used right now. Remove?
 */
func create(user User) error {
	var u User
	if notFound := db.Where(User{FacebookID: user.FacebookID}).First(&u).RecordNotFound(); !notFound {
		return errors.New("Facebook user already exists")
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

/**
 Get facebook.User based on Facebook ID, which is from Facebook access token.
 Create and return the user if it is not created yet.
 This function does *NOT* tell whether a user is new or existing.
 */
func Connect(token string) (User, error) {
	var u User
	if err := updateSession(token); err != nil {
		return u, err
	}
	if err := updateMe(); err != nil {
		return u, err
	}
	if me.FacebookID <= 0 {
		return u, errors.New("Invalid Facebook ID")
	}
	if notFound := db.Where(User{FacebookID: me.FacebookID}).First(&u).RecordNotFound(); !notFound {
		// existing user
		return u, nil
	}
	if err := db.Create(&me).Error; err != nil {
		return u, err
	}
	// new user
	return me, nil
}
