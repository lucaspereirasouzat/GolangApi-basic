package routes

import (
	fileController "docker.go/src/Controllers/Images"

	"github.com/gin-gonic/gin"
)

// UsersRoutes rotas dos usuarios
func FileRoutes(route *gin.Engine) {
	auth := route.Group("file")
	{
		//	auth.POST("store", fileController.Store)
		//auth.Use(middleware.Auth())
		auth.GET("index", fileController.Index)
		auth.GET("show", fileController.Show)
		//	auth.PUT("update", fileController.Update)
		//auth.POST()
	}
}
