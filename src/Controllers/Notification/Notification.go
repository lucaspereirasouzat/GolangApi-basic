package notification

import (
	base64 "encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	connection "docker.go/src/Connections"
	user "docker.go/src/Models/User"
	validators "docker.go/src/Validators"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/vmihailenco/msgpack"
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
	data, err := base64.StdEncoding.DecodeString("3wAAAAOodXNlcm5hbWWtbHVjYXMgUGVyZWlyYaVlbWFpbK9sdWNhc0B0ZXN0ZS5jb22ocGFzc3dvcmSkMTIzNA==")
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
}

/*
	Faz listagem de todos os usuarios
*/
func Index(c *gin.Context) {
	// type URI struct {
	// 	Page        uint64 `json:"page" `
	// 	RowsPerPage uint64 `json:"RowsPerPage"`
	// }
	// // pagaI, err := strconv.ParseUint(c.Param("page"), 10, 8)
	// // rowsPerPage1, err := strconv.ParseUint(c.Param("rowsPerPage"), 10, 8)
	// // data := Structure{
	// // 	pagaI,
	// // 	rowsPerPage1,
	// // }
	// // err = main.Validate.Struct(data)
	// // validationErrors := err.(validator.ValidationErrors)

	// // fmt.Println(validationErrors)

	// var uri URI

	// if err := c.BindQuery(&uri); err != nil {
	// 	fmt.Println(uri)
	// 	fmt.Println(err)
	// 	c.JSON(400, gin.H{"msg": err})
	// 	return
	// }
	// fmt.Println(uri)
	db := connection.CreateConnection()

	users := []user.User{}

	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)
	//fmt.Println(page, rowsPerPage)
	err = db.Select(&users, `SELECT * FROM users LIMIT ($1) OFFSET ($2)`, rowsPerPage, page*rowsPerPage)
	//fmt.Println(users)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	type IndexList struct {
		Page        uint64
		RowsPerPage uint64
		Table       []user.User
	}

	list := IndexList{page, rowsPerPage, users}

	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }
	c.IndentedJSON(http.StatusOK, list)
}

/*
	Cadastra um novo usuario no sistema
*/
var validate *validator.Validate

func Store(c *gin.Context) {

	db := connection.CreateConnection()
	tx := db.MustBegin()
	fmt.Println(c.Request.FormValue("code"))

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}

	type User struct {
		Username string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}

	var user User

	err = msgpack.Unmarshal(data, &user)

	if err != nil {
		fmt.Println("error in conversion")
		panic(err)
	}
	hasError, listError := validators.Validate(user)

	if hasError {
		c.JSON(400, listError)
		return
	}

	tx.MustExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)

	tx.Commit()

	db.Close()

	c.JSON(200, user)
}

/*
 Procura um novo usuario pelo id
*/
func Show(c *gin.Context) {
	db := connection.CreateConnection()
	user := user.User{}

	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)

	err = db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	db.Close()

	fmt.Printf("%#v\n", user)

	if err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, user)
}

/*
 Atualiza um novo usuario pelo id
*/
func Update(c *gin.Context) {
	db := connection.CreateConnection()
	//user := user.User{}

	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}

	var user user.User

	err = msgpack.Unmarshal(data, &user)
	if err != nil {
		fmt.Println("error in conversion")
		panic(err)
	}

	err = db.Get(&user, "UPDATE users SET username=$2, email=$3 WHERE id = $1", id, user.Username, user.Email)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()

	fmt.Printf("%#v\n", user)

	c.JSON(200, gin.H{
		"username": "lucas",
		"password": 1234,
		"email":    "lucas@teste.com",
	})
}

/*
 Deleta o usuario pelo id
*/
func Delete(c *gin.Context) {
	db := connection.CreateConnection()
	user := user.User{}

	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = db.Get(&user, "DELETE FROM users WHERE id = $1", id)
	db.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", user)

	c.JSON(200, gin.H{
		"username": "lucas",
		"password": 1234,
		"email":    "lucas@teste.com",
	})
}
