package payments

import (
	"errors"
	"github.com/devpayments/common/strategy"
	"github.com/devpayments/core/plugins/localwallet"
	"github.com/devpayments/core/plugins/paystackhosted"
)

type ChargeProviderFactory interface {
	Init()
	Create() strategy.Chargeable
}

type FundProviderFactory interface {
	Init()
	Create() strategy.Fundable
}

func NewChargeable(strategyName string) strategy.Chargeable {
	var factory ChargeProviderFactory

	if strategyName == "paystackhosted" {
		factory = &paystackhosted.Providerfactory{}
	} else {
		panic(errors.New("invalid charge strategy provider"))
	}

	factory.Init()
	return factory.Create()
}

func NewFundable(strategyName string) strategy.Fundable {
	var factory FundProviderFactory

	if strategyName == "localwallet" {
		factory = &localwallet.ProviderFactory{}
	} else {
		panic(errors.New("invalid fund strategy provider"))
	}

	factory.Init()
	return factory.Create()
}
