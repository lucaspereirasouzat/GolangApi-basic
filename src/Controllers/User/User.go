package user

import (
	"encoding/json"
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
)

type Lista struct {
	Page        uint64
	RowsPerPage uint64
	Total       uint64
	Table       []models.User
}

// Chamada para o Postgres dos Usuarios
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

	if err != nil {
		return Lista{}, err
	}

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
	page, err := strconv.ParseUint(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseUint(c.DefaultQuery("RowsPerPage", "50"), 10, 10)
	search := c.DefaultQuery("search", "")

	var list Lista

	if page == 0 && rowsPerPage == 50 && search == "" {
		// Pega os dados do redis e faz a validação
		result, err := connection.GetItemRedis("listUsers")
		// Se ouver erro no redis ele fara a chamada normalmente para o Postgres
		if err != nil {
			list, err = callTable(page, rowsPerPage, search)

			if err != nil {
				c.JSON(400, err)
				return
			}

			go func() {
				json, err := json.Marshal(list)
				if err != nil {
					//c.JSON(400, err)
					return
				}
				newJson := string(json)
				connection.SetItemRedis("listUsers", newJson)
			}()

		} else {
			err := json.Unmarshal([]byte(result), &list)
			if err != nil {
				c.JSON(400, err)
				return
			}
		}
	} else {
		list, err = callTable(page, rowsPerPage, search)

		if err != nil {
			c.JSON(400, err)
			return
		}
	}

	// responde para o front com o json
	c.JSON(200, list)
}

// Store Cadastra um novo usuario no sistema
func Store(c *gin.Context) {
	var user validatores.Register
	err := functions.FromMSGPACK(c.Request.FormValue("code"), &user)

	hasError, listError := validatores.Validate(user)

	if hasError {
		c.JSON(400, listError)
		return
	}

	// Gera o md5 da senha
	user.Password = functions.GenerateMD5(user.Password)

	// cria conexão com banco de dados
	db := connection.CreateConnection()
	err = db.Get(&user, "INSERT INTO "+table+" (username,email,password) VALUES ($1,$2,$3)  RETURNING username,email,password", user.Username, user.Email, user.Password)

	if err != nil {
		switch err.Error() {
		case `pq: duplicate key value violates unique constraint "users_email_key"`:
			var list [1]validatores.Error

			list[0] = validatores.Error{
				Field:   "email",
				Message: "E-mail duplicado",
			}

			type Errors struct {
				Errors [1]validatores.Error
			}

			listErrors := Errors{
				Errors: list,
			}
			// erro de não encontrado
			c.JSON(400, listErrors)
			break
		default:
			c.String(400, "%s", err)
			break
		}
		return
		//panic(err)
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
		c.JSON(400, "Não foram encotrados campos com este id")
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
	defer db.Close()

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, user)
}
