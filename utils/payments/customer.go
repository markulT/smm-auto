package payments

import (
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/customer"
)

type PaymentService interface {
	CreateCustomer(string) (string, error)
	CustomerExists(string) (bool, error)
}

type stripePaymentService struct {
}

func NewStripePaymentService() PaymentService {
	return &stripePaymentService{}
}

func (s *stripePaymentService) CreateCustomer(email string) (string, error) {
	params:=&stripe.CustomerParams{
		Email: &email,
	}
	c, err := customer.New(params)
	if err != nil {
		return "", err
	}
	return c.ID, err
}

func (s *stripePaymentService) CustomerExists(email string) (bool, error) {
	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}

	i := customer.List(params)

	for i.Next() {
		c := i.Customer()

		// If a customer with the specified email is found, return true
		if c.Email == email {
			return true, nil
		}
	}

	if err := i.Err(); err != nil {
		return false, err
	}

	// Customer with the specified email not found
	return false, nil
}