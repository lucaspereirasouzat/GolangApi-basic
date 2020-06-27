package user

import (
	base64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strconv"

	connection "docker.go/src/Connections"
	models "docker.go/src/Models"
	validatores "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"

	"github.com/vmihailenco/msgpack"
)

type Lista struct {
	Page        uint64
	RowsPerPage uint64
	Total       uint64
	Table       []models.User
}

func callTable(page uint64, rowsPerPage uint64, search string) (list Lista, err error) {
	query := functions.SearchFields(search, []string{"username", "email", "secureLevel"})
	selectFields := functions.SelectFields([]string{"id", "username", "email", "securelevel", "created_at"})

	db := connection.CreateConnection()

	users := []models.User{}
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
	total, err := connection.QueryTotalTable(
		db,
		table,
		query)

	defer db.Close()

	// if err != nil {
	// 	c.JSON(400, err)
	// 	panic(err)
	// }

	//models := rmfield(models.User, "Password")

	//var list Lista

	list.Page = page
	list.RowsPerPage = rowsPerPage
	list.Total = total
	list.Table = users

	return list, nil
}

const table string = "users"

// Index Faz listagem de todos os usuarios
func Index(c *gin.Context) {
	// Pega a pagina e a quantidade de campos que serão exibidos
	page, _ := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, _ := strconv.ParseUint(c.DefaultQuery("RowsPerPage", "50"), 10, 10)
	search := c.DefaultQuery("search", "")
	var list Lista
	//connection.SetItemRedis()
	if page == 0 && rowsPerPage == 50 && search == "" {
		result, err := connection.GetItemRedis("listUsers")
		if err != nil {
			list, err = callTable(page, rowsPerPage, search)

			if err != nil {
				c.JSON(400, err)
				panic(err)
			}

			go func() {
				json, err := json.Marshal(list)
				if err != nil {
					//return nil, err
				}
				newJson := string(json)
				connection.SetItemRedis("listUsers", newJson)
			}()

		} else {
			err := json.Unmarshal([]byte(result), &list)
			if err != nil {
				c.JSON(400, err)
				panic(err)
			}
		}
	} else {
		var err error
		list, err = callTable(page, rowsPerPage, search)

		if err != nil {
			c.JSON(400, err)
			panic(err)
		}
	}

	//fmt.Println(list)
	c.JSON(200, list)
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
		c.JSON(400, err)
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

	err = db.Get(&user, "INSERT INTO "+table+" (username,email,password) VALUES ($1,$2,$3)  RETURNING *", user.Username, user.Email, user.Password)

	if err != nil {
		c.String(400, "%s", err)
		panic(err)
	}
	defer db.Close()

	go func() {
		list, err := callTable(0, 50, "")

		json, err := json.Marshal(list)
		if err != nil {
			//return nil, err
		}
		newJson := string(json)
		connection.SetItemRedis("listUsers", newJson)
	}()

	c.JSON(200, user)
}

// Show Procura um usuario pelo id
func Show(c *gin.Context) {
	id, err := strconv.ParseInt(c.DefaultQuery("id", "1"), 10, 16)

	user := models.User{}
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

	var userMsgPack models.User

	functions.FromMSGPACK(c.Request.FormValue("code"), &userMsgPack)

	var fullUser models.User
	db := connection.CreateConnection()

	connection.ShowRow(db, table, &fullUser, "id", id)

	file, _, err := c.Request.FormFile("upload")
	path := "userfile_" + strconv.Itoa(int(id)) + ".png"
	filepath := "./tmp/" + path
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}

	if fullUser.Pathfile.Valid {
		db := connection.CreateConnection()

		err = db.Get(&userMsgPack, "UPDATE users SET username = ($2) pathfile = ($3)  WHERE id = ($1) RETURNING *", id, fullUser.Username, path)

		if err != nil {
			c.JSON(400, err)
			return
		}

		c.JSON(200, userMsgPack)
	} else {
		err = db.Get(&userMsgPack, "UPDATE users SET username = ($2) WHERE id = ($1) RETURNING *", id, userMsgPack.Username)

		if err != nil {
			c.JSON(400, err)
			return
		}
		c.JSON(200, userMsgPack)

	}

	go func() {
		list, err := callTable(0, 50, "")

		json, err := json.Marshal(list)
		if err != nil {
			//return nil, err
		}
		newJson := string(json)
		connection.SetItemRedis("listUsers", newJson)
	}()
	defer db.Close()
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
