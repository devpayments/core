package payments

import (
	"context"
	"errors"
	"fmt"
	"github.com/devpayments/common/entity"
	"github.com/devpayments/core/models"
	"github.com/google/uuid"
	"time"
)

type Repository[M entity.Model[E], E entity.Entity] interface {
	Create(ctx context.Context, entity *E) (int64, error)
	FindByID(ctx context.Context, id uuid.UUID) (*E, error)
	FindOne(ctx context.Context, whereMap map[string]any) (*E, error)
	FindAll(ctx context.Context, whereMap map[string]any) (*E, error)
}

type PaymentsRepository interface {
	Repository[models.PaymentModel, models.Payment]
	Update(ctx context.Context, p *models.Payment) (int64, error)
	GetPaymentTransaction(ctx context.Context, txnType string, paymentId uuid.UUID) (*models.Transaction, error)
	SetPaymentTransaction(ctx context.Context, paymentId uuid.UUID, transaction *models.Transaction) error
}

type TransactionsRepository interface {
	Repository[models.TransactionModel, models.Transaction]
}

type PaymentService struct {
	paymentRepository     PaymentsRepository
	transactionRepository TransactionsRepository
}

func NewPaymentService(paymentRepository PaymentsRepository, transactionRepository TransactionsRepository) PaymentService {
	return PaymentService{
		transactionRepository: transactionRepository,
		paymentRepository:     paymentRepository,
	}
}

func (ps *PaymentService) Initiate(ctx context.Context) (*models.Payment, error) {
	// Initiate payment
	payment := models.Payment{
		ID:        uuid.New().String(),
		Status:    "initiated",
		Currency:  entity.Currency("NGN"),
		Amount:    10000,
		Fee:       100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := ps.paymentRepository.Create(ctx, &payment)
	if err != nil {
		panic(err)
	}

	// Do some logic to select strategy or strategy would be provided to this function
	chargeHandler := NewChargeable("paystackhosted")

	// Initiate transaction
	t, err := chargeHandler.InitiateCharge(ctx, 1, payment.Amount, payment.Currency)
	initiateTransaction := models.Transaction{
		ID:                uuid.New().String(),
		PaymentID:         payment.ID,
		Currency:          t.Currency,
		Amount:            t.Amount,
		Provider:          chargeHandler.Name(),
		ProviderReference: t.Reference,
		Type:              "source",
		Status:            t.Status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = ps.paymentRepository.SetPaymentTransaction(ctx, uuid.MustParse(payment.ID), &initiateTransaction)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (ps *PaymentService) GetAuthorization(ctx context.Context, paymentId uuid.UUID) error {
	sourceTransaction, err := ps.paymentRepository.GetPaymentTransaction(ctx, "source", paymentId)
	if err != nil {
		return err
	}

	chargeHandler := NewChargeable(sourceTransaction.Provider)
	chargeHandler.GetChargeAuthorization(ctx, sourceTransaction.ProviderReference)

	fmt.Println(sourceTransaction)

	// Authorization type
	// otp, pin, redirect,
	return nil
}

func (ps *PaymentService) CompleteAuthorization(ctx context.Context, paymentId uuid.UUID, authorizationData any) error {
	sourceTransaction, err := ps.paymentRepository.GetPaymentTransaction(ctx, "source", paymentId)
	if err != nil {
		return err
	}

	chargeHandler := NewChargeable(sourceTransaction.Provider)
	chargeHandler.AuthorizeCharge(ctx, sourceTransaction.ProviderReference, authorizationData)

	return nil
}

func (ps *PaymentService) Complete(ctx context.Context, paymentId uuid.UUID) error {
	sourceTransaction, err := ps.paymentRepository.GetPaymentTransaction(ctx, "source", paymentId)
	if err != nil {
		return err
	}

	if sourceTransaction.Status != "success" {
		return errors.New("transaction is not successful")
	}

	wallet := entity.Wallet{
		Identifier: "0012923234",
	}
	fundable := NewFundable("localwallet")
	t, err := fundable.InitiateFunding(ctx, wallet, 10000, entity.Currency("NGN"))
	if err != nil {
		return err
	}

	destinationTransaction := models.Transaction{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Provider:  fundable.Name(),
		Type:      "destination",
		PaymentID: paymentId.String(),
		Status:    t.Status,
	}
	err = ps.paymentRepository.SetPaymentTransaction(ctx, paymentId, &destinationTransaction)
	if err != nil {
		return err
	}

	return nil
}
