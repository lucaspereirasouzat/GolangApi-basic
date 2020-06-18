package session

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"

	connection "docker.go/src/Connections"
	userModels "docker.go/src/Models/User"
	validators "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
)

// Session Faz login do usuario
func Session(c *gin.Context) {
	var user validators.Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	user.Password = functions.GenerateMD5(user.Password)

	db := connection.CreateConnection()

	var Fulluser userModels.User
	err := db.Get(&Fulluser, "SELECT * FROM users WHERE email=($1) AND password=($2)", user.Email, user.Password)

	defer db.Close()

	if err != nil {
		var list [1]validators.Error

		list[0] = validators.Error{
			Field:   "email",
			Message: "E-mail ou Senha inválidados",
		}

		type Errors struct {
			Errors [1]validators.Error
		}

		listErrors := Errors{
			Errors: list,
		}
		// erro de não encontrado
		c.JSON(400, listErrors)

		return
	}
	tokenString, err := functions.GenerateToken(Fulluser)

	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, tokenString)
}

// Logout Faz logut do user do database
func Logout(c *gin.Context) {
	type Header struct {
		Bearer string `header:"Bearer"`
	}

	h := Header{}

	if err := c.ShouldBindHeader(&h); err != nil {
		c.JSON(200, err)
	}
	db := connection.CreateConnection()
	tx := db.MustBegin()

	// Save in database the token
	tx.MustExec("DELETE FROM token WHERE token=($1); ", h.Bearer)
	tx.Commit()
	db.Close()

	c.JSON(200, "Concluido")
}

// ShowMyUser Mostra os dados do proprio usuario
func ShowMyUser(c *gin.Context) {
	var users, _ = c.Get("auth")
	c.JSON(200, users)
}

// UpdateMyUser atualiza os dados do proprio usuario
func UpdateMyUser(c *gin.Context) {
	// Dados do prorio usuario do auth
	UserGet, _ := c.Get("auth")
	myUser := UserGet.(userModels.User)

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		c.JSON(400, err)
		return
	}

	// Cria o arquivo
	file, _, err := c.Request.FormFile("upload")

	// Converte o id de int para string
	userid := strconv.Itoa(int(myUser.ID))
	filepath := "./tmp/userfile_" + userid + ".png"
	out, err := os.Create(filepath)

	if err != nil {
		c.JSON(400, err)
		return
	}

	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(400, err)
		return
	}

	if !myUser.FileId.Valid {
		db := connection.CreateConnection()

		var fileID int
		err := db.QueryRow("INSERT INTO file (path, user_id) VALUES ($1, $2) RETURNING id", filepath, userid).Scan(&fileID)

		if err != nil {
			c.JSON(400, err)
			return
		}

		var user userModels.User

		err = msgpack.Unmarshal(data, &user)

		if err != nil {
			c.JSON(400, err)
			panic(err)
		}

		err = db.Get(&user, "UPDATE users SET username = ($2) file_id=($3)  WHERE id = ($1) RETURNING *", myUser.ID, user.Username, fileID)

		if err != nil {
			c.JSON(400, err)
			return
		}

		c.JSON(200, user)

		defer db.Close()
	} else {
		var user userModels.User

		err = msgpack.Unmarshal(data, &user)
		if err != nil {
			c.JSON(400, err)
			panic(err)
		}

		db := connection.CreateConnection()
		err = db.Get(&user, "UPDATE users SET username = ($2) WHERE id = ($1) RETURNING *", myUser.ID, user.Username)

		if err != nil {
			c.JSON(400, err)
			return
		}

		c.JSON(200, user)
		defer db.Close()
	}
}
