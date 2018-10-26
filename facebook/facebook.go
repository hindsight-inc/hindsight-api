package main

import (
    "fmt"
    fb "github.com/huandu/facebook"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	fb.Debug = fb.DEBUG_ALL
	var globalApp = fb.New("", "")
    fmt.Println(globalApp)

	res, _ := fb.Get("/me", fb.Params{
		"access_token": "EAAFmf1wKk7MBAHNXpXaZBjfyyGQaPZCXM0ujFZC7iFtMQDCUc6L89TXyDYOZCKB7nDeSbcPXdWd3oZC0nuxJvtsLpV693i6dyDxvAxonv2LYiCYYPw2X2K6ZAyMRs08uzkG3g7hwkqYOxR819MB12Nssta5185RXhJwEk2KaSh7XtCp3kjtijdKSkx929AVGyZA0dZBEoHB5ZCjQtccxPODwupJ6GxZAN95CILPTFyiZAa8UgZDZD",
		"fields": "id,first_name,last_name,middle_name,name,name_format,picture,short_name",
	})
	debugInfo := res.DebugInfo()

	fmt.Println("http headers:", debugInfo.Header)
	fmt.Println("facebook api version:", debugInfo.FacebookApiVersion)
	spew.Dump(res)
	fmt.Println(res["id"])
	fmt.Println(res["first_name"])
	fmt.Println(res["last_name"])
	fmt.Println(res["middle_name"])
	fmt.Println(res["name"])
	fmt.Println(res["name_format"])
	fmt.Println(res["short_name"])
	spew.Dump(res["picture"])
	//fmt.Println(res["picture"]["width"])
	//fmt.Println(res["picture"]["height"])

	/*
    res, _ := fb.Get("/538744468", fb.Params{
        "fields": "first_name",
        "access_token": "a-valid-access-token",
    })
    fmt.Println("Here is my Facebook first name:", res["first_name"])
	*/
}
