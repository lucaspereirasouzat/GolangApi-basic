package routes

import (
	session "docker.go/src/Controllers/Session"
	middleware "docker.go/src/Middleware"
	"github.com/gin-gonic/gin"
)

// AuthRoutes rotas de autenticação
func AuthRoutes(route *gin.Engine) {
	auth := route.Group("auth")
	{
		auth.POST("session", session.Session)
		auth.POST("logout", session.Logout)
		auth.Use(middleware.Auth())
		auth.GET("myUser", session.ShowMyUser)
		//auth.POST()
	}
}
