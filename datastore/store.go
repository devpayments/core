package datastore

import (
	"github.com/devpayments/core/datastore/repository"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db           sqlx.DB
	Transactions *repository.TransactionRepository
	Payments     *repository.PaymentRepository
}

func NewStore(db sqlx.DB) *Store {
	return &Store{
		db:           db,
		Transactions: repository.NewTransactionRepository(&db, "transactions"),
		Payments:     repository.NewPaymentRepository(&db, "payments"),
	}
}

func (s *Store) GetDbConnection() sqlx.DB {
	return s.db
}
