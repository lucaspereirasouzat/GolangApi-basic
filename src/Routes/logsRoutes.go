package routes

import (
	logsController "docker.go/src/Controllers/Logs"
	middleware "docker.go/src/Middleware"

	"github.com/gin-gonic/gin"
)

// FileRoutes rotas dos arquivos
func LogsRoutes(route *gin.Engine) {
	auth := route.Group("logs")
	{
		//	auth.POST("store", fileController.Store)
		//auth.Use(middleware.Auth())
		//auth.Use(middleware.Auth([]string{"adm"}))
		//auth.GET("index", fileController.Index)

		auth.Use(middleware.Auth([]string{"adm"}))
		auth.GET("index", logsController.Index)
		//	auth.PUT("update", fileController.Update)
		//auth.POST()
	}
}
