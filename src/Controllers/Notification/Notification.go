package notification

import (
	"fmt"
	"strconv"

	connection "docker.go/src/Connections"
	notification "docker.go/src/Models/Notification"
	userModels "docker.go/src/Models/User"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const table string = "notification"

/*
	Faz listagem de todos os tokens de notificação
*/
func Index(c *gin.Context) {
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)
	search := c.DefaultQuery("search", "")

	query := ""
	query = functions.SearchFields(search, []string{"username", "email", "secureLevel"})
	selectFields := functions.SelectFields([]string{})

	db := connection.CreateConnection()

	notifications := []notification.Notification{}
	err = connection.QueryTable(db, table, selectFields, page, rowsPerPage, "", &notifications)

	if err != nil {
		c.JSON(400, err)
		panic(err)
	}

	total, err := connection.QueryTotalTable(db, table, query)

	if err != nil {
		c.JSON(400, err)
		panic(err)
	}

	list := struct {
		Page        uint64
		RowsPerPage uint64
		Total       uint64
		Table       []notification.Notification
	}{page,
		rowsPerPage,
		total,
		notifications}
	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }
	defer db.Close()
	c.IndentedJSON(200, list)
}

var validate *validator.Validate

/*
	Store Cadastra um novo token de notificação no sistema
*/
func Store(c *gin.Context) {

	UserGet, _ := c.Get("auth")
	us := UserGet.(userModels.User)

	//tx := db.MustBegin()

	var notificationItem = notification.Notification{}
	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)
	fmt.Println("code", c.Request.FormValue("code"))
	if err != nil {
		c.JSON(400, err)
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
	fmt.Println(notificationItem)
	db := connection.CreateConnection()
	db.Get(&notificationItem, "INSERT INTO "+table+" (tokennotification, user_id) VALUES ($1, $2) RETURNING *", notificationItem.TokenNotification, us.ID)

	//	tx.Commit()

	defer db.Close()

	c.JSON(200, notificationItem)
}

// Show Mostra um item notificação
func Show(c *gin.Context) {
	mynotification := notification.Notification{}

	id := c.Query("id")

	db := connection.CreateConnection()
	err := db.Get(&mynotification, "SELECT * FROM "+table+" WHERE tokennotification=($1)", id)
	defer db.Close()

	if err != nil {
		c.JSON(400, err)
		return
	}
	c.JSON(200, mynotification)
}

/*
 Atualiza um novo usuario pelo id
*/
func Update(c *gin.Context) {
	id := c.Query("id")

	var notificationItem notification.Notification

	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)
	if err != nil {
		c.JSON(400, err)
		panic(err)
	}
	db := connection.CreateConnection()
	err = db.Get(&notificationItem, "UPDATE "+table+" SET tokennotification = ($2) WHERE tokennotification = ($1)", id, notificationItem.TokenNotification)
	if err != nil {
		c.JSON(400, err)
		return
	}
	defer db.Close()

	c.JSON(200, notificationItem)
}

/*
 Deleta o usuario pelo id
*/
func Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		c.JSON(400, err)
		return
	}
	db := connection.CreateConnection()
	db.MustExec("DELETE FROM "+table+" WHERE tokennotification = ($1)", id)
	defer db.Close()

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, "OK")
}

func SendNotificationToUser(c *gin.Context) {
	search := c.DefaultQuery("user_id", "1")
	query := ""
	query = " WHERE user_id = '" + search + "'"
	selectFields := functions.SelectFields([]string{"tokennotification"})

	db := connection.CreateConnection()

	notifications := []notification.Notification{}
	err := connection.QueryTable(db, table, selectFields, 0, 100, query, &notifications)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(notifications)
	// var list []string

	// for _, v := range notifications {
	// 	list = append(list, v.TokenNotification)
	// }
	// fmt.Println(list)
	//list = functions.Remove(list, 0)

	functions.SendNotification(notifications[0].TokenNotification)
}
