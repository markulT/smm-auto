package payments

import (
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/plan"
	"github.com/stripe/stripe-go/v75/subscription"
)

type SubscriptionDoesNotExistException struct {}
func (e SubscriptionDoesNotExistException) Error() string {
	return "Subscription does not exist"
}

func CreateStripeSubscription(email string, subscriptionID string) (stripeSubscriptionID string, err error) {
	if err:=checkIfPlanExists(subscriptionID);!err {
		return "",SubscriptionDoesNotExistException{}
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(email),
		Items: []*stripe.SubscriptionItemsParams{{Price: stripe.String(subscriptionID)}},
		ProrationBehavior: stripe.String("always_invoice"),
	}

	subscriptionInfo, err := subscription.New(params)
	if err != nil {
		return "",err
	}
	return subscriptionInfo.ID,nil
}

func checkIfPlanExists(planID string) bool {
	_, err := plan.Get(planID, nil)
	if err != nil {
		return false
	}
	return true
}
