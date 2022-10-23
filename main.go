package main

import (
	_ "github.com/devpayments/common/contracts"
	// "github.com/devpayments/core/payments"
)

func main() {
	// get source
	// validate source
	// get destination
	// validate destination
	// check if payment can happen between destination and source // should be done on caller
	// charge source based on appropriate strategy
	// fund temp - a new wallet for each payment
	// on payment source complete fund destination
	// charge temp

	// payment := payments.NewPaymentService()
	// payment.Ser
}

//initiate payment(source, destination)
// confirm paymentsSource(source)
// fund payment wallet
// complete payment

// set charge strategy
// Set fund strategy

// register sources
// register destinations
