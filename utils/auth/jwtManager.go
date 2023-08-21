package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"golearn/utils"
	"os"
	"time"
)

type EmailReadingError struct {
	Message string
}
func (e EmailReadingError)Error() string {
	return e.Message
}

type Tokens struct {
	AccessToken string
	RefreshToken string
}

func createToken(body map[string]interface{}, expirationTime time.Time, secretKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range body {
		claims[key] = value
	}
	claims["exp"] = expirationTime.Unix()
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return signedToken, err
}

func CreateAccessToken(body map[string]interface{}) (string, error) {

	expirationTime := utils.GetEnvInt("refreshExpirationTimeDays", 30)
	expirationTimeUnix := time.Now().Add(time.Duration(expirationTime) * time.Hour)

	return createToken(body, expirationTimeUnix, []byte(os.Getenv("secretKey")))
}
func CreateRefreshToken(body map[string]interface{}) (string, error) {
	expirationTime := utils.GetEnvInt("refreshExpirationTimeDays", 30)
	expirationTimeUnix := time.Now().Add(time.Duration(expirationTime) * 24 * time.Hour)
	return createToken(body,expirationTimeUnix, []byte(os.Getenv("secretKeyRefresh")))
}

func Validate(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err!=nil {
		return nil,err
	}
	if _,ok := token.Claims.(jwt.Claims);!ok && !token.Valid {
		return nil, fmt.Errorf("Invalid token")
	}
	return token, nil
}

func GetSubject(tokenString string, secretKey string) (string, error)  {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err!=nil {
		return "" ,err
	}
	claims:=token.Claims.(jwt.MapClaims)
	return claims["email"].(string), nil
}