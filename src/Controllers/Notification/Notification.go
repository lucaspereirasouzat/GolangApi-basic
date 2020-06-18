package notification

import (
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

	db := connection.CreateConnection()
	//tx := db.MustBegin()

	var notificationItem = notification.Notification{}
	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)

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

	db.Get(&notificationItem, "INSERT INTO "+table+" (tokennotification, user_id) VALUES ($1, $2) RETURNING *", notificationItem.TokenNotification, us.ID)

	//	tx.Commit()

	defer db.Close()

	c.JSON(200, notificationItem)
}

// Show Mostra um item notificação
func Show(c *gin.Context) {
	db := connection.CreateConnection()
	mynotification := notification.Notification{}

	id := c.Query("id")

	err := db.Get(&mynotification, "SELECT * FROM "+table+" WHERE tokennotification=($1)", id)

	if err != nil {
		c.JSON(400, err)
		return
	}
	defer db.Close()
	c.JSON(200, mynotification)
}

/*
 Atualiza um novo usuario pelo id
*/
func Update(c *gin.Context) {
	db := connection.CreateConnection()

	id := c.Query("id")

	var notificationItem notification.Notification

	err := functions.FromMSGPACK(c.Request.FormValue("code"), &notificationItem)
	if err != nil {
		c.JSON(400, err)
		panic(err)
	}

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
	db := connection.CreateConnection()

	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		c.JSON(400, err)
		return
	}
	db.MustExec("DELETE FROM "+table+" WHERE tokennotification = ($1)", id)
	defer db.Close()

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, "OK")
}
