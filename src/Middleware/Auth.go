package middleware

import (
	"docker.go/src/functions"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	models "docker.go/src/Models/User"
)

// Auth Faz o Log do sistema
func Auth(list []string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if len(list) == 0 {
			c.Next()
			return
		}
		// faz a verificação do token enviado no bear token do header
		token, err := functions.VerifyToken(c.Request)

		if err != nil {

			if functions.Contains(list, "guest") {
				c.Next()
				return
			}

			// Envia o erro caso não consiga verificar o erro
			c.JSON(404, "Token incorreto ou não enviado")
			panic(err)
			return
		}
		// faz o map do Usuario e verifica se está ok
		claims, ok := token.Claims.(jwt.MapClaims)
		// Cria o usuario que será enviado
		var user models.User
		// Verifica se está valido apos fazer o map
		if ok && token.Valid {

			data := claims["User"]

			var usermaped = data.(map[string]interface{})

			securelevel := usermaped["Securelevel"].(string)

			if !functions.Contains(list, securelevel) {
				c.String(404, "Você não tem permissão para fazer isso")
				c.Abort()
				panic(err)
				return
			}

			user.ID = uint64(usermaped["ID"].(float64))
			user.Username = usermaped["Username"].(string)
			user.Email = usermaped["Email"].(string)
			user.Securelevel = usermaped["Securelevel"].(string)
			fileResult := usermaped["FileId"].(map[string]interface{})
			user.FileId.Valid = fileResult["Valid"].(bool)
			user.FileId.Int64 = int64(fileResult["Int64"].(float64))
			//fmt.Println(usermaped["CreatedAt"].(string))
			//user.CreatedAt, err = time.Parse("0001-01-01T00:00:00Z", usermaped["CreatedAt"].(string))
		}
		//fmt.Println(user)

		c.Set("auth", user)
		// before request

		//c.Next()

		// access the status we are sending
		// status := c.Writer.Status()
		// log.Println(status)
	}
}
