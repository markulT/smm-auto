package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/paymentmethod"
	"github.com/stripe/stripe-go/v75/price"
	"github.com/stripe/stripe-go/v75/product"
	"github.com/stripe/stripe-go/v75/webhook"
	"golearn/models"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
	"io"
	"os"
)

type subscriptionController struct {
	paymentService payments.PaymentService
	paymentRepo mongoRepository.PaymentRepo
}

func SetupPaymentRoutes(r *gin.Engine, ps payments.PaymentService, pr mongoRepository.PaymentRepo) {
	//paymentRepo := mongoRepository.NewPaymentRepo()
	//paymentService := payments.NewStripePaymentService(paymentRepo)
	sc := subscriptionController{paymentService: ps, paymentRepo: pr}
	paymentGroup := r.Group("/payment")
	paymentGroup.GET("/subscriptions/plans", jsonHelper.MakeHttpHandler(sc.getSubscriptionPlans))
	paymentGroup.Use(auth.AuthMiddleware)
	//paymentGroup.Use(payments.PaymentMiddleware)
	//paymentGroup.POST("/intent", jsonHelper.MakeHttpHandler(createIntentHandler))
	paymentGroup.GET("/subscription", jsonHelper.MakeHttpHandler(sc.getSubscriptionInfo))
	paymentGroup.POST("/subscription", jsonHelper.MakeHttpHandler(sc.subscriptionСreationHandler))
	paymentGroup.POST("/paymentMethod/setDefault", jsonHelper.MakeHttpHandler(sc.setDefaultPaymentMethod))
	paymentGroup.DELETE("/subscription/cancel", jsonHelper.MakeHttpHandler(sc.deleteSubscription))
	paymentGroup.GET("/customerExists", sc.customerExistsHandler)
	paymentGroup.POST("/paymentMethod/add", jsonHelper.MakeHttpHandler(sc.addPaymentMethodHandler))
	paymentGroup.POST("/setupIntent", jsonHelper.MakeHttpHandler(initSetupIntentHandler))
	paymentGroup.POST("/setupIntent/create", jsonHelper.MakeHttpHandler(sc.createSetupIntent))
	paymentGroup.GET("/paymentMethod/getAll", jsonHelper.MakeHttpHandler(sc.paymentMethodsHandler))
	paymentGroup.GET("/paymentMethod/getDefault", jsonHelper.MakeHttpHandler(sc.getDefaultPaymentMethod))

	webHookGroup := r.Group("/stripeWebhook")

	webHookGroup.POST("/subscribe", jsonHelper.MakeHttpHandler(subscriptionWebhookHandler))
	webHookGroup.POST("/setupIntent", jsonHelper.MakeHttpHandler(setupIntentWebhookHandler))
}

type CreateSetupIntentRequest struct{}
type CreateSetupIntentResponse struct {
}

type SetDefaultPaymentMethodRequest struct {
	PaymentMethodID string `json:"paymentMethodId"`
}

func (sc *subscriptionController) getDefaultPaymentMethod(c *gin.Context) error {

	var err error

	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	paymentMethodID, err := sc.paymentService.GetDefaultPaymentMethod(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"paymentMethodId":paymentMethodID})
	return nil
}

func (sc *subscriptionController) setDefaultPaymentMethod(c *gin.Context) error {
	var err error
	var body SetDefaultPaymentMethodRequest
	jsonHelper.BindWithException(&body, c)
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	err = sc.paymentService.SetDefaultPaymentMethod(user.CustomerID, body.PaymentMethodID)
	c.JSON(200, gin.H{})
	return nil
}

func (sc *subscriptionController) getSubscriptionInfo(c *gin.Context)  error {

	fmt.Println("nigga")

	var err error
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	sub,err := sc.paymentService.GetSubscriptionByCustomerID(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{"subscriptions":sub})
	return nil
}

func (sc *subscriptionController) deleteSubscription(c *gin.Context) error {
	var err error
	authUserEmail, exists := c.Get("userEmail")
	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthorized",
			Status: 417,
		}
	}
	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 417,
		}
	}

	err = sc.paymentService.DeleteSubscriptionByCustomerID(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{})
	return nil
}

func (sc *subscriptionController) createSetupIntent(c *gin.Context) error {

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

	si, err := sc.paymentService.CreateSetupIntent(user.CustomerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	customer, err := sc.paymentService.GetCustomerByID(user.CustomerID)
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
	paymentRepo := mongoRepository.NewPaymentRepo()
	paymentsService := payments.NewStripePaymentService(paymentRepo)

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	endpointSecret := os.Getenv("checkoutSessionCompletedSecret")

	event, err := webhook.ConstructEvent(requestBody, c.GetHeader("Stripe-Signature"), endpointSecret)

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
	paymentRepo := mongoRepository.NewPaymentRepo()
	paymentService := payments.NewStripePaymentService(paymentRepo)

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

type AddPaymentMethodRequest struct {
	CardNumber string `json:"cardNumber"`
	ExpMonth   int64  `json:"expMonth"`
	ExpYear    int64  `json:"expYear"`
	CVC        string `json:"cvc"`
}

// @Summary Add payment method
// @Tags payments
// @Description Add payment method (card)
// @ID AddPaymentMethod
// @Accept json
// @Produce json
// @Param request body controllers.AddPaymentMethodRequest true "Card data"
// @Success 200 {string} string "Returns array of subscription plans (subscription type can be checked here : https://stripe.com/docs/api/plans)"
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /payments/paymentMethod/add [post]
func (sc *subscriptionController) addPaymentMethodHandler(c *gin.Context) error {

	var body AddPaymentMethodRequest

	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")

	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthenticated",
			Status: 400,
		}
	}

	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", userEmail))
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	pm, err := sc.paymentService.AddPaymentMethod(payments.CardData{
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
	err = sc.paymentService.AttachPaymentMethodToCustomer(pm.ID, user.CustomerID)
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
	PriceID string `json:"priceId"`
}

// @Summary Create subscription
// @Tags payments
// @Description Create subscription
// @ID CreateSub
// @Accept json
// @Produce json
// @Param request body controllers.CreateSubscriptionRequest true "Card data"
// @Success 200 {string} string "OK"
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError
// @Failure default {object} jsonHelper.ApiError
// @Router /payments/subscription [post]
func (sc *subscriptionController) subscriptionСreationHandler(c *gin.Context) error {

	var body CreateSubscriptionRequest

	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")

	if !exists {
		return jsonHelper.ApiError{
			Err:    "User unauthenticated",
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

	subExists := sc.paymentService.CustomerSubscribed(user.CustomerID)
	if subExists {
		return jsonHelper.ApiError{
			Err:    "User already subscribed",
			Status: 400,
		}
	}

	subscription, err := sc.paymentService.CreateSubscription(user.CustomerID, body.PriceID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error : " + err.Error(),
			Status: 500,
		}
	}
	subModel, err := models.NewSubscriptionFromStripe(subscription)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	fmt.Println(subModel)
	err = sc.paymentRepo.SaveSubscription(context.Background(), *subModel)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "Internal server error : " + err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, gin.H{"message": "successfully added subscription"})
	return nil
}

func (sc *subscriptionController) customerExistsHandler(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		c.JSON(403, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	customerExists, _ := sc.paymentService.CustomerExists(fmt.Sprintf("%d", userEmail))
	if !customerExists {
		c.JSON(404, gin.H{"customerExists": "false"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"customerExists": "true"})
}

func (sc *subscriptionController) paymentMethodsHandler(c *gin.Context) error {
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

func (sc *subscriptionController) getSubscriptionPlans(c *gin.Context) error {
	var prices []stripe.Product

	params := &stripe.ProductListParams{}

	i:=product.List(params)
	for i.Next() {
		p := i.Product()
		priceParams := &stripe.PriceParams{}
		fullPrice, _ := price.Get(i.Product().DefaultPrice.ID, priceParams)
		p.DefaultPrice = fullPrice
		prices = append(prices, *p)
	}

	c.JSON(200, gin.H{"subLevels": prices})
	return nil
}

func subscriptionWebhookHandler(c *gin.Context) error {
	var err error
	//var body models.Subscription
	//jsonHelper.BindWithException(&body, c)

	paymentRepo := mongoRepository.NewPaymentRepo()

	//err = paymentRepo.SaveSubscription(context.Background(), body)
	//if err != nil {
	//	return jsonHelper.ApiError{
	//		Err:    "Error saving subscription",
	//		Status: 500,
	//	}
	//}

	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	event, err := webhook.ConstructEvent(requestBody, c.GetHeader("Stripe-Signature"), "whsec_OfBPfcD0lNo0PNqYdOQmOdrlsBcLD8Gt")
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	//jsonData, err := json.MarshalIndent(event.Data, "", "  ")
	//if err != nil {
	//	return jsonHelper.ApiError{
	//		Err:    "Aboba",
	//		Status: 500,
	//	}
	//}

	switch event.Type {
	case "customer.subscription.created":
		// Then define and call a function to handle the event customer.subscription.created
		subscription, err := models.NewSubscriptionFromEventData(event.Data)
		if err != nil {
			return jsonHelper.ApiError{
				Err:    err.Error(),
				Status: 500,
			}
		}
		err = paymentRepo.SaveSubscription(context.Background(), *subscription)
		//fmt.Println()
	case "customer.subscription.deleted":
		// Then define and call a function to handle the event customer.subscription.deleted
	// ... handle other event types
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	c.String(200, "aboba")
	return nil
}
