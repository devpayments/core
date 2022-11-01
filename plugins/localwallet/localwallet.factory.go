package localwallet

import (
	"github.com/devpayments/common/strategy"
	"github.com/devpayments/core/config"
	"github.com/devpayments/core/datastore/db"
	"github.com/jmoiron/sqlx"
)

type ProviderFactory struct {
	dbCon *sqlx.DB
}

func (f *ProviderFactory) Init() {
	dbConfig := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		User:     "remi",
		Password: "root1234",
		SSLMode:  "disable",
		Name:     "payments",
		Port:     "5432",
	}
	d := db.New(dbConfig)
	dbCon, err := d.GetDbConnection()
	if err != nil {
		panic(err)
	}
	f.dbCon = dbCon
}

func (f *ProviderFactory) Create() strategy.Fundable {

	return NewService(f.dbCon)
}
