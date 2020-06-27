package routes

import (
	fileController "docker.go/src/Controllers/Images"
	middleware "docker.go/src/Middleware"

	"github.com/gin-gonic/gin"
)

// FileRoutes rotas dos arquivos
func FileRoutes(route *gin.Engine) {
	auth := route.Group("file")
	{
		//	auth.POST("store", fileController.Store)
		//auth.Use(middleware.Auth())
		//auth.Use(middleware.Auth([]string{"adm"}))
		//auth.GET("index", fileController.Index)

		auth.Use(middleware.Auth([]string{"adm", "user"}))
		auth.GET("show", fileController.Show)
		//	auth.PUT("update", fileController.Update)
		//auth.POST()
	}
}
