package session

import (
	"fmt"
	"net/http"

	connection "docker.go/src/Connections"
	user "docker.go/src/Models/User"
	"docker.go/src/functions"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validator "github.com/go-playground/validator/v10"
)

// 	"github.com/go-playground/validator/v10"
// )

// type Login struct {
// 	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
// 	Password string `form:"password" json:"password" xml:"password" binding:"required"`
// }

type Login struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type Token struct {
	User user.User
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
	var user user.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	// validate := validator.New()
	// validateStruct(user, validate)

	db := connection.CreateConnection()
	err := db.Get(&user, "SELECT * FROM users WHERE email=($1) AND password=($2)", user.Email, user.Password)
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		// erro de não encontrado
		c.JSON(404, err)

		return
	}

	tokenString, err := functions.GenerateToken(user)

	if err != nil {
		c.JSON(400, tokenString)
		return
	}
	// db = connection.CreateConnection()
	// tx := db.MustBegin()

	// // Save in database the token
	// tx.MustExec("INSERT INTO token (token, user_id) VALUES ($1, $2)", tokenString, user.ID)
	// tx.Commit()
	// db.Close()

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
 Procura um novo usuario pelo id
*/
func ShowMyUser(c *gin.Context) {
	var users, _ = c.Get("auth")
	c.JSON(200, users)
}

func validateStruct(user user.User, validate *validator.Validate) {

	//	validate := validator.New()

	// returns nil or ValidationErrors ( []FieldError )
	err := validate.Struct(user)
	fmt.Println(err)
	//return err
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			fmt.Println(err)
			return
		}

		for _, err := range err.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}

		// from here you can create your own error messages in whatever language you wish
		return
	}

	// save user to database
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
