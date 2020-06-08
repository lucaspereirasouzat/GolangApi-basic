package routes

import (
	session "docker.go/src/Controllers/Session"
	"github.com/gin-gonic/gin"
)

// AuthRoutes rotas de autenticação
func AuthRoutes(route *gin.Engine) {
	auth := route.Group("auth")
	{
		auth.POST("session", session.Session)
		auth.POST("logout", session.Logout)
		//auth.POST()
	}
}
