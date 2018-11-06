package config

import (
	//"fmt"
	"github.com/spf13/viper"
)

type Configuration struct {
	MySQL_database string
	MySQL_password string

	JWT_Realm string
	JWT_Key string

	Facebook_app_id	string
	Facebook_app_secret string
	Facebook_access_token string
}

var Config *Configuration

func loadViper() error {
	//viper.SetDefault("fb_app_id", "FACEBOOK_APP_ID")
	//viper.SetDefault("fb_app_secret", "FACEBOOK_SECRET")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath("../config/")
	viper.SetConfigType("yaml")

	//	load normal config
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		//panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return err
	}

	//	load secret config
	viper.SetConfigName("secret")
	if err := viper.MergeInConfig(); err != nil {
		//panic(fmt.Errorf("Fatal error secret file: %s \n", err))
		return err
	}
	return nil
}

func Init() (*Configuration, error) {
	// `viper` is used as our configuration management library
	if err := loadViper(); err != nil {
		return nil, err
	}
	Config = &Configuration{
		MySQL_database: viper.GetString("mysql_database"),
		MySQL_password: viper.GetString("mysql_password"),

		JWT_Realm: viper.GetString("jwt_realm"),
		JWT_Key: viper.GetString("jwt_key"),

		Facebook_app_id: viper.GetString("fb_app_id"),
		Facebook_app_secret: viper.GetString("fb_app_secret"),
		Facebook_access_token: viper.GetString("fb_access_token"),
	}
	return Config, nil
}

func Shared() *Configuration {
	return Config
}