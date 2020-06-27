package functions

import (
	"fmt"
	"log"
	"os"

	"github.com/edganiukov/fcm"
)

func SendNotification(tokenNotification string, list []string, title string, body string) {

	fmt.Println(list)

	itemNotification := fcm.Notification{}

	itemNotification.Title = title
	itemNotification.Body = body
	//itemNotification.Subtitle = "Pereira"
	// itemNotification.Icon = "Pereira"

	msg := &fcm.Message{
		Token: tokenNotification,
		Data: map[string]interface{}{
			"foo":   "bar",
			"title": "lucas",
			"body":  "Pereira",
			"msg":   "DDD",
		},
		Notification: &itemNotification,
	}

	// Create a FCM client to send the message.
	client, err := fcm.NewClient(os.Getenv("SERVER_NOTIFICATION_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// Send the message and receive the response without retries.
	response, err := client.Send(msg)
	if err != nil {
		/* ... */
	}
	fmt.Println(response)
}
