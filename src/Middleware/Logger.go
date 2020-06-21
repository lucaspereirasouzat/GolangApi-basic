package middleware

import (
	"fmt"
	"time"

	connection "docker.go/src/Connections"
	"github.com/gin-gonic/gin"
)

// type Log struct {
// 	time    string      `bson:"time"`
// 	latency string      `bson:"latency"`
// 	status  int         `bson:"status"`
// 	url     string      `bson:"url"`
// 	user    interface{} `bson:"user"`
// }

// Logger Faz o Log do sistema
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		go func() {

			t := time.Now()
			timeString := t.String()
			//after request
			//	latency := time.Since(t)
			latencyString := time.Since(t).String()

			// access the status we are sending
			status := c.Writer.Status()

			url := c.Request.URL.String()
			fmt.Println(url)
			//url = string(url)

			var auth, _ = c.Get("auth")

			fmt.Println("Method", c.Request.Method)
			fmt.Println(c.Writer.Header())
			log := connection.Log{
				Time:    timeString,
				Latency: latencyString,
				Status:  status,
				Url:     url,
				User:    auth,
				Method:  c.Request.Method,
			}

			connection.InsertMongoDB("Log", "logs", log)
		}()

	}
}
