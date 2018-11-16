package config

type Configuration struct {
	MySQL_database string
	MySQL_password string

	JWT_Realm string
	JWT_Key   string

	Facebook_disable_test bool
	Facebook_app_id       string
	Facebook_app_secret   string
	Facebook_access_token string
}

var Config *Configuration

func loadConfiguration(config configurationLoader) error {

	config.setConfig("config", "yaml", ".")
	config.setConfig("config", "yaml", "./config/")
	config.setConfig("config", "yaml", "../config/")
	if err := config.read(); err != nil {
		return err
	}

	config.setConfig("secret", "yaml", ".")
	if err := config.merge(); err != nil {
		return err
	}
	return nil
}

func readConfiguration(config configurationReader) (*Configuration, error) {
	Config = &Configuration{
		MySQL_database: config.getString("mysql_database"),
		MySQL_password: config.getString("mysql_password"),

		JWT_Realm: config.getString("jwt_realm"),
		JWT_Key:   config.getString("jwt_key"),

		Facebook_disable_test: config.getBool("fb_disable_test"),
		Facebook_app_id:       config.getString("fb_app_id"),
		Facebook_app_secret:   config.getString("fb_app_secret"),
		Facebook_access_token: config.getString("fb_access_token"),
	}
	return Config, nil
}

func Init(config configurator) (*Configuration, error) {

	if err := loadConfiguration(config); err != nil {
		return nil, err
	}

	return readConfiguration(config)
}

func Shared() *Configuration {
	return Config
}
