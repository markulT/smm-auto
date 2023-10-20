package payments

import (
	"github.com/stripe/stripe-go/v75/plan"
)

type SubscriptionDoesNotExistException struct {}
func (e SubscriptionDoesNotExistException) Error() string {
	return "Subscription plan does not exist"
}

func checkIfPlanExists(planID string) bool {

	_, err := plan.Get(planID, nil)
	if err != nil {
		return false
	}
	return true
}
