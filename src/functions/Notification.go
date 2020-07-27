package functions

import (
	"fmt"
	"os"

	"github.com/NaySoftware/go-fcm"
	//"github.com/edganiukov/fcm"
)

func SendNotification(tokenNotification string, list []string, title string, body string) error {

	//fmt.Println("list", list)

	// for i := 0; i < count(list); i++ {

	// }

	itemNotification := fcm.NotificationPayload{}

	itemNotification.Title = title
	itemNotification.Body = body
	//itemNotification.Subtitle = "Pereira"
	// itemNotification.Icon = "Pereira"

	msg := &fcm.FcmMsg{
		// Token: tokenNotification,
		//RegistrationIDs: list,
		//	To: tokenNotification,
		Data: map[string]interface{}{
			"foo":   "bar",
			"title": "lucas",
			"body":  "Pereira",
			"msg":   "DDD",
		},

		// Notification: &itemNotification,
	}

	ids := []string{
		tokenNotification,
	}

	// Create a FCM client to send the message.
	client := fcm.NewFcmClient(os.Getenv("SERVER_NOTIFICATION_KEY"))
	client.NewFcmRegIdsMsg(ids, msg)
	// if err != nil {
	// 	return err
	// }

	// Send the message and receive the response without retries.
	_, err := client.Send()
	if err != nil {
		fmt.Println("Error on send", err)
		return err
		/* ... */
	}
	return nil
}
