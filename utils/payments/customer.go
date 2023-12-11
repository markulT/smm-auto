package payments

import (
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/checkout/session"
	"github.com/stripe/stripe-go/v75/customer"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/plan"
	"github.com/stripe/stripe-go/v75/setupintent"
	"github.com/stripe/stripe-go/v75/subscription"
	"golearn/models"
)

type PaymentService interface {
	GetSubscriptionByCustomerID(cid string) ([]*models.Subscription,error)
	CreateCustomer(string) (string, error)
	CustomerExists(string) (bool, error)
	AddPaymentMethod(cd CardData) (*stripe.PaymentMethod, error)
	CreateSubscription(customerID string, priceID string) (*stripe.Subscription, error)
	AttachPaymentMethodToCustomer(pmid string, customerID string) error
	GetSubPlans() []*stripe.Plan
	InitSetupIntent(cid string, paymentMethodType string) error
	GetSetupIntent(setid string) (*stripe.SetupIntent, error)
	CreateSetupIntent(cid string) (*stripe.SetupIntent, error)
	GetCustomerByID(cid string) (*stripe.Customer, error)
	DeleteSubscriptionByCustomerID(cid string) error

	SetDefaultPaymentMethod(cID string,pmID string) error
	GetDefaultPaymentMethod(cID string) (string, error)

	CustomerSubscribed(cID string) bool
}

type PaymentRepo interface {
	DeleteSubscriptionByID(c context.Context, subID string) error
	FindSubscriptionByCustomerID(c context.Context, cID string) (*models.Subscription , error)
}

type stripePaymentService struct {
	paymentRepo PaymentRepo
}

func NewStripePaymentService(pr PaymentRepo) PaymentService {
	return &stripePaymentService{paymentRepo: pr}
}

func (s *stripePaymentService) CustomerSubscribed(cID string) bool {

	subFromDB, _ := s.paymentRepo.FindSubscriptionByCustomerID(context.Background(),cID)

	subFromStripe, _ := s.GetSubscriptionByCustomerID(cID)

	if subFromDB != nil || subFromStripe != nil {
		return true
	}
	return false
}

func (s *stripePaymentService) GetDefaultPaymentMethod(cID string) (string, error) {
	a := "invoice_settings.default_payment_method"
	fmt.Println("aboba 0")
	params := &stripe.CustomerParams{
		Expand: []*string{&a},
	}
	fmt.Println("aboba")
	c, err := customer.Get(cID, params)
	if err != nil {
		return "", err
	}
	fmt.Println(cID)
	fmt.Println(c.InvoiceSettings.DefaultPaymentMethod)
	return c.InvoiceSettings.DefaultPaymentMethod.ID, nil
}

func (s *stripePaymentService) SetDefaultPaymentMethod(cID string,pmID string) error {

	params := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String("pm_12345"),
		},
	}

	_, err := customer.Update("cus_12345", params)
	return err
}

func (s *stripePaymentService) GetSubscriptionByCustomerID(cid string) ([]*models.Subscription,error) {
	var subList []*models.Subscription

	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(cid),
	}
	i := subscription.List(params)
	for i.Next() {
		sub := i.Subscription()
		subModel, err := models.NewSubscriptionFromStripe(sub)
		if err != nil {
			return subList,err
		}
		subList = append(subList, subModel)

	}
	return subList,nil
}

func (s *stripePaymentService) DeleteSubscriptionByCustomerID(cid string) (error) {
	params := &stripe.SubscriptionListParams{
		Customer: stripe.String(cid),
	}
	i := subscription.List(params)
	for i.Next() {
		sub := i.Subscription()

		params := &stripe.SubscriptionCancelParams{}
		_, err := subscription.Cancel(sub.ID, params)
		if err != nil {
			return err
		}
		err = s.paymentRepo.DeleteSubscriptionByID(context.Background(),sub.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *stripePaymentService) CreateSubscription(customerID string, priceID string) (*stripe.Subscription, error) {
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			&stripe.SubscriptionItemsParams{
				Price: stripe.String(priceID),
			},
		},
	}
	result, err := subscription.New(params)
	return result, err
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
