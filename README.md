# how to mock

As a javascript developer, I was used to just overwrite the functions. Things like `a = b` make sense if I want to overwrite `a`. Packages like _sinon_ were really helpful.

I started developing in GO a few months ago. Way more different, is like a new world. I started working on a project, and you know, things like TDD should not depend on technology, they should be "agnostic". So, I start trying to make things in TDD way.

 When I started touching other developers code if found the called DB everywhere, and not just DB but external modules. So I thought, ok, let's do it on a Javascript way, and then boom! it's not possible. Please read that when I say "it's not possible" I mean **NOT POSSIBLE**. You CAN NOT say something like.
 
 ## Code
 ```go
package main
import "net/http"

func main(){
  http.get("google.com")
}
 ```
## Test 
 ```go
package main
import "testing"

func mock(){...}

func TestMain(t *testing.T){
  http.get = mock
}
 ```

 So I went into an existential crisis, I couldn't believe that GO, a brand new google language doesn't has the mocking system. But it actually has one. But you have to think in a testing way, have a good architecture and finally apply some design patterns.

For this example, I will use a _stripe_ library as an external library

```go
package main

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

/* PUBLIC INTERFACE */

// ChargeSource create a new stripe charge
func ChargeSource(lib StripeInterface, amount int64, currency string, sourceID string) error {
  chargeParams := &stripe.ChargeParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}
	chargeParams.SetSource(sourceID)
	chargeParams := lib.newChargeParams(amount, currency, sourceID)
	chargeResponse, err := charge.New(params)
	return err
}
```

Now imagine calling the `ChargeSource` Function, every time this function will call Stripe library, so, if you run tests a lot, maybe you get blocked and test will stop working, but more importantly, when you make TDD you just need to test chucks of code, no the whole behavior, for that, use integration tests. So, having that in mind, that you just need to test YOUR function, not Stripe code we have to make some changes.

First of all, you have to use Interfaces, in this way we can deceive the go package. Second, we need to use some kind of dependency injection.

Please read the code starting from below 

```go
// 4
type StripeStruct struct{}
// 3
func (stripeObject StripeStruct) newCharge(params *stripe.ChargeParams) (*stripe.Charge, error) {
	chargeResponse, err := charge.New(params)
	return chargeResponse, err
}
// 2
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
// 1
func ChargeSource(lib StripeInterface, amount int64, currency string, sourceID string) error {
	chargeParams := lib.newChargeParams(amount, currency, sourceID)
	_, err := lib.newCharge(chargeParams)
	return err
}
``` 

The first thing is to create a new public interface (Upper case) function which for any reason **ANY** should call directly third party modules. This function should call the `2 and 3` functions we wrote above.
Here is where you have to think about architecture for a second, make a clean architecture. Go back to your college years and remember **GRASP**. Do you remember the High cohesion/Low coupling? if you don't shame on you, stop this reading and google it.

As you can see, the `ChargeSource` function gets a `StripeInterface` argument. This is the interface we are about to implement, and this is the Dependency injection I mentioned before. Now our function will call the functions we want to be called if they are of type `StripeInterface`

Now, we have a low coupling code. If we want to use a different library, not Stripe, the external interface won't be affected at all.

But we need to do something else, we have to create an interface 
```go
type StripeInterface interface {
	newCharge(params *stripe.ChargeParams) (*stripe.Charge, error)
	newChargeParams(amount int64, currency string, sourceID string) *stripe.ChargeParams
}
```

This interface will allow me to use a fake struct on the test and make it looks like `StripeStruct`

## Test:


First 
``` go
type fakeIdeal struct {
	hasError bool
}
```
This will be our fake struct

These are our stubs:
```go
func (stripeObject fakeIdeal) newCharge(params *goStripe.ChargeParams) (*goStripe.Charge, error) {
	if stripeObject.hasError == true {
		return &goStripe.Charge{}, errors.New("stub_error")
	}
	return &goStripe.Charge{}, nil
}

func (stripeObject fakeIdeal) newChargeParams(amount int64, currency string, sourceID string) *goStripe.ChargeParams {
	return &goStripe.ChargeParams{}
}
```

Note that when we use the receiver we make that receiver part of interface `StripeInterface` The one our method is expected.

The final code looks like this.
```go
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

``` 