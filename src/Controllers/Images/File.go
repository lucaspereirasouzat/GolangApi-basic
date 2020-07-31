package notification

import (
	"github.com/gin-gonic/gin"
)

const table = "file"

/*
 Procura uma imagem pelo id
*/
func Show(c *gin.Context) {
	path := c.Query("path")
	c.File("./tmp/" + path)
}
