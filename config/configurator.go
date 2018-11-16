package config

import "github.com/spf13/viper"

type configurationLoader interface {
	setConfig(name string, extension string, path string)
	read() error
	merge() error
}

type configurationReader interface {
	getBool(key string) bool
	getString(key string) string
}

type configurator interface {
	configurationLoader
	configurationReader
}

type ConfigProvider struct {
}

func (config ConfigProvider) setConfig(name string, extension string, path string) {
	viper.AddConfigPath(path)
	viper.SetConfigType(extension)
	viper.SetConfigName(name)
}

func (config ConfigProvider) read() error {
	return viper.ReadInConfig()
}

func (config ConfigProvider) merge() error {
	return viper.MergeInConfig()
}

func (config ConfigProvider) getBool(key string) bool {
	return viper.GetBool(key)
}

func (config ConfigProvider) getString(key string) string {
	return viper.GetString(key)
}
