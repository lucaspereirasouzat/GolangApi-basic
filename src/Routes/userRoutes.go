package routes

import (
	userController "docker.go/src/Controllers/User"
	"github.com/gin-gonic/gin"
)

func UsersRoutes(route *gin.Engine) {
	auth := route.Group("user")
	{

		auth.GET("index", userController.Index)
		auth.POST("store", userController.Store)
		auth.GET("show", userController.Show)
		auth.PUT("update", userController.Update)
		//auth.POST()
	}
}
