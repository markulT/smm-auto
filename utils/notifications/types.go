package notifications

import (
	"context"
	"firebase.google.com/go/messaging"
)

type NotificationService interface {
	SendNotification(title, content string, token string) error

}

type DefaultNotificationService struct {
	fcmClient *messaging.Client
}

//func NewDefaultNotificationService() NotificationService {
//	return &DefaultNotificationService{}
//}

func (ns *DefaultNotificationService) SendNotification(title, content string, token string) error {

	_, err := ns.fcmClient.Send(context.TODO(),&messaging.Message{
		Notification: &messaging.Notification{
			Title:    title,
			Body:     content,
		},
		Token: token,
	})

	if err != nil {
		return err
	}

	return nil
}
