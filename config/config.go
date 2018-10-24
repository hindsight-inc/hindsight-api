package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	viper.SetDefault("fb_app_id", "FACEBOOK_APP_ID")
	viper.SetDefault("fb_app_secret", "FACEBOOK_SECRET")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetConfigName("secret")
	if err := viper.MergeInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	fmt.Println("app id: " + viper.GetString("fb_app_id"))
	fmt.Println("app secret: " + viper.GetString("fb_app_secret"))
}
