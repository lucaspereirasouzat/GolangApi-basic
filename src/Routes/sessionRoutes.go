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
		auth.POST("newPassword", session.RequestNewPassword)
		auth.POST("changePassword", session.ChangePassword)
		auth.Use(middleware.Auth([]string{"user", "adm"}))
		auth.PUT("update", session.UpdateMyUser)
		auth.GET("myUser", session.ShowMyUser)
		//auth.POST()
	}
}
