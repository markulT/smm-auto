package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
	//paymentGroup.Use(payments.PaymentMiddleware)
	//paymentGroup.POST("/intent", jsonHelper.MakeHttpHandler(createIntentHandler))
	paymentGroup.POST("/subscription", jsonHelper.MakeHttpHandler(subscriptionСreationHandler))
	paymentGroup.GET("/customerExists", customerExistsHandler)
	paymentGroup.POST("/paymentMethod/add", jsonHelper.MakeHttpHandler(addPaymentMethodHandler))

	webHookGroup := r.Group("/stripeWebhook")
	webHookGroup.POST("/subscribe",subscriptionWebhookHandler)
}

func getSubPlans(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()
	planList := paymentsService.GetSubPlans()
	c.JSON(200, gin.H{"plans":planList})
	return nil
}

func addPaymentMethodHandler(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	var body struct {
		SubscriptionType string `json:"subscriptionType"`
		CardNumber string `json:"cardNumber"`
		ExpMonth int64 `json:"expMonth"`
		ExpYear int64 `json:"expYear"`
		CVC string `json:"cvc"`
	}

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

//func createIntentHandler(c *gin.Context) error {
//
//	authUserEmail, exists := c.Get("userEmail")
//	if !exists {
//		return jsonHelper.ApiError{
//			Err:    "User unauthorized",
//			Status: 401,
//		}
//	}
//	user, err := mongoRepository.GetUserByEmail(fmt.Sprintf("%s", authUserEmail))
//	if err != nil {
//		return jsonHelper.ApiError{
//			Err:    "User does not exist",
//			Status: 400,
//		}
//	}
//
//
//
//	c.JSON(200, gin.H{"message":"successfully added subscription"})
//	return nil
//}

func subscriptionСreationHandler(c *gin.Context) error {

	paymentsService := payments.NewStripePaymentService()

	var body struct {
		SubscriptionType string `json:"subscriptionType"`
	}

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