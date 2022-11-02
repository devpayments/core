package localwallet

import (
	"fmt"
	"github.com/devpayments/common/strategy"
	"github.com/jmoiron/sqlx"
)

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

type ProviderFactory struct {
	dbCon *sqlx.DB
}

func (f *ProviderFactory) Init() {
	dbConfig := DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		User:     "remi",
		Password: "root1234",
		SSLMode:  "disable",
		Name:     "payments",
		Port:     "5432",
	}

	datasource := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.Driver,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSLMode,
	)

	dbCon, err := sqlx.Open(dbConfig.Driver, datasource)
	if err != nil {
		panic(err)
	}

	err = dbCon.Ping()
	if err != nil {
		panic(err)
	}

	f.dbCon = dbCon
}

func (f *ProviderFactory) Create() strategy.Fundable {

	return NewService(f.dbCon)
}
