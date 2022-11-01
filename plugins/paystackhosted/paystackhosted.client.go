package paystackhosted

import (
	"context"
	"github.com/devpayments/common/entity"
	"github.com/devpayments/common/httpclient"
)

type InitiatePaymentRequestData struct {
	Currency  entity.Currency `json:"currency"`
	Amount    int64           `json:"amount"`
	Reference string          `json:"reference"`
	Email     string          `json:"email"`
}

type InitiatePaymentResponseData struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type APIClient struct {
	secretKey string
	publicKey string
	baseURL   string
}

func NewAPIClient(secretKey string, publicKey string, baseURL string) *APIClient {
	return &APIClient{secretKey: secretKey, publicKey: publicKey, baseURL: baseURL}
}

func (a APIClient) InitiateTransaction(ctx context.Context, request InitiatePaymentRequestData) (*InitiatePaymentResponseData, error) {
	httpRequest := httpclient.NewHttpRequest("https://api.paystack.co/transaction/initialize", "POST")

	httpRequest.SetBody(request)
	httpRequest.SetAuthToken(a.secretKey)

	var response InitiatePaymentResponseData
	err := httpclient.MakeApiCall(ctx, httpRequest, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
