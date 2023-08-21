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

func CreateAccessToken(body map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	for key, value := range body {
		claims[key] = value
	}
	expirationTime := utils.GetEnvInt("expirationTimeHours", 1)
	expirationTimeUnix := time.Now().Add(time.Duration(expirationTime) * time.Hour)
	claims["exp"] = expirationTimeUnix.Unix()

	secretKey := []byte(os.Getenv("secretKey"))
	signedToken, err := token.SignedString(secretKey)
	if err!=nil {
		return "", err
	}
	return signedToken, nil
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