package user

import (
	base64 "encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"

	connection "docker.go/src/Connections"
	userModels "docker.go/src/Models/User"
	validatores "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"

	"github.com/vmihailenco/msgpack"
)

const table = "users"

// Index Faz listagem de todos os usuarios
func Index(c *gin.Context) {
	// Pega a pagina e a quantidade de campos que serão exibidos
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("RowsPerPage", "10"), 10, 10)
	search := c.DefaultQuery("search", "")

	query := functions.SearchFields(search, []string{"username", "email", "secureLevel"})
	selectFields := functions.SelectFields([]string{"id", "username", "email", "securelevel", "created_at"})

	db := connection.CreateConnection()

	users := []userModels.User{}
	//Faz a query principal que retorna os usuarios com paginação
	err = connection.QueryTable(
		db,
		table,
		selectFields,
		page,
		rowsPerPage,
		query,
		&users)
	//Faz a query que retorna o total de usuarios
	total, err := connection.QueryTotalTable(db, table, query)

	defer db.Close()

	if err != nil {
		c.String(400, "%s", err)
		panic(err)
	}

	type IndexList struct {
		Page        uint64
		RowsPerPage uint64
		Total       uint64
		Table       []userModels.User
	}

	list := IndexList{page, rowsPerPage, total, users}

	//	b := functions.ToMSGPACK(list)
	// b, err := msgpack.Marshal(list)
	// if err != nil {
	// 	panic(err)
	// }
	c.JSON(http.StatusOK, list)
}

// Store Cadastra um novo usuario no sistema
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

	user.Password = functions.GenerateMD5(user.Password)
	db := connection.CreateConnection()

	result, err := connection.InserIntoTable(db, table, []string{"UserName", "Email", "Password"}, user.Username, user.Email, user.Password)

	fmt.Println("result", result, err)
	// tx.Commit()

	if err != nil {
		c.String(400, "%s", err)
		panic(err)
		fmt.Println(err)
	}
	defer db.Close()
	// fmt.Println(err)
	c.JSON(200, user)
}

// Show Procura um usuario pelo id
func Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)

	user := userModels.User{}
	db := connection.CreateConnection()

	err = connection.ShowRow(db, table, &user, "id", id)
	defer db.Close()
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, user)
}

// Update o usuario pelo id
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

	var userMsgPack userModels.User

	err = msgpack.Unmarshal(data, &userMsgPack)
	if err != nil {
		c.JSON(400, err)
		return
	}

	var fullUser userModels.User
	db := connection.CreateConnection()

	connection.ShowRow(db, table, &fullUser, "id", id)

	file, _, err := c.Request.FormFile("upload")
	userid := strconv.Itoa(int(id))
	filepath := "./tmp/userfile_" + userid + ".png"
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("full", fullUser)
	if !fullUser.FileId.Valid {
		result, err := connection.InserIntoTable(db, "file", []string{"path", "userid"}, filepath, userid)
		fmt.Println(result)
		if err != nil {
			fmt.Println(err)
			return
		}

		type FileResult struct {
			ID uint64
		}

		var fileres FileResult

		err = connection.ShowRow(db, "file", &fileres, "path", filepath)
		if err != nil {
			c.JSON(400, err)
			return
		}
		result, err = connection.UpdateRow(db, table, []string{"username", "file_id"}, "id", fullUser.ID, userMsgPack.Username, fileres.ID)
		if err != nil {
			c.JSON(400, err)
			return
		}
		fmt.Printf("%#v\n", userMsgPack)

		c.JSON(200, userMsgPack)
	} else {
		fmt.Printf("%#v\n", "errasd")
		connection.UpdateRow(db, table, []string{"username"}, "id", fullUser.ID, userMsgPack.Username)
		c.JSON(200, userMsgPack)

		if err != nil {
			fmt.Println("err", err)
			return
		}
	}
	defer db.Close()
	// connection.UpdateRow(table, []string{"username"}, "ID", id, user.Username)

	c.JSON(200, userMsgPack)
}

//Delete o usuario pelo id
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
