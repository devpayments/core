package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/devpayments/common/entity"
	"github.com/devpayments/core/datastore/db"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		SourceTransactionID:      db.NewNullUUID(p.SourceTransactionID),
		DestinationTransactionID: db.NewNullUUID(p.DestinationTransactionID),
		Status:                   p.Status,
		CreatedAt:                p.CreatedAt,
		UpdatedAt:                p.UpdatedAt,
	}
}

type PaymentRepository struct {
	*BaseRepository[PaymentModel, Payment]
}

func NewPaymentRepository(dbCon *sqlx.DB, tableName string) *PaymentRepository {
	baseRepo := NewBaseRepository[PaymentModel, Payment](dbCon, tableName)
	return &PaymentRepository{BaseRepository: baseRepo}
}

func (r *PaymentRepository) Update(ctx context.Context, p *Payment) (int64, error) {
	p.UpdatedAt = time.Now()

	var m PaymentModel
	model := m.FromEntity(*p).(PaymentModel)

	ds := goqu.Update(r.tableName).
		Set(model).
		Where(goqu.Ex{
			"id": model.ID,
		})

	updateSQL, _, err := ds.ToSQL()
	if err != nil {
		panic(err)
	}
	fmt.Println(updateSQL)

	res, err := r.db.ExecContext(ctx, updateSQL)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func (r *PaymentRepository) GetPaymentTransaction(ctx context.Context, txnType string, paymentId uuid.UUID) (*Transaction, error) {
	var transactionField string
	if txnType == "destination" {
		transactionField = "destination_transaction_id"
	} else if txnType == "source" {
		transactionField = "source_transaction_id"
	} else {
		return nil, errors.New("invalid transaction type")
	}

	query, _, err := goqu.
		From(r.tableName).
		Join(
			goqu.T("transactions"),
			goqu.On(
				goqu.I(transactionField).Eq(goqu.I("transactions.id")),
			),
		).
		Select("transactions.*").
		Where(goqu.Ex{"payments.id": paymentId.String()}).
		ToSQL()

	row := r.db.QueryRowxContext(ctx, query)

	var m TransactionModel
	err = row.StructScan(&m)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	t := m.ToEntity()
	return &t, nil
}

func (r *PaymentRepository) SetPaymentTransaction(ctx context.Context, paymentId uuid.UUID, transaction *Transaction) error {
	var transactionField string
	if transaction.Type == "destination" {
		transactionField = "destination_transaction_id"
	} else if transaction.Type == "source" {
		transactionField = "source_transaction_id"
	} else {
		return errors.New("invalid transaction type")
	}

	var m TransactionModel
	model := m.FromEntity(*transaction)

	// Create New Transaction
	insertTransactionSQL, _, err := goqu.Insert("transactions").Rows(
		model,
	).ToSQL()
	if err != nil {
		panic(err)
	}

	_, err = r.db.ExecContext(ctx, insertTransactionSQL)
	if err != nil {
		return err
	}

	// Update Payment Transaction
	updatePaymentSQL, _, err := goqu.
		Update("payments").
		Set(map[string]interface{}{
			transactionField: transaction.ID,
		}).
		Where(goqu.Ex{"id": paymentId.String()}).
		ToSQL()
	_, err = r.db.ExecContext(ctx, updatePaymentSQL)
	if err != nil {
		panic(err)
	}

	return nil
}
