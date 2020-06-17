package session

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	connection "docker.go/src/Connections"
	userModels "docker.go/src/Models/User"
	validators "docker.go/src/Validators"
	"docker.go/src/functions"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
)

type Login struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type Token struct {
	User userModels.User
	jwt.StandardClaims
}

// Session Faz login do usuario
func Session(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user.Password = functions.GenerateMD5(user.Password)
	var Fulluser userModels.User

	db := connection.CreateConnection()
	//err := connection.ShowRow(db, "users", &Fulluser, "email", user.Email)
	err := db.Get(&Fulluser, "SELECT * FROM users WHERE email=($1) AND password=($2)", user.Email, user.Password)

	defer db.Close()

	if err != nil {
		fmt.Println(err)
		myerror := validators.Error{
			Field:   "email",
			Message: "E-mail ou Senha inválidados",
		}

		var list [1]validators.Error
		list[0] = myerror

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
	//fmt.Println(Fulluser)

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

/*
 Mostra os dados do proprio usuario do auth
*/
func ShowMyUser(c *gin.Context) {
	var users, _ = c.Get("auth")
	c.JSON(200, users)
}

/*
 Atualiza um novo usuario pelo id
*/
func UpdateMyUser(c *gin.Context) {

	UserGet, _ := c.Get("auth")
	us := UserGet.(userModels.User)

	data, err := base64.StdEncoding.DecodeString(c.Request.FormValue("code"))
	if err != nil {
		panic(err)
	}

	file, _, err := c.Request.FormFile("upload")
	userid := strconv.Itoa(int(us.ID))
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

	if !us.FileId.Valid {
		db := connection.CreateConnection()
		tx := db.MustBegin()

		tx.MustExec("INSERT INTO file (path, user_id) VALUES ($1, $2)", filepath, userid)
		tx.Commit()

		if err != nil {
			fmt.Println(err)
			return
		}

		type FileResult struct {
			ID uint64
		}

		var fileres FileResult
		err = db.Get(&fileres, "SELECT id FROM file WHERE path = ($1)", filepath)
		if err != nil {
			fmt.Println("error on select", err)
			return
		}
		fmt.Println(fileres)
		var user userModels.User

		err = msgpack.Unmarshal(data, &user)
		if err != nil {
			fmt.Println("error in conversion")
			panic(err)
		}

		tx = db.MustBegin()
		result := tx.MustExec("UPDATE users SET username = ($2), file_id = ($3) WHERE id = ($1)", us.ID, user.Username, fileres.ID)
		fmt.Println(result)
		if err != nil {
			fmt.Println(err)
			return
		}

		tx.Commit()

		db.Close()

		fmt.Printf("%#v\n", user)

		c.JSON(200, user)
	} else {
		var user userModels.User

		err = msgpack.Unmarshal(data, &user)
		if err != nil {
			fmt.Println("error in conversion")
			panic(err)
		}
		db := connection.CreateConnection()

		tx := db.MustBegin()
		result := tx.MustExec("UPDATE users SET username = ($2) WHERE id = ($1)", us.ID, user.Username)
		fmt.Println(result)
		if err != nil {
			fmt.Println(err)
			return
		}

		tx.Commit()

		db.Close()

		fmt.Printf("%#v\n", user)

		c.JSON(200, user)
	}

}
