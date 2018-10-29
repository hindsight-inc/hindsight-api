package main

import (
    "fmt"
	"strconv"
	//"reflect"
    fb "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	//"github.com/davecgh/go-spew/spew"
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
    fmt.Println(app)
}

func UpdateSession() {
	session = app.Session(cfg.Facebook_access_token)
	if err := session.Validate(); err != nil {
		panic(err)
	}
}

func UpdateMe() {
	//res, _ := session.Get("/me", nil)
	res, _ := session.Get("/me", fb.Params{
		"fields": "id,first_name,last_name,middle_name,name,name_format,picture,short_name",
	})
	me = User{}
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

	/*
	res, _ := fb.Get("/me", fb.Params{
		"access_token": "EAAFWnwbpvnwBAFBGhgosZBmpiKKPOdjLe6PiswomwuejvOEd9M77qPDEMhcCWgQsZCsJJbZBiDstXq4ZAlK0HFZBuR71shdUOpKQa8lFWh1qZC83ct5m8qtsLY9y0l6wEZCZBiCmUbMhfEXEeqD65RK1Yo1qx28WrSjqeEdRxBaISNeZA82BvxEcJNhYdInZAPwPJ5IQdyNyaqngZDZD",
		"fields": "id,first_name,last_name,middle_name,name,name_format,picture,short_name",
	})
	debugInfo := res.DebugInfo()

	fmt.Println("http headers:", debugInfo.Header)
	fmt.Println("facebook api version:", debugInfo.FacebookApiVersion)
	*/
	//spew.Dump(res)
	fmt.Println(res["id"])
	fmt.Println(res["first_name"])
	fmt.Println(res["last_name"])
	fmt.Println(res["middle_name"])
	fmt.Println(res["name"])
	fmt.Println(res["name_format"])
	fmt.Println(res["short_name"])
	//spew.Dump(res["picture"])
	//fmt.Println(res["picture"])

	pic, ok := res["picture"].(map[string] interface {})
	if !ok {
		panic("invalid picture")
	}
	//fmt.Println(pic["data"])
	data, ok := pic["data"].(map[string] interface {})
	if !ok {
		panic("invalid picture")
	}
	fmt.Println(data["width"])
	fmt.Println(data["height"])
	fmt.Println(data["url"])
	fmt.Println(data["is_silhouette"])
}

func main() {

	Init()
	UpdateSession()
	UpdateMe()
	fmt.Println(me)

	/*
    res, _ := fb.Get("/538744468", fb.Params{
        "fields": "first_name",
        "access_token": "a-valid-access-token",
    })
    fmt.Println("Here is my Facebook first name:", res["first_name"])
	*/
}
