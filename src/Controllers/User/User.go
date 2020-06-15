package user

import (
	base64 "encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	connection "docker.go/src/Connections"
	user "docker.go/src/Models/User"
	validatores "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"

	"github.com/vmihailenco/msgpack"
)

var table = "users"

/*
	Faz listagem de todos os usuarios
*/
func Index(c *gin.Context) {
	users := []user.User{}

	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)

	err = connection.QueryTable(table, page, rowsPerPage, &users)
	total, err := connection.QueryTotalTable(table)
	if err != nil {
		c.String(400, "%s", err)
		panic(err)
	}

	type IndexList struct {
		Page        uint64
		RowsPerPage uint64
		Total       uint64
		Table       []user.User
	}

	list := IndexList{page, rowsPerPage, total, users}

	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }
	c.IndentedJSON(http.StatusOK, list)
}

/*
	Cadastra um novo usuario no sistema
*/

func Store(c *gin.Context) {

	code := c.Request.FormValue("code")

	data, err := base64.StdEncoding.DecodeString(code)

	if err != nil {
		c.JSON(400, err)
		return
	}

	var user validatores.Register

	err = msgpack.Unmarshal(data, &user)

	if err != nil {
		c.String(400, "%s", err)
		panic(err)
	}
	// hasError, listError := validatores.Validate(user)
	// fmt.Println(hasError, listError)

	// if hasError {
	// 	c.JSON(400, listError)
	// 	return
	// }

	user.Password, _ = functions.GeneratePassword(user.Password)
	result, err := connection.InserIntoTable(table, []string{"UserName", "Email", "Password"}, user.Username, user.Email, user.Password)
	// db := connection.CreateConnection()
	// tx := db.MustBegin()

	// result, err := tx.Exec("INSERT INTO "+table+" (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)
	fmt.Println("result", result, err)
	// tx.Commit()

	if err != nil {
		c.String(400, "%s", err)
		panic(err)
		fmt.Println(err)
	}
	//defer db.Close()
	// fmt.Println(err)
	c.JSON(200, user)
}

/*
 Procura um novo usuario pelo id
*/
func Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)

	user := user.User{}
	err = connection.ShowRow(table, &user, "id", id)

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, user)
}

/*
 Atualiza um novo usuario pelo id
*/
func Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		c.JSON(400, err)
		return
	}

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		c.JSON(400, err)
		return
	}

	var user user.User

	err = msgpack.Unmarshal(data, &user)
	if err != nil {
		c.JSON(400, err)
		return
	}
	//	connection.InserIntoTable(table)
	db := connection.CreateConnection()
	tx := db.MustBegin()
	tx.MustExec("UPDATE "+table+"  SET username=$2 WHERE id = $1", id, user.Username)

	if err != nil {
		c.JSON(400, err)
		return
	}
	defer db.Close()

	c.JSON(200, user)
}

/*
 Deleta o usuario pelo id
*/
func Delete(c *gin.Context) {
	db := connection.CreateConnection()
	user := user.User{}

	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		c.JSON(400, err)
		return
	}
	err = db.Get(&user, "DELETE FROM "+table+"  WHERE id = $1", id)
	db.Close()

	if err != nil {
		c.JSON(400, err)
		return
	}
	fmt.Printf("%#v\n", user)

	c.JSON(200, user)
}
