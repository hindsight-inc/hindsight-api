package main

import (
    "fmt"
    fb "github.com/huandu/facebook"
)

func main() {
	fb.Debug = fb.DEBUG_ALL
	var globalApp = fb.New("394172167787443", "")
    fmt.Println(globalApp)

	res, _ := fb.Get("/me", fb.Params{"access_token": "xxx"})
	debugInfo := res.DebugInfo()

	fmt.Println("http headers:", debugInfo.Header)
	fmt.Println("facebook api version:", debugInfo.FacebookApiVersion)

	/*
    res, _ := fb.Get("/538744468", fb.Params{
        "fields": "first_name",
        "access_token": "a-valid-access-token",
    })
    fmt.Println("Here is my Facebook first name:", res["first_name"])
	*/
}