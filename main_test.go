package main

import (
	"errors"
	"testing"

	goStripe "github.com/stripe/stripe-go"
)

type fakeIdeal struct {
	hasError bool
}

func (stripeObject fakeIdeal) newCharge(params *goStripe.ChargeParams) (*goStripe.Charge, error) {
	if stripeObject.hasError == true {
		return &goStripe.Charge{}, errors.New("stub_error")
	}
	return &goStripe.Charge{}, nil
}

func (stripeObject fakeIdeal) newChargeParams(amount int64, currency string, sourceID string) *goStripe.ChargeParams {
	return &goStripe.ChargeParams{}
}

func TestChargeSourceError(t *testing.T) {
	testCases := []struct {
		Name     string
		HasError bool
		Validate func(t *testing.T, err error)
	}{
		{
			Name:     "Return error nil if stub return nil",
			HasError: false,
			Validate: func(t *testing.T, err error) {
				if err != nil {
					t.Fatal("Test fail because error is not nil")
				}
			},
		},
		{
			Name:     "Return error if stub returns it",
			HasError: true,
			Validate: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("Test fail because error is nil")
				}
			},
		},
	}
	for _, test := range testCases {
		stub := fakeIdeal{test.HasError}
		err := ChargeSource(stub, 1, "EUR", "123456789")
		test.Validate(t, err)
	}
}
