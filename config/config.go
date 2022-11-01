package config

type Config struct {
	Database DatabaseConfig
	App      AppConfig
}

type DatabaseConfig struct {
	Driver         string `envconfig:"DB_DRIVER"`
	Host           string `envconfig:"DB_HOST"`
	Port           string `envconfig:"DB_PORT"`
	User           string `envconfig:"DB_USER"`
	Password       string `envconfig:"DB_PASSWORD"`
	Name           string `envconfig:"DB_NAME"`
	SSLMode        string `envconfig:"DB_SSL_MODE"`
	SearchPath     string `envconfig:"DB_SEARCH_PATH"`
	RedisURL       string `envconfig:"REDIS_URL"`
	RedisPort      string `envconfig:"REDIS_PORT"`
	RedisPassword  string `envconfig:"REDIS_PASSWORD"`
	RedisNamespace string `envconfig:"REDIS_NAMESPACE"`
}

type AppConfig struct {
	Environment string `envconfig:"APP_ENV"`
	PORT        string `envconfig:"APP_PORT"`
}

//func Load() (*Config, error) {
//	if err := godotenv.Load(".env"); err != nil {
//		panic(err)
//	}
//
//	var con Config
//	if err := envconfig.Process("", &con); err != nil {
//		panic(err)
//	}
//
//	return &con, nil
//}
