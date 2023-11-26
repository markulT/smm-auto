package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golearn/models"
	"golearn/repository"
	mongoRepository "golearn/repository"
	"golearn/utils/auth"
	"golearn/utils/jsonHelper"
	"golearn/utils/payments"
)

func SetupAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	authGroup.POST("/login", jsonHelper.MakeHttpHandler(loginHandler))
	authGroup.POST("/refresh", jsonHelper.MakeHttpHandler(refresh))
	authGroup.POST("/signup", jsonHelper.MakeHttpHandler(signup))

	authGroup.Use(auth.AuthMiddleware)

	authGroup.POST("/setDeviceToken", jsonHelper.MakeHttpHandler(setDeviceTokenHandler))
	authGroup.GET("/profile", jsonHelper.MakeHttpHandler(getProfileHandler))
}

type LoginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
}

type SetDeviceTokenRequest struct {
	DeviceToken string `json:"deviceToken"`
}

type ProfileResponse struct {
	Email            string `json:"email"`
	ChannelList      []models.Channel `json:"channelList"`
	SubscriptionID   string `json:"subscriptionID"`
	SubscriptionType int    `json:"subscriptionType"`
}


// @Summary Get profile handler
// @Tags auth
// @Description Get user's profile data
// @ID getProfile
// @Accept json
// @Produce json
// @Success 200 {object} controllers.ProfileResponse
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 417 {object} jsonHelper.ApiError "Error identifying user"
// @Failure 500 {object} jsonHelper.ApiError "Internal server error"
// @Failure default {object} jsonHelper.ApiError
// @Router /archive/ [get]
func getProfileHandler(c *gin.Context) error {
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
			Err:    "User does not exist",
			Status: 404,
		}
	}
	chRepo:= mongoRepository.NewChannelRepo()
	usersChannelList, err := chRepo.FindAllByUserID(context.Background(), user.ID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    "User does not have any post",
			Status: 404,
		}
	}
	userProfile := ProfileResponse{
		Email:            user.Email,
		SubscriptionID:   user.SubscriptionID,
		SubscriptionType: user.SubscriptionType,
		ChannelList: *usersChannelList,
	}

	fmt.Println(user)

	c.JSON(200, userProfile)
	return nil
}

func setDeviceTokenHandler(c *gin.Context) error {

	var body SetDeviceTokenRequest
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
			Err:    "User does not exist",
			Status: 404,
		}
	}

	err = mongoRepository.SetUsersDeviceToken(user.ID, body.DeviceToken)

	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}

	c.JSON(200, nil)
	return nil
}

// @Summary Login
// @Tags auth
// @Description Login with email&password. Returns jwt tokens that should be saved in application. the jwt token should be pinned to each request with header (Example - "Authorization": Bearer jwtToken). If the given token is invalid - 401 status error always gets thrown
// @ID Login
// @Accept json
// @Produce json
// @Param body body LoginRequestBody true "account email"
// @Success 200 {object} controllers.AuthResponse
// @Failure 403 {object} jsonHelper.ApiError "Wrong email/password"
// @Failure 404 {object} jsonHelper.ApiError "Wrong email/password"
// @Failure default {object} jsonHelper.ApiError
// @Router /auth/login [post]
func loginHandler(c *gin.Context) error {
	var body LoginRequestBody

	jsonHelper.BindWithException(&body, c)
	userFromDB, err := repository.GetUserByEmail(body.Email)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 404,
		}
	}
	fmt.Println(userFromDB.Password)
	if err := bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(body.Password)); err != nil {
		fmt.Println(err.Error())
		return jsonHelper.ApiError{
			Err:    "Invalid password",
			Status: 403,
		}
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email":    userFromDB.Email,
		"password": userFromDB.Password,
	}, c)

	c.JSON(200, gin.H{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
	return nil
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// @Summary Refresh
// @Tags auth
// @Description Refresh jwt token
// @ID Refresh
// @Accept json
// @Produce json
// @Param body body controllers.RefreshRequest true "Account info"
// @Success 200 {object} controllers.AuthResponse
// @Failure 400 {object} jsonHelper.ApiError "Error identifying user from token"
// @Failure default {object} jsonHelper.ApiError
// @Router /auth/refresh [post]
func refresh(c *gin.Context) error {
	var body RefreshRequest
	jsonHelper.BindWithException(&body, c)
	var userFromDb models.User
	email, err := auth.GetSubject(body.RefreshToken)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	userFromDb, err = repository.GetUserByEmail(email)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	if _, err := auth.Validate(body.RefreshToken); err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 400,
		}
	}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email": userFromDb.Email,
	}, c)
	c.JSON(200, gin.H{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
	return nil
}

// @Summary Signup
// @Tags auth
// @Description Signup
// @ID Signup
// @Accept json
// @Produce json
// @Param body body controllers.LoginRequestBody true "Account info"
// @Success 200 {object} controllers.AuthResponse
// @Failure 400 {object} jsonHelper.ApiError "User with such email already exists"
// @Failure 500 {object} jsonHelper.ApiError "Internal server error (might be issue with stripe)"
// @Failure default {object} jsonHelper.ApiError
// @Router /auth/signup [post]
func signup(c *gin.Context) error {
	var body struct {
		Email    string `json:"email"`
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
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	userId, err := uuid.NewRandom()
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	newUser := models.User{ID: userId, Email: body.Email, Password: string(hashedPassword)}
	tokens := auth.GenerateTokens(map[string]interface{}{
		"email": newUser.Email,
	}, c)

	if err := repository.SaveUser(&newUser); err != nil {
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
			Status: 500,
		}
	}
	err = repository.UpdateCustomerIDByEmail(body.Email, customerID)
	if err != nil {
		return jsonHelper.ApiError{
			Err:    err.Error(),
			Status: 500,
		}
	}
	c.JSON(200, gin.H{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
	})
	return nil
}
