package messages

import (
	"context"
	"fmt"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

//func getRecipientIDAndHash() (int, int64, error) {
//	client := telegram.NewClient(25826350, "6b3bb341938fc1e24dd909f8c419325f", telegram.Options{})
//	resp, err := client.
//}

func SendMessage(text string) {
	fmt.Println("A")
	client := telegram.NewClient(25826350, "6b3bb341938fc1e24dd909f8c419325f", telegram.Options{})

	err := client.Run(context.Background(), func(ctx context.Context) error {

		api := client.API()

		//recipientID := tg.InputPeerUser{
		//	UserID: 831160444,
		//	AccessHash:
		//}

		fmt.Println("B")
		//userArr, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUser{UserID: 831160444}})

		message, err := api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
			Message: text,
			Peer: &tg.InputPeerUser{UserID: 831160444},
			SendAs: &tg.InputPeerSelf{},
		})

		fmt.Println("C")
		if err != nil {
			fmt.Println(message)
			fmt.Println(err)
			return err
		}
		fmt.Println(message)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
