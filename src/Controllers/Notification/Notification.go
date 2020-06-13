package notification

import (
	base64 "encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	connection "docker.go/src/Connections"
	notification "docker.go/src/Models/Notification"
	userModels "docker.go/src/Models/User"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/vmihailenco/msgpack"
)

/*
	Faz listagem de todos os tokens de notificação
*/
func Index(c *gin.Context) {
	notifications := []notification.Notification{}

	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)

	err = connection.QueryTable("notification", page, rowsPerPage, &notifications)

	if err != nil {
		return
	}

	type IndexList struct {
		Page        uint64
		RowsPerPage uint64
		Table       []notification.Notification
	}

	list := IndexList{page, rowsPerPage, notifications}
	fmt.Println(list)
	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }
	c.IndentedJSON(http.StatusOK, list)
}

var validate *validator.Validate

/*
	Store Cadastra um novo token de notificação no sistema
*/
func Store(c *gin.Context) {

	UserGet, _ := c.Get("auth")
	us := UserGet.(userModels.User)

	db := connection.CreateConnection()
	tx := db.MustBegin()

	var notificationItem = notification.Notification{}
	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)

	if err != nil {
		panic(err)
	}
	// data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))

	// err = msgpack.Unmarshal(data, &notificationItem)

	// if err != nil {
	// 	fmt.Println("error in conversion")
	// 	panic(err)
	// }
	//	hasError, listError := validators.Validate(notificationItem)

	// if hasError {
	// 	c.JSON(400, listError)
	// 	return
	// }

	tx.MustExec("INSERT INTO notification (tokennotification, user_id) VALUES ($1, $2)", notificationItem.TokenNotification, us.ID)

	tx.Commit()

	db.Close()

	c.JSON(200, notificationItem)
}

// /*
//  Procura um novo usuario pelo id
// */
func Show(c *gin.Context) {
	db := connection.CreateConnection()
	mynotification := notification.Notification{}

	id := c.Query("id")

	err := db.Get(&mynotification, "SELECT * FROM notification WHERE tokennotification=($1)", id)
	db.Close()

	fmt.Println(mynotification)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, mynotification)
}

/*
 Atualiza um novo usuario pelo id
*/
func Update(c *gin.Context) {
	db := connection.CreateConnection()
	//user := user.User{}

	id := c.Query("id")

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}

	var notificationItem notification.Notification

	err = msgpack.Unmarshal(data, &notificationItem)
	if err != nil {
		fmt.Println("error in conversion")
		panic(err)
	}

	err = db.Get(&notificationItem, "UPDATE notification SET tokennotification = ($2) WHERE tokennotification = ($1)", id, notificationItem.TokenNotification)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()

	fmt.Printf("%#v\n", notificationItem)

	c.JSON(200, notificationItem)
}

/*
 Deleta o usuario pelo id
*/
// func Delete(c *gin.Context) {
// 	db := connection.CreateConnection()
// 	user := user.User{}

// 	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	err = db.Get(&user, "DELETE FROM users WHERE id = $1", id)
// 	db.Close()

// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Printf("%#v\n", user)

// 	c.JSON(200, gin.H{
// 		"username": "lucas",
// 		"password": 1234,
// 		"email":    "lucas@teste.com",
// 	})
// }
