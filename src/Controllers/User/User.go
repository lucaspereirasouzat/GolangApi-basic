package user

import (
	connection "docker.go/src/Connections"
	user "docker.go/src/Models/User"
	base64 "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
	"net/http"
	"strconv"
)

func ExampleMarshal() []byte {
	type Item struct {
		Foo string
	}

	b, err := msgpack.Marshal(&Item{Foo: "lucas"})
	if err != nil {
		panic(err)
	}

	var item Item
	err = msgpack.Unmarshal(b, &item)
	if err != nil {
		panic(err)
	}
	// array bytenumber
	fmt.Println(b)

	// decodifica do base64 string msgpack
	data, err := base64.StdEncoding.DecodeString("gaNGb2+lbHVjYXM=")
	if err != nil {
		panic(err)
	}
	fmt.Printf("% x", data)

	var item2 Item
	var err2 = msgpack.Unmarshal(data, &item)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(item2)

	return b

	// Output: bar
}

func Index(c *gin.Context) {
	db := connection.CreateConnection()

	users := []user.User{}

	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 16)
	rowsPerPage, err := strconv.ParseInt(c.DefaultQuery("rowsPerPage", "10"), 10, 16)

	err = db.Select(&users, `SELECT * FROM users LIMIT ($1) OFFSET ($2)`, rowsPerPage, page*rowsPerPage)

	if err != nil {
		fmt.Println(err)
		return
	}

	type IndexList struct {
		Page        int8
		RowsPerPage int16
		Table       []user.User
	}

	list := IndexList{1, 10, users}

	b, err := msgpack.Marshal(list)
	if err != nil {
		panic(err)
	}
	c.IndentedJSON(http.StatusOK, b)
}

func Store(c *gin.Context) {

	db := connection.CreateConnection()
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", "Jason", "jmoiron@jmoiron.net", "Moiron")

	tx.Commit()
	c.JSON(200, gin.H{
		"username": "lucas",
		"password": 1234,
		"email":    "lucas@teste.com",
	})
}

func Show(c *gin.Context) {
	db := connection.CreateConnection()
	user := user.User{}
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", c.Query("id"))
	fmt.Printf("%#v\n", user)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, user)
}

func Update(c *gin.Context) {

	c.JSON(200, gin.H{
		"username": "lucas",
		"password": 1234,
		"email":    "lucas@teste.com",
	})
}

func Delete(c *gin.Context) {

	c.JSON(200, gin.H{
		"username": "lucas",
		"password": 1234,
		"email":    "lucas@teste.com",
	})
}
