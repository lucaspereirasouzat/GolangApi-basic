package middleware

import (
	"fmt"
	"log"
	"time"

	"docker.go/src/functions"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	models "docker.go/src/Models/User"
)

// Auth Faz o Log do sistema
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// faz a verificação do token enviado no bear token do header
		token, err := functions.VerifyToken(c.Request)

		if err != nil {
			fmt.Println("error", err)
			// Envia o erro caso não consiga verificar o erro
			c.JSON(404, "Token incorreto ou não enviado")
			// Quebra a aplicação
			panic(err)
		}
		// faz o map do Usuario e verifica se está ok
		claims, ok := token.Claims.(jwt.MapClaims)
		// Cria o usuario que será enviado
		var user models.User
		// Verifica se está valido apos fazer o map
		if ok && token.Valid {

			data := claims["User"]

			var usermaped = data.(map[string]interface{})

			user.ID = uint64(usermaped["ID"].(float64))
			user.Username = usermaped["Username"].(string)
			user.Email = usermaped["Email"].(string)
			user.CreatedAt, err = time.Parse("2006/01/02 03:04:05", usermaped["CreatedAt"].(string))
		}
		c.Set("auth", user)
		// before request

		c.Next()

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
