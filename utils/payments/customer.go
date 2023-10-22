package payments

import (
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/plan"
	"github.com/stripe/stripe-go/v75/subscription"
)

type PaymentService interface {
	CreateCustomer(string) (string, error)
	CustomerExists(string) (bool, error)
	CreateSubscription(email string, subscriptionID string) (subID string, err error)
	AddPaymentMethod(cd CardData) (*stripe.PaymentMethod,error)
	AttachPaymentMethodToCustomer(pmid string, customerID string) error
	GetSubPlans() []*stripe.Plan
}

type stripePaymentService struct {
}

func NewStripePaymentService() PaymentService {
	return &stripePaymentService{}
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

func (s *stripePaymentService) AddPaymentMethod(cd CardData) (*stripe.PaymentMethod,error) {

	params := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number: stripe.String(cd.CardNumber),
			ExpMonth: &cd.ExpMonth,
			ExpYear: &cd.ExpYear,
			CVC: stripe.String(cd.CVC),
		},
		Type: stripe.String("card"),
	}
	pm, err := paymentmethod.New(params)

	if err != nil {
		return &stripe.PaymentMethod{},err
	}

	return pm, nil

}

func (s *stripePaymentService) CreateSubscription(customerID string, subscriptionID string) (stripeSubscriptionID string, err error) {
	if err:=checkIfPlanExists(subscriptionID);!err {
		return "",SubscriptionDoesNotExistException{}
	}
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{{Price: stripe.String(subscriptionID)}},
		ProrationBehavior: stripe.String("always_invoice"),
	}

	subscriptionInfo, err := subscription.New(params)
	if err != nil {
		return "",err
	}
	return subscriptionInfo.ID,nil
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
		if c.Email == email {
			return true, nil
		}
	}

	if err := i.Err(); err != nil {
		return false, err
	}
	return false, nil
}