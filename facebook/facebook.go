package main

import (
    "fmt"
	"strconv"
    fb "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"hindsight/config"
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

var cfg *config.Configuration
var app *fb.App
var session *fb.Session
var me User

func Init() {
	if _, err := config.Init(); err != nil {
		panic(err)
	}
	cfg = config.Shared()

	//fb.Debug = fb.DEBUG_ALL
	//fb.Version = "v3.0"

	app = fb.New(cfg.Facebook_app_id, cfg.Facebook_app_secret)
    //fmt.Println(app)
}

func UpdateSession() {
	session = app.Session(cfg.Facebook_access_token)
	if err := session.Validate(); err != nil {
		panic(err)
	}
}

func UpdateMe() {
	res, _ := session.Get("/me", fb.Params{
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

	// Picture
	if pic, ok := res["picture"].(map[string] interface {}); ok {
		if data, ok := pic["data"].(map[string] interface {}); ok {
			// Available fields: width, height, url, is_silhouette
			if s, ok := data["url"].(string); ok {
				me.AvatarURL = s
			}
		}
	}
}

func main() {

	Init()
	UpdateSession()
	UpdateMe()
	fmt.Println(me)
}
