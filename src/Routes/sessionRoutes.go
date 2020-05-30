package routes

import (
	session "docker.go/src/Controllers/Session"
	"github.com/gin-gonic/gin"
)

// rotas de autenticação
func AuthRoutes(route *gin.Engine) {
	auth := route.Group("auth")
	{
		auth.POST("session", session.Session)
		// auth.GET("index/uers", userController.Index)
		//auth.POST()
	}
}