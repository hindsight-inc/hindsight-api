package config

type Configuration struct {
	HTTPPort string

	MySQLDatabase string
	MySQLPassword string

	JWTRealm string
	JWTKey   string

	Facebook_disable_test bool
	FacebookAppID       string
	FacebookAppSecret   string
	FacebookAccessToken string
}

var Config *Configuration

func loadConfiguration(config loader) error {

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

func readConfiguration(config reader) (*Configuration, error) {
	Config = &Configuration{
		HTTPPort: config.getString("http_port"),

		MySQLDatabase: config.getString("mysql_database"),
		MySQLPassword: config.getString("mysql_password"),

		JWTRealm: config.getString("jwt_realm"),
		JWTKey:   config.getString("jwt_key"),

		Facebook_disable_test: config.getBool("fb_disable_test"),
		FacebookAppID:       config.getString("fb_app_id"),
		FacebookAppSecret:   config.getString("fb_app_secret"),
		FacebookAccessToken: config.getString("fb_access_token"),
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
