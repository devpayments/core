package db

import (
	"fmt"
	"github.com/devpayments/core/config"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Db struct {
	DbConfig config.DatabaseConfig
}

func New(dbConfig config.DatabaseConfig) *Db {
	return &Db{DbConfig: dbConfig}
}

func (d *Db) Datasource() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		d.DbConfig.Driver,
		d.DbConfig.User,
		d.DbConfig.Password,
		d.DbConfig.Host,
		d.DbConfig.Port,
		d.DbConfig.Name,
		d.DbConfig.SSLMode,
	)
}

func (d *Db) GetDbConnection() (*sqlx.DB, error) {
	db, err := sqlx.Open(d.DbConfig.Driver, d.Datasource())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db, nil
}

func NewNullUUID(uuidString string) *uuid.NullUUID {
	parsedUUID, err := uuid.Parse(uuidString)
	if err != nil {
		return nil
	}
	return &uuid.NullUUID{UUID: parsedUUID, Valid: true}
}
