package paystackhosted

import (
	"context"
	"github.com/devpayments/common/entity"
	"github.com/google/uuid"
)

type Service struct {
	apiClient APIClient
}

func NewService(apiClient APIClient) *Service {
	return &Service{apiClient: apiClient}
}

func (s *Service) Name() string {
	return "paystackhosted"
}

func (s *Service) InitiateCharge(ctx context.Context, source any, amount int64, currency entity.Currency) (entity.Transaction, error) {
	requestData := InitiatePaymentRequestData{
		Currency:  currency,
		Amount:    amount,
		Reference: uuid.New().String(),
		Email:     "test@mail.com",
	}
	s.apiClient.InitiateTransaction(ctx, requestData)

	t := entity.Transaction{
		Type:      s.Name(),
		Status:    "initiated",
		Currency:  currency,
		Amount:    amount,
		Reference: requestData.Reference,
	}
	return t, nil
}

func (s *Service) GetChargeAuthorization(ctx context.Context, reference string) error {
	return nil
}

func (s *Service) AuthorizeCharge(ctx context.Context, reference string, authorizationData any) error {
	return nil
}

func (s *Service) CheckChargeStatus(ctx context.Context, reference string) error {
	return nil
}

func (s *Service) CompleteCharge(ctx context.Context, reference string) error {
	return nil
}
