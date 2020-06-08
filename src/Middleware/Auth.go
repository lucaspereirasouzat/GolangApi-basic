package middleware

import (
	"fmt"
	"log"

	"docker.go/src/functions"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	models "docker.go/src/Models/User"
)

// Auth Faz o Log do sistema
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, err := functions.VerifyToken(c.Request)
		//	fmt.Println(token)
		if err != nil {
			fmt.Println("error", err)
			panic(err)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		//fmt.Println("claims and ok", claims, ok)

		// type User struct {
		// 	ID              int64
		// 	Username, Email string
		// }
		var user models.User
		if ok && token.Valid {
			//fmt.Println(claims["User"])
			// value, err := strconv.ParseInt(claims["User"]["ID"], 0, 8)
			// if err != nil {
			// 	panic(err)
			// }
			// user := User{
			// 	value,
			// }

			data := claims["User"] //make(, 1000)

			// var m map[string]string
			// var ss []string
			// ss = strings.Split(data, "&")
			// m = make(map[string]string)
			// for _, pair := range ss {
			// 	z := strings.Split(pair, "=")
			// 	m[z[0]] = z[1]
			// }
			//delete(data, "Password")
			//v,ok := in.(data models.User)
			// for key, value := range data {
			// 	fmt.Println("Key:", key, "Value:", value)
			// }
			//like := make(data *models.User)
			//	data2 :=
			//user = make(data[string], interface{})
			// for key, value := range data.MapKeys() {
			// 	fmt.Println("Key:", key, "Value:", value)
			// }
			fmt.Println(data.(map[string]interface{})["Username"])

			// user = models.User{
			// 	Username: data.(map[string]interface{"Username"}),
			// }
			// e := data.(map[string]interface{})["ID"]
			// // id := strconv.ParseInt(, 0, 64)
			// user = models.User{
			// 	ID:       id,
			// 	Username: data.(map[string]interface{})["Username"],
			// 	data.(map[string]interface{})["Password"],
			// }

			// if !ok {
			// 	panic(err)
			// 	return
			// }
			// id, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["id"]), 10, 64)
			// if err != nil {
			// 	panic(err)
			// 	return
			// }
			// return &AccessDetails{
			// 	AccessUuid: accessUuid,
			// 	UserId:     userId,
			// }, nil
		}
		//return nil, err

		// if err := c.ShouldBindHeader(&h); err != nil {
		// 	c.JSON(200, err)
		// }
		// fmt.Println(h)

		// Save in database the token
		//user = db.QueryRow("SELECT * FROM users INNER JOIN token ON token.user_id = users.id WHERE token.token=($1); ", h.Bearer)

		// if err != nil {
		// 	fmt.Println("err", err)
		// 	//	panic(err)
		// }
		// tx.Commit()

		// db.Close()

		//c.JSON(200, "Concluido")

		c.Set("auth", user)
		// before request

		c.Next()

		// after request
		//	latency := time.Since(t)
		//	log.Print(latency)

		// access the status we are sending
		status := c.Writer.Status()
		log.Println(status)
	}
}
