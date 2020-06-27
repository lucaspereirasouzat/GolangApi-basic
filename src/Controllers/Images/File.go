package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const table = "file"

/*
	Faz listagem de todos os usuarios
*/
// func Index(c *gin.Context) {

// 	files := []file.File{}

// 	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
// 	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("rowsPerPage", "10"), 10, 10)
// 	search := c.DefaultQuery("search", "")

// 	query := ""
// 	query = functions.SearchFields(search, []string{"username", "email", "secureLevel"})
// 	selectFields := functions.SelectFields([]string{})

// 	db := connection.CreateConnection()

// 	err = connection.QueryTable(db, table, selectFields, page, rowsPerPage, " ", &files)
// 	total, err := connection.QueryTotalTable(db, table, query)

// 	if err != nil {
// 		c.String(400, "%s", err)
// 		panic(err)
// 	}
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer db.Close()

// 	type IndexList struct {
// 		Page        uint64
// 		RowsPerPage uint64
// 		Total       uint64
// 		Table       []file.File
// 	}

// 	list := IndexList{page, rowsPerPage, total, files}

// 	c.IndentedJSON(http.StatusOK, list)
// }

/*
	Cadastra um novo usuario no sistema
*/
var validate *validator.Validate

// func Store(c *gin.Context) {

// 	db := connection.CreateConnection()
// 	tx := db.MustBegin()
// 	fmt.Println(c.Request.FormValue("code"))

// 	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	type User struct {
// 		Username string `validate:"required"`
// 		Email    string `validate:"required,email"`
// 		Password string `validate:"required"`
// 	}

// 	var user User

// 	err = msgpack.Unmarshal(data, &user)

// 	if err != nil {
// 		fmt.Println("error in conversion")
// 		panic(err)
// 	}
// 	hasError, listError := validators.Validate(user)

// 	if hasError {
// 		c.JSON(400, listError)
// 		return
// 	}

// 	tx.MustExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, user.Password)

// 	tx.Commit()

// 	db.Close()

// 	c.JSON(200, user)
// }

/*
 Procura uma imagem pelo id
*/
func Show(c *gin.Context) {
	path := c.Query("path")
	c.File("./tmp/" + path)
}

/*
 Atualiza um novo usuario pelo id
*/
// func Update(c *gin.Context) {
// 	db := connection.CreateConnection()
// 	//user := user.User{}

// 	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
// 	if err != nil {
// 		panic(err)
// 	}

// 	var user file.File

// 	err = msgpack.Unmarshal(data, &user)
// 	if err != nil {
// 		fmt.Println("error in conversion")
// 		panic(err)
// 	}

// 	err = db.Get(&user, "UPDATE users SET username=$2, email=$3 WHERE id = $1", id, user.Username, user.Email)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	db.Close()

// 	fmt.Printf("%#v\n", user)

// 	c.JSON(200, gin.H{
// 		"username": "lucas",
// 		"password": 1234,
// 		"email":    "lucas@teste.com",
// 	})
// }

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
