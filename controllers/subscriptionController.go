package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
)

func SetupPaymentRoutes(r *gin.Engine) {
	paymentGroup := r.Group("/payment")
	paymentGroup.Use(auth.AuthMiddleware)
	paymentGroup.Use(payments.PaymentMiddleware)
	paymentGroup.POST("/intent", createIntentHandler)
	paymentGroup.POST("/subscription", subscriptionСreationHandler)
	paymentGroup.GET("/customerExists", customerExistsHandler)

	webHookGroup := r.Group("/stripeWebhook")
	webHookGroup.POST("/subscribe",subscriptionWebhookHandler)
}

func createIntentHandler(c *gin.Context)  {
	c.JSON(200, gin.H{"message":"successfully added subscription"})
}

func subscriptionСreationHandler(c *gin.Context)  {

	var body struct {
		SubscriptionType string `json:"subscriptionType"`
	}

	jsonHelper.BindWithException(&body, c)

	userEmail, exists := c.Get("userEmail")

	if exists != false {
		c.JSON(401, gin.H{"error":"Invalid token extraction (auth middleware error)"})
		c.Abort()
		return
	}

	subscriptionID,err := payments.CreateStripeSubscription(fmt.Sprintf("%d", userEmail), body.SubscriptionType)
	if err != nil {
		c.JSON(401, gin.H{"error":"Invalid token extraction (auth middleware error)"})
		c.Abort()
		return
	}

	err = repository.UpdateUserSubscriptionID(fmt.Sprintf("%d", userEmail), body.SubscriptionType, subscriptionID)
	if err != nil {
		c.JSON(401, gin.H{"error":err.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message":"successfully added subscription"})
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