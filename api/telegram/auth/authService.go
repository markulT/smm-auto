package auth

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type SendCodeAuthResult struct {
	PhoneCodeHash string
}

func SendCode(phoneNumber string) SendCodeAuthResult {
	var result SendCodeAuthResult
	client := telegram.NewClient(25826350, "6b3bb341938fc1e24dd909f8c419325f", telegram.Options{})
	err := client.Run(context.Background(), func(ctx context.Context) error {
		api := client.API()
		auth, err := api.AuthSendCode(ctx, &tg.AuthSendCodeRequest{
			PhoneNumber: phoneNumber,
			APIID:       25826350,
			APIHash:     "6b3bb341938fc1e24dd909f8c419325f",
			Settings:    tg.CodeSettings{},
		})
		result.PhoneCodeHash = auth.String()
		fmt.Println(auth)
		if err!= nil {
			fmt.Println(err)
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return result
}

func ConfirmCode(code string, phoneNumber string, phoneCodeHash string) {
	client := telegram.NewClient(25826350, "6b3bb341938fc1e24dd909f8c419325f", telegram.Options{})

	err := client.Run(context.Background(), func(ctx context.Context) error {
		api := client.API()
		fmt.Println(code)
		fmt.Println(phoneNumber)
		fmt.Println(phoneCodeHash)
		auth, err := api.AuthSignIn(ctx, &tg.AuthSignInRequest{
			PhoneNumber:       phoneNumber,
			PhoneCode:         code,
			PhoneCodeHash: 	   phoneCodeHash,
		})
		fmt.Println(auth)
		fmt.Println(err)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
