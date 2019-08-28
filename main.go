package main

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

type StripeInterface interface {
	newCharge(params *stripe.ChargeParams) (*stripe.Charge, error)
	newChargeParams(amount int64, currency string, sourceID string) *stripe.ChargeParams
}

type StripeStruct struct{}

func (stripeObject StripeStruct) newCharge(params *stripe.ChargeParams) (*stripe.Charge, error) {
	chargeResponse, err := charge.New(params)
	return chargeResponse, err
}

func (stripeObject StripeStruct) newChargeParams(amount int64, currency string, sourceID string) *stripe.ChargeParams {
	chargeParams := &stripe.ChargeParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	chargeParams.SetSource(sourceID)
	return chargeParams
}

/* PUBLIC INTERFACE */

// ChargeSource create a new stripe charge
func ChargeSource(lib StripeInterface, amount int64, currency string, sourceID string) error {
	chargeParams := lib.newChargeParams(amount, currency, sourceID)
	_, err := lib.newCharge(chargeParams)
	return err
}
