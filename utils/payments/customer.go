package payments

import (
	"fmt"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/plan"
	"github.com/stripe/stripe-go/v75/setupintent"
	"github.com/stripe/stripe-go/v75/subscription"
)

type PaymentService interface {
	CreateCustomer(string) (string, error)
	CustomerExists(string) (bool, error)
	CreateSubscription(email string, subscriptionID string) (subID string, err error)
	AddPaymentMethod(cd CardData) (*stripe.PaymentMethod, error)
	AttachPaymentMethodToCustomer(pmid string, customerID string) error
	GetSubPlans() []*stripe.Plan
	InitSetupIntent(cid string, paymentMethodType string) error
	GetSetupIntent(setid string) (*stripe.SetupIntent, error)
	CreateSetupIntent(cid string) (*stripe.SetupIntent, error)
	GetCustomerByID(cid string) (*stripe.Customer, error)
}

type stripePaymentService struct {
}

func NewStripePaymentService() PaymentService {
	return &stripePaymentService{}
}

func (s *stripePaymentService) GetCustomerByID(cid string) (*stripe.Customer, error) {
	c, err := customer.Get(cid, nil)
	if err != nil {
		return &stripe.Customer{}, err
	}
	return c, nil
}

func (s *stripePaymentService) CreateSetupIntent(cid string) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{
		AutomaticPaymentMethods: &stripe.SetupIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Customer: stripe.String(cid),
	}
	si, err := setupintent.New(params)
	if err != nil {
		return &stripe.SetupIntent{}, err
	}
	return si, nil
}

func (s *stripePaymentService) GetSetupIntent(setid string) (*stripe.SetupIntent, error) {
	params := &stripe.SetupIntentParams{}
	result, err := setupintent.Get(setid, params)
	if err != nil {
		return &stripe.SetupIntent{}, err
	}
	return result, nil
}

func (s *stripePaymentService) InitSetupIntent(cid string, paymentMethodType string) error {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			paymentMethodType,
		}),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSetup)),
		Customer: stripe.String(cid),
	}
	_, err := session.New(params)
	if err != nil {
		return err
	}
	return nil
}

func (s *stripePaymentService) GetSubPlans() []*stripe.Plan {
	var planList []*stripe.Plan
	params := &stripe.PlanListParams{}
	params.Filters.AddFilter("limit", "", "100")
	i := plan.List(params)
	for i.Next() {
		p := i.Plan()
		planList = append(planList, p)
	}
	return planList
}

func (s *stripePaymentService) AttachPaymentMethodToCustomer(pmid string, customerID string) error {
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	}
	_, err := paymentmethod.Attach(
		pmid,
		params,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *stripePaymentService) AddPaymentMethod(cd CardData) (*stripe.PaymentMethod, error) {

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(cd.CardNumber),
			ExpMonth: &cd.ExpMonth,
			ExpYear:  &cd.ExpYear,
			CVC:      stripe.String(cd.CVC),
		},
		Type: stripe.String("card"),
	}
	pm, err := paymentmethod.New(params)

	if err != nil {
		return &stripe.PaymentMethod{}, err
	}

	return pm, nil

}

func (s *stripePaymentService) CreateSubscription(customerID string, subscriptionID string) (stripeSubscriptionID string, err error) {
	if err := checkIfPlanExists(subscriptionID); !err {
		return "", SubscriptionDoesNotExistException{}
	}
	params := &stripe.SubscriptionParams{
		Customer:          stripe.String(customerID),
		Items:             []*stripe.SubscriptionItemsParams{{Price: stripe.String(subscriptionID)}},
		ProrationBehavior: stripe.String("always_invoice"),
	}

	subscriptionInfo, err := subscription.New(params)
	if err != nil {
		return "", err
	}
	return subscriptionInfo.ID, nil
}

func (s *stripePaymentService) CreateCustomer(email string) (string, error) {
	params := &stripe.CustomerParams{
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
	fmt.Print(email + "email")

	i := customer.List(params)

	fmt.Print(params)

	for i.Next() {
		c := i.Customer()
		if c.Email == email {
			return true, nil
		}
	}

	if err := i.Err(); err != nil {
		return false, err
	}
	return false, nil
}
