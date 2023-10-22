package payments

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golearn/repository"
)

func PaymentMiddleware(c *gin.Context) {

	stripeService := NewStripePaymentService()

	userEmail, exists := c.Get("userEmail")

	if !exists {
		c.JSON(500, gin.H{"error":"Internal server error"})
		c.Abort()
		return
	}
	user, err := repository.GetUserByEmail(fmt.Sprintf("%s",userEmail))
	customerExists, err := stripeService.CustomerExists(user.CustomerID)
	if err != nil {
		c.JSON(500, gin.H{"Error":err.Error()})
		c.Abort()
		return
	}
	if !customerExists {
		c.JSON(403, gin.H{"Error":"Customer's card is not registered"})
		c.Abort()
		return
	}
	c.Next()

}
