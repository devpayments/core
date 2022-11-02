package models

import (
	"github.com/devpayments/common/entity"
	"github.com/google/uuid"
	"time"
)

type Payment struct {
	ID                       string
	SourceTransactionID      string
	DestinationTransactionID string
	Currency                 entity.Currency
	Amount                   int64
	Fee                      int64
	Status                   string
	CreatedAt                time.Time
	UpdatedAt                time.Time
}

type PaymentModel struct {
	ID                       uuid.UUID      `db:"id" goqu:"skipupdate"`
	SourceTransactionID      *uuid.NullUUID `db:"source_transaction_id"`
	DestinationTransactionID *uuid.NullUUID `db:"destination_transaction_id"`
	Currency                 string         `db:"currency" goqu:"skipupdate"`
	Amount                   int64          `db:"amount" goqu:"skipupdate"`
	Fee                      int64          `db:"fee" goqu:"skipupdate"`
	Status                   string         `db:"status"`
	CreatedAt                time.Time      `db:"created_at" goqu:"skipupdate"`
	UpdatedAt                time.Time      `db:"updated_at"`
}

func (m PaymentModel) ToEntity() Payment {
	return Payment{
		ID:                       m.ID.String(),
		SourceTransactionID:      m.SourceTransactionID.UUID.String(),
		DestinationTransactionID: m.DestinationTransactionID.UUID.String(),
		Currency:                 entity.Currency(m.Currency),
		Amount:                   m.Amount,
		Fee:                      m.Fee,
		Status:                   m.Status,
		CreatedAt:                m.CreatedAt,
		UpdatedAt:                m.UpdatedAt,
	}
}

func (m PaymentModel) FromEntity(p Payment) any {
	return PaymentModel{
		ID:                       uuid.MustParse(p.ID),
		Currency:                 string(p.Currency),
		SourceTransactionID:      entity.NewNullUUID(p.SourceTransactionID),
		DestinationTransactionID: entity.NewNullUUID(p.DestinationTransactionID),
		Status:                   p.Status,
		CreatedAt:                p.CreatedAt,
		UpdatedAt:                p.UpdatedAt,
	}
}

type Transaction struct {
	ID                string
	PaymentID         string
	Currency          entity.Currency
	Amount            int64
	Provider          string
	ProviderReference string
	Type              string
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TransactionModel struct {
	ID                uuid.UUID `db:"id"`
	PaymentID         uuid.UUID `db:"payment_id"`
	Currency          string    `db:"currency"`
	Amount            int64     `db:"amount"`
	Provider          string    `db:"provider"`
	ProviderReference string    `db:"provider_reference"`
	Type              string    `db:"type"`
	Status            string    `db:"status"`
	CreatedAt         time.Time `db:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"`
}

func (m TransactionModel) ToEntity() Transaction {
	return Transaction{
		ID:                m.ID.String(),
		PaymentID:         m.PaymentID.String(),
		Currency:          entity.Currency(m.Currency),
		Amount:            m.Amount,
		Provider:          m.Provider,
		ProviderReference: m.ProviderReference,
		Type:              m.Type,
		Status:            m.Status,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func (m TransactionModel) FromEntity(t Transaction) any {
	return TransactionModel{
		ID:                uuid.MustParse(t.ID),
		PaymentID:         uuid.MustParse(t.PaymentID),
		Currency:          string(t.Currency),
		Amount:            t.Amount,
		Provider:          t.Provider,
		ProviderReference: t.ProviderReference,
		Type:              t.Type,
		Status:            t.Status,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}
}
