package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"time"
)

type Transaction struct {
	ID                string
	PaymentID         string
	Currency          string
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
		Currency:          m.Currency,
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
		Currency:          m.Currency,
		Amount:            m.Amount,
		Provider:          t.Provider,
		ProviderReference: t.ProviderReference,
		Type:              t.Type,
		Status:            t.Status,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}
}

type TransactionRepository struct {
	*BaseRepository[TransactionModel, Transaction]
}

func NewTransactionRepository(dbCon *sqlx.DB, tableName string) *TransactionRepository {
	baseRepo := NewBaseRepository[TransactionModel, Transaction](dbCon, tableName)
	return &TransactionRepository{BaseRepository: baseRepo}
}
