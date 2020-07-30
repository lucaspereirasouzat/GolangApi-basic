package session

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	connection "docker.go/src/Connections"
	models "docker.go/src/Models"
	validators "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
)

// Session Faz login do usuario
func Session(c *gin.Context) {
	var user validators.Login
	user.Email = c.Request.FormValue("email")
	user.Password = c.Request.FormValue("password")
	// if err := c.ShouldBindQuery(&user); err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	// 	return
	// }
	fmt.Println("user", user)
	user.Password = functions.GenerateMD5(user.Password)
	db := connection.CreateConnection()
	var Fulluser models.User
	err := db.Get(&Fulluser, "SELECT * FROM users WHERE email=($1) AND password=($2)", user.Email, user.Password)

	defer db.Close()

	if err != nil {
		fmt.Println(err)
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
	defer db.Close()

	c.JSON(200, "Concluido")
}

// ShowMyUser Mostra os dados do proprio usuario
func ShowMyUser(c *gin.Context) {
	var user, _ = c.Get("auth")
	c.JSON(200, user)
}

// UpdateMyUser atualiza os dados do proprio usuario
func UpdateMyUser(c *gin.Context) {
	// Dados do prorio usuario do auth
	UserGet, _ := c.Get("auth")
	myUser := UserGet.(models.User)

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		c.JSON(400, err)
		return
	}

	file, _, err := c.Request.FormFile("upload")
	path := "userfile_" + strconv.Itoa(int(myUser.ID)) + ".png"
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

	if len(myUser.Pathfile.String) == 0 {
		var user models.User

		err = msgpack.Unmarshal(data, &user)

		if err != nil {
			c.JSON(400, err)
			panic(err)
		}

		db := connection.CreateConnection()
		err = db.Get(&user, "UPDATE users SET username = ($2), pathfile = ($3)  WHERE id = ($1) RETURNING *", myUser.ID, user.Username, path)
		defer db.Close()

		if err != nil {
			c.JSON(400, err)
			return
		}

		c.JSON(200, user)

	} else {
		var user models.User

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

// RequestNewPassword envia o email
func RequestNewPassword(c *gin.Context) {
	email := c.Request.FormValue("email")
	user := models.User{}

	db := connection.CreateConnection()
	err := connection.ShowRow(db, "users", &user, "email", email)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, err)
		return
	}

	item := struct {
		token  string
		userID uint64
	}{functions.RandStringBytesRmndr(30), user.ID}

	_, err = db.Exec("INSERT INTO token (token,user_id) values ($1,$2)", item.token, item.userID)

	if err != nil {
		fmt.Println(err)
		c.JSON(400, err)
		return
	}

	defer db.Close()

	// send email
	c.JSON(200, item.token)
}

// ChangePassword Faz a troca de senha com o token
func ChangePassword(c *gin.Context) {
	token := c.Request.FormValue("token")
	db := connection.CreateConnection()

	var userID int
	err := db.Get(&userID, "UPDATE token SET is_revoked = $2 WHERE token = $1 RETURNING user_id", token, true)

	if err != nil {
		c.JSON(400, err)
		return
	}

	password := functions.GenerateMD5(c.Request.FormValue("password"))

	user := models.User{}

	err = db.Get(&user, "UPDATE users SET password = ($2) WHERE id = ($1) RETURNING *", userID, password)
	if err != nil {
		c.JSON(400, err)
		return
	}

	defer db.Close()
	c.JSON(200, user)
}
