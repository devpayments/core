package paystackhosted

import "github.com/devpayments/common/strategy"

type Providerfactory struct {
	client APIClient
}

func (f *Providerfactory) Init() {

}

func (f *Providerfactory) Create() strategy.Chargeable {
	apiClient := NewAPIClient(
		"sk_test_b0c6698dc9c3f462afb975f000544b8a9bd27983",
		"pk_test_5054c85c1bc06f096fe5f57a208cabf11a1ca15a",
		"https://api.paystack.co",
	)

	return NewService(*apiClient)
}
