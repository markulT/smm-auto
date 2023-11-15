package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/webhook"
	"golearn/repository"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
	"io"
	"os"
)

func SetupPaymentRoutes(r *gin.Engine) {
	paymentGroup := r.Group("/payment")
	paymentGroup.GET("/plans/", jsonHelper.MakeHttpHandler(getSubPlans))
	paymentGroup.Use(auth.AuthMiddleware)
	//paymentGroup.Use(payments.PaymentMiddleware)
	//paymentGroup.POST("/intent", jsonHelper.MakeHttpHandler(createIntentHandler))
	paymentGroup.POST("/subscription", jsonHelper.MakeHttpHandler(subscriptionСreationHandler))
	paymentGroup.GET("/customerExists", customerExistsHandler)
	paymentGroup.POST("/paymentMethod/add", jsonHelper.MakeHttpHandler(addPaymentMethodHandler))
	paymentGroup.POST("/setupIntent", jsonHelper.MakeHttpHandler(initSetupIntentHandler))
	paymentGroup.POST("/setupIntent/create", jsonHelper.MakeHttpHandler(createSetupIntent))
	paymentGroup.GET("/paymentMethod/getAll", jsonHelper.MakeHttpHandler(paymentMethodsHandler))
	paymentGroup.GET("/subscriptions/plans", jsonHelper.MakeHttpHandler(getSubscriptionPlans))

	webHookGroup := r.Group("/stripeWebhook")

	webHookGroup.POST("/subscribe", subscriptionWebhookHandler)
	webHookGroup.POST("/setupIntent", jsonHelper.MakeHttpHandler(setupIntentWebhookHandler))
}

type CreateSetupIntentRequest struct{}
type CreateSetupIntentResponse struct {
	SetupClientSecret string `json:"setupClientSecret"`
	CustomerID        string `json:"customerID"`
}

func createSetupIntent(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	var body CreateSetupIntentRequest
	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "Unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	si, err := paymentsService.CreateSetupIntent(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	customer, err := paymentsService.GetCustomerByID(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{
		"setupClientSecret": si.ClientSecret,
		"customerID":        customer.ID,
	})
	return nil
}

func setupIntentWebhookHandler(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	endpointSecret := os.Getenv("checkoutSessionCompletedSecret")

	event, err := webhook.ConstructEvent(requestBody, c.GetHeader("Stripe-Signature"), endpointSecret)
	fmt.Println(event)
	fmt.Println(event.Type)
	switch event.Type {
	case "setup_intent.succeeded":
		fmt.Println(event)
	case "checkout.session.completed":
		intent, err := paymentsService.GetSetupIntent(event.Data.Object["setup_intent"].(string))
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
		user, err := mongoRepository.GetUserByEmail(event.Data.Object["customer_email"].(string))
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
		err = paymentsService.AttachPaymentMethodToCustomer(intent.PaymentMethod.ID, user.CustomerID)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
		c.JSON(200, gin.H{})
	default:
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}

	return nil
}

type InitSetupIntentRequest struct {
	PaymentMethod string `json:"paymentMethod"`
}
type InitSetupIntentResponse struct {
}

func initSetupIntentHandler(c *gin.Context) error {

	paymentService := payments.NewStripePaymentService()

	var body InitSetupIntentRequest
	jsonHelper.BindWithException(&body, c)
	userEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "Unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	err = paymentService.InitSetupIntent(user.CustomerID, body.PaymentMethod)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"status": "success"})
	return nil
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
	c.JSON(200, gin.H{"plans": planList})
	return nil
}

type AddPaymentMethodRequest struct {
	CardNumber string `json:"cardNumber"`
	ExpMonth   int64  `json:"expMonth"`
	ExpYear    int64  `json:"expYear"`
	CVC        string `json:"cvc"`
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
	c.JSON(200, gin.H{"status": "Successfully added payment method"})
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

	subscriptionID, err := paymentsService.CreateSubscription(user.CustomerID, body.SubscriptionType)
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

	c.JSON(200, gin.H{"message": "successfully added subscription"})
	return nil
}

func customerExistsHandler(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(403, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	stripeService := payments.NewStripePaymentService()
	customerExists, _ := stripeService.CustomerExists(fmt.Sprintf("%d", userEmail))
	if !customerExists {
		c.JSON(401, gin.H{"customerExists": "false"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"customerExists": "true"})
}

func paymentMethodsHandler(c *gin.Context) error {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User does not exist",
			Status: 404,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	var paymentMethods []stripe.PaymentMethod

	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(user.CustomerID),
		Type:     stripe.String("card"),
	}
	i := paymentmethod.List(params)
	for i.Next() {
		pm := i.PaymentMethod()
		paymentMethods = append(paymentMethods, *pm)
	}

	c.JSON(200, gin.H{"paymentMethods": paymentMethods})
	return nil
}

func getSubscriptionPlans(c *gin.Context) error {
	var prices []stripe.Price

	params := &stripe.PriceListParams{}

	params.Filters.AddFilter("", "", "")

	i := price.List(params)
	for i.Next() {
		p := i.Price()
		prices = append(prices, *p)
	}

	c.JSON(200, gin.H{"subLevels": prices})
	return nil
}

func subscriptionWebhookHandler(c *gin.Context) {
	c.String(200, "aboba")
}
