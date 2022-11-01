package localwallet

import (
	"context"
	"database/sql"
	"errors"
	"github.com/devpayments/common/entity"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	dbCon *sqlx.DB
}

func NewService(dbCon *sqlx.DB) *Service {
	return &Service{dbCon: dbCon}
}

func loadWallet(ctx context.Context, dbCon *sqlx.DB, wallet *entity.Wallet) error {
	selectSQL, _, err := goqu.From("wallets").
		Where(goqu.Ex{
			"identifier": wallet.Identifier,
		}).
		ToSQL()
	if err != nil {
		return err
	}

	row := dbCon.QueryRowxContext(ctx, selectSQL)
	err = row.StructScan(wallet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return nil
}

func (ls *Service) Name() string {
	return "localwallet"
}

func (ls *Service) InitiateFunding(ctx context.Context, destination any, amount int64, currency entity.Currency) (entity.Transaction, error) {
	t := entity.Transaction{
		Type:      ls.Name(),
		Status:    "initiated",
		Amount:    amount,
		Currency:  currency,
		Reference: uuid.New().String(),
	}

	// get wallet
	wallet := destination.(entity.Wallet)
	loadWallet(ctx, ls.dbCon, &wallet)

	return t, nil
}

func (ls *Service) CompleteFunding(ctx context.Context, reference string) error {
	selectSQL, _, err := goqu.From("transactions").
		Where(goqu.Ex{
			"reference": reference,
		}).
		ToSQL()
	if err != nil {
		return err
	}

	row := ls.dbCon.QueryRowxContext(ctx, selectSQL)
	tr := entity.Transaction{}
	err = row.StructScan(tr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return nil
}
