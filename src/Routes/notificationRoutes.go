package routes

import (
	notificationController "docker.go/src/Controllers/Notification"
	middleware "docker.go/src/Middleware"

	"github.com/gin-gonic/gin"
)

// UsersRoutes rotas dos usuarios
func NotificationRoutes(route *gin.Engine) {
	auth := route.Group("notification")
	{
		auth.Use(middleware.Auth())
		auth.POST("store", notificationController.Store)
		auth.GET("index", notificationController.Index)
		auth.GET("show", notificationController.Show)
		auth.PUT("update", notificationController.Update)
		auth.DELETE("delete", notificationController.Delete)
		//auth.POST()
	}
}
