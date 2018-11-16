package config

import "github.com/spf13/viper"

type loader interface {
	setConfig(name string, extension string, path string)
	read() error
	merge() error
}

type reader interface {
	getBool(key string) bool
	getString(key string) string
}

type configurator interface {
	loader
	reader
}

type ViperProvider struct{}

func (config ViperProvider) setConfig(name string, extension string, path string) {
	viper.AddConfigPath(path)
	viper.SetConfigType(extension)
	viper.SetConfigName(name)
}

func (config ViperProvider) read() error {
	return viper.ReadInConfig()
}

func (config ViperProvider) merge() error {
	return viper.MergeInConfig()
}

func (config ViperProvider) getBool(key string) bool {
	return viper.GetBool(key)
}

func (config ViperProvider) getString(key string) string {
	return viper.GetString(key)
}
