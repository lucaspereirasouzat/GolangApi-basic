package session

import (
	// "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

//userModel "docker.go/src/Models/User"
// func Session(c *gin.Context) {

// 	c.JSON(200, gin.H{
// 		"username": "lucas",
// 		"password": 1234,
// 		"email":    "lucas@teste.com",
// 	})
// }

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//A sample use
var user = User{
	ID:       1,
	Username: "username",
	Password: "password",
}

// Faz login do usuario
func Session(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	// token, err := CreateToken(user.ID)
	// if err != nil {
	// 	c.JSON(http.StatusUnprocessableEntity, err.Error())
	// 	return
	// }
	// c.JSON(http.StatusOK, token)
}
