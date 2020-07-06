package logs

import (
	"strconv"

	connection "docker.go/src/Connections"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Lista struct {
	Page        int64
	RowsPerPage int64
	Total       int64
	Table       []primitive.M
}

func CallTableLog(page int64, rowsPerPage int64, search string) (list Lista, err error) {
	// query := functions.SearchFields(search, []string{"username", "email", "secureLevel"})
	// selectFields := functions.SelectFields([]string{"id", "username", "email", "securelevel", "created_at"})
	result, total, err := connection.SelectMongoDB("Log", "logs", search, page*rowsPerPage, rowsPerPage)
	if err != nil {
		return Lista{}, err
	}
	return Lista{Page: page, RowsPerPage: rowsPerPage, Total: total, Table: result}, nil
}

/*
	Faz listagem de todos os tokens de notificação
*/
func Index(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "0"), 10, 8)
	rowsPerPage, err := strconv.ParseInt(c.DefaultQuery("rowsPerPage", "50"), 10, 10)
	search := c.DefaultQuery("search", "")

	if err != nil {
		c.JSON(400, err)
		panic(err)
	}
	// result, total, err := connection.SelectMongoDB("Log", "logs", search, page*rowsPerPage, rowsPerPage)
	// if err != nil {

	// 	c.JSON(400, err)
	// }
	var list Lista

	list, err = CallTableLog(page, rowsPerPage, search)
	// if page == 0 && rowsPerPage == 50 && search == "" {
	// 	result, err := connection.GetItemRedis("listLogs")
	// 	//fmt.Println("result", result)
	// 	if err != nil {
	// 		list, err = CallTableLog(page, rowsPerPage, search)

	// 		if err != nil {
	// 			c.JSON(400, err)
	// 			return
	// 		}

	// 		go func() {
	// 			json, err := json.Marshal(list)
	// 			if err != nil {
	// 				c.JSON(400, err)
	// 				return
	// 			}
	// 			newJson := string(json)
	// 			connection.SetItemRedis("listLogs", newJson)
	// 		}()

	// 	} else {
	// 		err := json.Unmarshal([]byte(result), &list)
	// 		if err != nil {
	// 			c.JSON(400, err)
	// 			return
	// 		}
	// 	}
	// } else {
	// 	list, err = CallTableLog(page, rowsPerPage, search)

	// 	if err != nil {
	// 		c.JSON(400, err)
	// 		return
	// 	}
	// }
	// // b, err := msgpack.Marshal(list)
	// // if err != nil {
	// // 	panic(err)
	// // }

	c.JSON(200, list)

}
