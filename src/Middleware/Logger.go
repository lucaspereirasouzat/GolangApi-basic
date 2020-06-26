package middleware

import (
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

			//url = string(url)

			var auth, _ = c.Get("auth")

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
