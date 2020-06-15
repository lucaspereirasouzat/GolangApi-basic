package routes

import (
	userController "docker.go/src/Controllers/User"
	middleware "docker.go/src/Middleware"

	"github.com/gin-gonic/gin"
)

// UsersRoutes rotas dos usuarios
func UsersRoutes(route *gin.Engine) {
	//var register validatores.Register
	auth := route.Group("user")
	{
		auth.POST("store", userController.Store)
		auth.Use(middleware.Auth([]string{"user", "ADM"}))
		auth.GET("index", userController.Index)
		auth.GET("show", userController.Show)
		auth.PUT("update", userController.Update)
		//auth.POST()
	}
}
