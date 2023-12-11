package models

import (
	"fmt"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/product"
	"os"
)

type Subscription struct {
	ID string `bson:"_id" json:"id"`
	SubLevel int `bson:"subLevel" json:"subLevel"`
	CustomerID string `json:"customerId" bson:"customerId"`
	CurrentPeriodEnd int `json:"currentPeriodEnd" bson:"currentPeriodEnd,omitempty"`
}

func NewSubscriptionFromEventData(e *stripe.EventData) (*Subscription,error) {
	var err error

	productID, err := getProductIDFromEventData(e)
	if err != nil {
		return nil, err
	}
	subLevel, err := getSubLevelFromProductID(productID)
	if err != nil {
		return nil, err
	}
	customerID, ok := e.Object["customer"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid event object")
	}
	s := Subscription{ID: e.Object["id"].(string), SubLevel: subLevel, CustomerID: customerID}
	return &s, nil
}

func NewSubscriptionFromStripe(s *stripe.Subscription) (*Subscription, error) {
	var err error

	p, err := getProductByPriceID(s.Items.Data[0].Price.ID)
	if err != nil {
		return nil, err
	}

	subLevel, err := getSubLevelFromProductID(p.ID)
	if err != nil {
		return nil, err
	}
	subModel := Subscription{ID: s.ID, SubLevel: subLevel, CustomerID: s.Customer.ID, CurrentPeriodEnd: int(s.CurrentPeriodEnd)}

	return &subModel, nil
}

func getProductIDFromEventData(e *stripe.EventData) (string, error) {
	obj, ok := e.Object["items"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid event")
	}
	fmt.Println(obj)
	prodID, ok := obj["data"].([]interface{})[0].(map[string]interface{})["plan"].(map[string]interface{})["product"].(string)
	if !ok {
		return "", fmt.Errorf("invalid event")
	}
	fmt.Println(prodID)
	return prodID, nil
}

func getProductByPriceID(priceID string) (*stripe.Product, error) {
	p, err := price.Get(priceID, nil)
	if err != nil {
		return nil, err
	}
	return p.Product, nil
}

func getSubLevelFromProductID(prodID string) (int, error) {

	var err error
	var subLevel int
	p, err := product.Get(prodID, nil)
	if err != nil {
		return 0, err
	}

	fmt.Println(p.Name)

	switch p.Name {
	case os.Getenv("minimalSubName"):
		subLevel = 1
	case os.Getenv("standardSubName"):
		subLevel = 2
	case os.Getenv("premiumSubName"):
		subLevel = 3
	default:
		subLevel = 0
	}
	return subLevel, nil
}

