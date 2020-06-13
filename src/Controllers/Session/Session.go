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

// type User struct {
// 	ID       uint64 `json:"id"`
// 	Username string `json:"username"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

//A sample use
// var user = user.User{
// 	ID:       1,
// 	Username: "username",
// 	Password: "password",
// }

// type Error struct {
// 	field   string
// 	message string
// }

// type ListErrors struct {
// 	errors []Error
// }

// func generateError(fields []string) []string {
// 	var listErrors ListErrors
// 	for _, field := range fields {
// 		errorer := Error{field, "erro"}
// 		listErrors.errors = append(listErrors.errors, errorer)
// 	}
// 	return listErrors.errors
// }

// func validate(email string, password string) ListErrors {
// 	var user user.User
// 	db := connection.CreateConnection()

// 	err := db.Get(&user, "SELECT * FROM users WHERE email=$1 AND password=$2", email, password)
// 	fmt.Printf("%#v\n", user)

// 	if err != nil {
// 		fmt.Println("error")
// 		fmt.Println(err)
// 		return ListErrors{}
// 	}

// 	var arrayError []string
// 	arrayError = append(arrayError, "lucas")

// 	return generateError(arrayError)

// }

// Session Faz login do usuario
func Session(c *gin.Context) {
	var user Login
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user.Password, _ = functions.GeneratePassword(user.Password)

	var Fulluser userModels.User
	db := connection.CreateConnection()
	err := db.Get(&Fulluser, "SELECT * FROM users WHERE email=($1) ", user.Email)
	db.Close()

	var resultPassword = []byte(Fulluser.Password)

	equal := functions.ComparePasswords(user.Password, resultPassword)

	if err != nil || equal {
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
	fmt.Println(Fulluser)

	tokenString, err := functions.GenerateToken(Fulluser)

	if err != nil {
		c.JSON(400, tokenString)
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

// // Create the Signin handler
// func Signin(w http.ResponseWriter, r *http.Request) {
// 	var creds Credentials
// 	// Get the JSON body and decode into credentials
// 	err := json.NewDecoder(r.Body).Decode(&creds)
// 	if err != nil {
// 		// If the structure of the body is wrong, return an HTTP error
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

// 	// Get the expected password from our in memory map
// 	expectedPassword, ok := users[creds.Username]

// 	// If a password exists for the given user
// 	// AND, if it is the same as the password we received, the we can move ahead
// 	// if NOT, then we return an "Unauthorized" status
// 	if !ok || expectedPassword != creds.Password {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	// Declare the expiration time of the token
// 	// here, we have kept it as 5 minutes
// 	expirationTime := time.Now().Add(5 * time.Minute)
// 	// Create the JWT claims, which includes the username and expiry time
// 	claims := &Claims{
// 		Username: creds.Username,
// 		StandardClaims: jwt.StandardClaims{
// 			// In JWT, the expiry time is expressed as unix milliseconds
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	// Declare the token with the algorithm used for signing, and the claims
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	// Create the JWT string
// 	tokenString, err := token.SignedString(jwtKey)
// 	if err != nil {
// 		// If there is an error in creating the JWT return an internal server error
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	// Finally, we set the client cookie for "token" as the JWT we just generated
// 	// we also set an expiry time which is the same as the token itself
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Value:   tokenString,
// 		Expires: expirationTime,
// 	})
// }

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
