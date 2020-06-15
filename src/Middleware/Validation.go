package middleware

import (
	"encoding/base64"
	"fmt"

	validatores "docker.go/src/Validators"

	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
)

// Validation valida
func Validation(dataItem interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Request.FormValue("code")
		fmt.Println(code)
		fmt.Println(dataItem.(validatores.Register).Username)

		data, err := base64.StdEncoding.DecodeString(code)
		if err != nil {

			panic(err)
		}

		//var user validatores.Register

		err = msgpack.Unmarshal(data, &dataItem)
		fmt.Println(data)
		if err != nil {
			fmt.Println("error in conversion")

			panic(err)
		}

		// // faz a verificação do token enviado no bear token do header
		// token, err := functions.VerifyToken(c.Request)

		// if err != nil {

		// 	fmt.Println("error", err)
		// 	// Envia o erro caso não consiga verificar o erro
		// 	c.JSON(404, "Token incorreto ou não enviado")
		// 	return
		// 	// Quebra a aplicação
		// 	panic(err)
		// }
		// // faz o map do Usuario e verifica se está ok
		// claims, ok := token.Claims.(jwt.MapClaims)
		// // Cria o usuario que será enviado
		// var user models.User
		// // Verifica se está valido apos fazer o map
		// if ok && token.Valid {

		// 	data := claims["User"]

		// 	var usermaped = data.(map[string]interface{})

		// 	//	securelevel := usermaped["Securelevel"].(string)

		// 	user.ID = uint64(usermaped["ID"].(float64))
		// 	user.Username = usermaped["Username"].(string)
		// 	user.Email = usermaped["Email"].(string)
		// 	user.Securelevel = usermaped["Securelevel"].(string)
		// 	fileResult := usermaped["FileId"].(map[string]interface{})
		// 	user.FileId.Valid = fileResult["Valid"].(bool)
		// 	user.FileId.Int64 = int64(fileResult["Int64"].(float64))
		// 	//fmt.Println(usermaped["CreatedAt"].(string))
		// 	//user.CreatedAt, err = time.Parse("0001-01-01T00:00:00Z", usermaped["CreatedAt"].(string))
		// }
		// fmt.Println(user)
		// c.Set("auth", user)
		// // before request

		c.Next()

		// // access the status we are sending
		// status := c.Writer.Status()
		// log.Println(status)
	}
}
