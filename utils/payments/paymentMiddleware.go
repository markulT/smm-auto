package payments

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func PaymentMiddleware(c *gin.Context) {
	userEmail, exists := c.Get("userEmail")

	if !exists {
		c.JSON(500, gin.H{"error":"Internal server error"})
	}
	stripeService := NewStripePaymentService()
	customerExists, err := stripeService.CustomerExists(fmt.Sprintf("%d", userEmail))
	if err != nil {
		return
	}

	if !customerExists {
		c.JSON(403, gin.H{"Error":"Customer's card is not registered"})
	}
	c.Next()

}
