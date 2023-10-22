package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"golearn/repository"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
)

func SetupPaymentRoutes(r *gin.Engine) {
	paymentGroup := r.Group("/payment")
	paymentGroup.GET("/plans/", jsonHelper.MakeHttpHandler(getSubPlans))
	paymentGroup.Use(auth.AuthMiddleware)
	paymentGroup.Use(payments.PaymentMiddleware)
	//paymentGroup.POST("/intent", jsonHelper.MakeHttpHandler(createIntentHandler))
	paymentGroup.POST("/subscription", jsonHelper.MakeHttpHandler(subscriptionСreationHandler))
	paymentGroup.GET("/customerExists", customerExistsHandler)
	paymentGroup.POST("/paymentMethod/add", jsonHelper.MakeHttpHandler(addPaymentMethodHandler))

	webHookGroup := r.Group("/stripeWebhook")
	webHookGroup.POST("/subscribe",subscriptionWebhookHandler)
}

type GetSubPlansResponse struct {
	Plans []*stripe.Plan `json:"plans"`
}

// @Summary Get subscription plans
// @Tags posts
// @Description Get all available subscription plans
// @ID GetSubPlans
// @Accept json
// @Produce json
// @Success 200 controllers.GetSubPlansResponse
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /payments/plans [get]
func getSubPlans(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()
	planList := paymentsService.GetSubPlans()
	c.JSON(200, gin.H{"plans":planList})
	return nil
}

type AddPaymentMethodRequest struct {
	CardNumber string `json:"cardNumber"`
	ExpMonth int64 `json:"expMonth"`
	ExpYear int64 `json:"expYear"`
	CVC string `json:"cvc"`
}

// @Summary Add payment method
// @Tags posts
// @Description Add payment method (card)
// @ID AddPaymentMethod
// @Accept json
// @Produce json
// @Param request body controllers.AddPaymentMethodRequest true "Card data"
// @Success 200 controllers.GetSubPlansResponse
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /payments/paymentMethod/add [post]
func addPaymentMethodHandler(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	var body AddPaymentMethodRequest

	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")

	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthenticated",
			Status: 401,
		}
	}

	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	pm, err := paymentsService.AddPaymentMethod(payments.CardData{
		CardNumber: body.CardNumber,
		ExpMonth:   body.ExpMonth,
		ExpYear:    body.ExpYear,
		CVC:        body.CVC,
	})
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	err = paymentsService.AttachPaymentMethodToCustomer(pm.ID, user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"status":"Successfully added payment method"})
	return nil
}

type CreateSubscriptionRequest struct {
	SubscriptionType string `json:"subscriptionType"`
}

// @Summary Create subscription
// @Tags posts
// @Description Create subscription
// @ID CreateSub
// @Accept json
// @Produce json
// @Param request body controllers.CreateSubscriptionRequest true "Card data"
// @Success 200 {string} a
// @Failure 400, 417 {object} jsonHelper.ApiError
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /payments/subscription [post]
func subscriptionСreationHandler(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	var body CreateSubscriptionRequest

	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")

	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthenticated",
			Status: 401,
		}
	}

	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	subscriptionID,err := paymentsService.CreateSubscription(user.CustomerID, body.SubscriptionType)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error : " + err.Error(),
			Status: 500,
		}
	}

	err = repository.UpdateUserSubscriptionID(fmt.Sprintf("%d", userEmail), body.SubscriptionType, subscriptionID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error : " + err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"message":"successfully added subscription"})
	return nil
}

func customerExistsHandler(c *gin.Context)  {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(403, gin.H{"error":"Invalid token"})
		c.Abort()
		return
	}
	stripeService := payments.NewStripePaymentService()
	customerExists, _ := stripeService.CustomerExists(fmt.Sprintf("%d", userEmail))
	if !customerExists {
		c.JSON(401, gin.H{"customerExists":"false"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"customerExists":"true"})
}

func subscriptionWebhookHandler(c *gin.Context) {
	c.String(200, "aboba")
}