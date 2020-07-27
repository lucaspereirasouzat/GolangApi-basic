package notification

import (
	"fmt"
	"strconv"

	connection "docker.go/src/Connections"
	models "docker.go/src/Models"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const table string = "notification"

/*
	Faz listagem de todos os tokens de notificação
*/
func Index(c *gin.Context) {

	/* dados do usuario */
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)
	search := c.DefaultQuery("search", "")
	/* fim dados do usuario */

	/* Seleção de campos */
	query := functions.SearchFields(search, []string{"username", "email", "secureLevel"})
	selectFields := functions.SelectFields([]string{})
	/* fim Seleção de campos */

	notifications := []models.Notification{}

	db := connection.CreateConnection()
	err = connection.QueryTable(db, table, selectFields, page, rowsPerPage, "", &notifications)

	if err != nil {
		c.JSON(400, err)
		panic(err)
	}

	total, err := connection.QueryTotalTable(db, table, query)
	defer db.Close()

	if err != nil {
		c.JSON(400, err)
		panic(err)
	}

	list := struct {
		Page        uint64
		RowsPerPage uint64
		Total       uint64
		Table       []models.Notification
	}{page,
		rowsPerPage,
		total,
		notifications}
	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }

	// resposta para o usuario
	c.IndentedJSON(200, list)
}

var validate *validator.Validate

/*
	Store Cadastra um novo token de notificação no sistema
*/
func Store(c *gin.Context) {

	UserGet, _ := c.Get("auth")
	us := UserGet.(models.User)

	//tx := db.MustBegin()

	var notificationItem = models.Notification{}
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
	mynotification := models.Notification{}

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

	var notificationItem models.Notification

	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)
	if err != nil {
		c.JSON(400, err)
		panic(err)
	}
	db := connection.CreateConnection()
	err = db.Get(&notificationItem, "UPDATE "+table+" SET tokennotification = ($2) WHERE tokennotification = ($1)", id, notificationItem.TokenNotification)
	defer db.Close()
	if err != nil {
		c.JSON(400, err)
		return
	}

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

// SendNotificationToUser envia notificação para um usuario
func SendNotificationToUser(c *gin.Context) {
	search := c.DefaultQuery("user_id", "1")
	query := " WHERE user_id = '" + search + "'"
	selectFields := functions.SelectFields([]string{"tokennotification"})

	var notifications []models.Notification
	db := connection.CreateConnection()
	err := connection.QueryTable(db, table, selectFields, 0, 100, query, &notifications)
	defer db.Close()
	if err != nil {
		fmt.Println("error query", err)
	}
	//fmt.Println(notifications)
	var list []string

	for _, v := range notifications {
		list = append(list, v.TokenNotification)
	}
	fmt.Println(list)
	//list = functions.Remove(list, 0)
	go func() {
		title := c.Query("title")
		body := c.Query("body")

		for _, element := range list {
			err := functions.SendNotification(element, list, title, body)
			if err == nil {
				db := connection.CreateConnection()
				db.MustExec("DELETE FROM "+table+" WHERE tokennotification = ($1)", element)
				defer db.Close()
			}
			// index is the index where we are
			// element is the element from someSlice for where we are
		}
		//functions.SendNotification(notifications[0].TokenNotification, list, title, body)

		db := connection.CreateConnection()
		db.MustExec("INSERT INTO datanotification (title,body, user_id) VALUES ($1, $2,$3)", title, body, search)
		defer db.Close()
	}()

	c.String(200, "enviado para ")
}

func MyNotifications(c *gin.Context) {
	search := c.DefaultQuery("user_id", "1")
	query := " WHERE user_id = '" + search + "'"

	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("RowsPerPage", "50"), 10, 10)
	selectFields := functions.SelectFields([]string{})
	db := connection.CreateConnection()

	notifications := []models.DataNotification{}
	err = connection.QueryTable(db, "datanotification", selectFields, page, rowsPerPage, query+" ORDER BY created_at DESC", &notifications)
	total, err := connection.QueryTotalTable(db, "datanotification", query)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(notifications)

	list := struct {
		Page        uint64
		RowsPerPage uint64
		Total       uint64
		Table       []models.DataNotification
	}{page,
		rowsPerPage,
		total,
		notifications}

	c.JSON(200, list)
	// var list []string

	// for _, v := range notifications {
	// 	list = append(list, v.TokenNotification)
	// }
	// fmt.Println(list)
	//list = functions.Remove(list, 0)

	//functions.SendNotification(notifications[0].TokenNotification)
}
