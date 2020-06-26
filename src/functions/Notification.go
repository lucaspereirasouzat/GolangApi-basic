package functions

import (
	"fmt"

	"github.com/NaySoftware/go-fcm"
)

const (
	serverKey = "AAAAX770V_c:APA91bGIsNe-zeBt--88AXrD-9htQ3bbizsGbqHkTi-BF2d5zAa26IfzuwOZk495CQin6fTabMZ1FYkVdaVSe7BN6PCBLrOm2Xw_L1yAE_qrVIdPw8MttCGQISx6CnBRlrdz53EfkIj-"
)

func SendNotification(tokenNotification string) {

	data := map[string]string{
		"msg": "Hello World1",
		"sum": "Happy Day",
	}

	ids := []string{
		tokenNotification,
	}

	// xds := []string{
	// 	"token5",
	// 	"token6",
	// 	"token7",
	// }

	c := fcm.NewFcmClient(serverKey)
	c.NewFcmRegIdsMsg(ids, data)
	// if len(tokens) >= 1 {
	// 	c.AppendDevices(tokens)
	// }
	//	c.AppendDevices(xds)

	status, err := c.Send()

	if err == nil {
		fmt.Println(err)
		status.PrintResults()
	} else {
		fmt.Println(err)
	}

}
