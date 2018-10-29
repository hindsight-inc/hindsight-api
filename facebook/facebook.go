package main

import (
    "fmt"
	"strconv"
	"errors"
    "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"hindsight/config"
	"hindsight/database"
)

type User struct {
	gorm.Model
	FacebookID int64 `json:"facebook_id"`
	Name string `json:"name"`
	ShortName string `json:"short_name"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	MiddleName string `json:"middle_name"`
	NameFormat string `json:"name_format"`
	AvatarURL string `json:"avatar_url"`
}

func (User) TableName() string {
    return "facebook_user"
}

var cfg *config.Configuration
var app *facebook.App
var session *facebook.Session
var me User

func Init() error {
	if _, err := config.Init(); err != nil {
		return err
	}
	cfg = config.Shared()

	//facebook.Debug = facebook.DEBUG_ALL
	//facebook.Version = "v3.0"

	app = facebook.New(cfg.Facebook_app_id, cfg.Facebook_app_secret)
    //fmt.Println(app)
	return nil
}

func UpdateSession(token string) error {
	session = app.Session(token)
	return session.Validate()
}

func UpdateMe() {
	res, _ := session.Get("/me", facebook.Params{
		"fields": "id,first_name,last_name,middle_name,name,name_format,picture,short_name",
	})
	me = User{}

	// Default permissions: https://developers.facebook.com/docs/facebook-login/permissions/#reference-default
	if s, ok := res["id"].(string); ok {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			me.FacebookID = i
		}
	}
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
}

func Create(user User) error {
	var u User
	db := database.GetDB()
	if err := db.Where("facebook_id = ?", user.FacebookID).First(&u).Error; err != nil {
		return err
	}
	if u.ID == 0 {
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	} else {
		return errors.New("Facebook user already exists")
	}
	return nil
}

func Test() {
	db := database.Init()
	db.AutoMigrate(&User{})
	defer db.Close()

	if err := Init(); err == nil {
		if err := UpdateSession(cfg.Facebook_access_token); err == nil {
			UpdateMe()
			fmt.Println(me)
			if err := Create(me); err == nil {
				Create(me)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func main() {
	Test()
}
