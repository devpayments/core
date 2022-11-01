package payments

import (
	"context"
	"errors"
	"fmt"
	"github.com/devpayments/common/entity"
	"github.com/devpayments/core/datastore"
	"github.com/devpayments/core/datastore/repository"
	"github.com/google/uuid"
	"time"
)

type PaymentService struct {
	store datastore.Store
}

func NewPaymentService(store datastore.Store) PaymentService {
	return PaymentService{
		store: store,
	}
}

func (ps *PaymentService) Initiate(ctx context.Context) (*repository.Payment, error) {

	// Initiate payment
	payment := repository.Payment{
		ID:        uuid.New().String(),
		Status:    "initiated",
		Currency:  entity.Currency("NGN"),
		Amount:    10000,
		Fee:       100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := ps.store.Payments.Create(ctx, &payment)
	if err != nil {
		panic(err)
	}

	// Do some logic to select strategy or strategy would be provided to this function
	chargeHandler := NewChargeable("paystackhosted")

	// Initiate transaction
	t, err := chargeHandler.InitiateCharge(ctx, 1, payment.Amount, payment.Currency)
	initiateTransaction := repository.Transaction{
		ID:                uuid.New().String(),
		PaymentID:         payment.ID,
		Currency:          string(t.Currency),
		Amount:            t.Amount,
		Provider:          chargeHandler.Name(),
		ProviderReference: t.Reference,
		Type:              "source",
		Status:            t.Status,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err = ps.store.Payments.SetPaymentTransaction(ctx, uuid.MustParse(payment.ID), &initiateTransaction)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (ps *PaymentService) GetAuthorization(ctx context.Context, paymentId uuid.UUID) error {
	sourceTransaction, err := ps.store.Payments.GetPaymentTransaction(ctx, "source", paymentId)
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
	sourceTransaction, err := ps.store.Payments.GetPaymentTransaction(ctx, "source", paymentId)
	if err != nil {
		return err
	}

	chargeHandler := NewChargeable(sourceTransaction.Provider)
	chargeHandler.AuthorizeCharge(ctx, sourceTransaction.ProviderReference, authorizationData)

	return nil
}

func (ps *PaymentService) Complete(ctx context.Context, paymentId uuid.UUID) error {
	sourceTransaction, err := ps.store.Payments.GetPaymentTransaction(ctx, "source", paymentId)
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

	destinationTransaction := repository.Transaction{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Provider:  fundable.Name(),
		Type:      "destination",
		PaymentID: paymentId.String(),
		Status:    t.Status,
	}
	err = ps.store.Payments.SetPaymentTransaction(ctx, paymentId, &destinationTransaction)
	if err != nil {
		return err
	}

	return nil
}
