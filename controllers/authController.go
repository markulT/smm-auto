package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golearn/models"
	"golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", jsonHelper.MakeHttpHandler(login))
	authGroup.POST("/refresh", jsonHelper.MakeHttpHandler(refresh))
	authGroup.POST("/signup", jsonHelper.MakeHttpHandler(signup))
}

type LoginRequestBody struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func login(c *gin.Context) error {
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	jsonHelper.BindWithException(&body, c)
	userFromDB, err := repository.GetUserByEmail(body.Email)
	if err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 404,
		}
	}
	fmt.Println(userFromDB.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(body.Password));err!=nil {
		fmt.Println(err.Error())
		return jsonHelper.ApiError{
			Err:    "Invalid password",
			Status: 403,
		}
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":userFromDB.Email,
		"password": userFromDB.Password,
	}, c)

	c.JSON(200, gin.H{
		"accessToken": tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"email": userFromDB.Email,
	})
	return nil
}


func refresh(c *gin.Context) error {
	var body struct {
		RefreshToken string `json:"refreshToken"`
	}
	jsonHelper.BindWithException(&body, c)
	var userFromDb models.User
	email, err := auth.GetSubject(body.RefreshToken)
	if err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	userFromDb, err = repository.GetUserByEmail(email)
	if err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	if _, err := auth.Validate(body.RefreshToken);err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":userFromDb.Email,
	}, c)
	c.JSON(200, gin.H{
		"accessToken":tokens.AccessToken,
		"refreshToken":tokens.RefreshToken,
		"email": userFromDb.Email,
	})
	return nil
}



func signup(c *gin.Context) error {
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	jsonHelper.BindWithException(&body, c)
	userExists := auth.UserExists(body.Email)
	if userExists {
		return jsonHelper.ApiError{
			Err:    "User already exists",
			Status: 400,
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	userId, err := uuid.NewRandom()
	if err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	newUser := models.User{ID: userId, Email: body.Email, Password: string(hashedPassword)}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":    newUser.Email,
	}, c)

	if err := repository.SaveUser(&newUser); err!=nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	stripeService := payments.NewStripePaymentService()
	customerID, err := stripeService.CreateCustomer(body.Email)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	err = repository.UpdateCustomerIDByEmail(body.Email, customerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	c.JSON(200, gin.H{
		"accessToken": tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"email":newUser.Email,
	})
	return nil
}