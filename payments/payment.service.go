package payments

import (
	"context"

	"github.com/devpayments/common/contracts"
	"github.com/devpayments/common/strategy"
)

type PaymentService struct {
	FundingStrategy strategy.Fundable
	ChargeStrategy  strategy.Chargeable
}

func NewPaymentService() PaymentService {
	return PaymentService{}
}

func (ps *PaymentService) Initiate(source contracts.PaymentSource, destination contracts.PaymentDestination, amount int32, currency string) error {
	return nil
}

func (ps *PaymentService) SetFundingStrategy(ctx context.Context, strategy strategy.Fundable) {
	ps.FundingStrategy = strategy
}
