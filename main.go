package main

import (
	"github.com/gin-gonic/gin"
	"golearn/controllers"
	"golearn/utils"
)

func init() {
	utils.LoadEnvVariables()
	utils.ConnectToDb()
}



func main() {
	r := gin.Default()

	controllers.SetupAuthRoutes(r)
	controllers.SetupTelegramRoutes(r)
	r.Run()

	//client := telegram.NewClient(25826350, "6b3bb341938fc1e24dd909f8c419325f", telegram.Options{})
	//var result SendCodeAuthResult
	//err := client.Run(context.Background(), func(ctx context.Context) error {
	//	api := client.API()
	//	auth1, err := api.AuthSendCode(ctx, &tg.AuthSendCodeRequest{
	//		PhoneNumber: "380987997410",
	//		APIID:       25826350,
	//		APIHash:     "6b3bb341938fc1e24dd909f8c419325f",
	//		Settings:    tg.CodeSettings{},
	//	})
	//	result.PhoneCodeHash = auth1.String()
	//	fmt.Println(auth1)
	//	if err!= nil {
	//		fmt.Println(err)
	//		panic(err)
	//	}
	//	return nil
	//})
	//if err != nil {
	//	panic(err)
	//}
	//var inputCode string
	//var inputCodeHash string
	//fmt.Println("Input code :")
	//fmt.Scanln(&inputCode)
	//fmt.Println("Input code hash :")
	//fmt.Scanln(&inputCodeHash)
	//err1 := client.Run(context.Background(), func(ctx context.Context) error {
	//	api := client.API()
	//	authResult, err := api.AuthSignIn(ctx, &tg.AuthSignInRequest{
	//		PhoneNumber:       "380987997410",
	//		PhoneCode:         inputCode,
	//		PhoneCodeHash: 	   inputCodeHash,
	//	})
	//	fmt.Println(authResult)
	//	fmt.Println(err)
	//	return nil
	//})
	//if err1 != nil {
	//	panic(err)
	//}
}

//type SendCodeAuthResult struct {
//	PhoneCodeHash string
//}


