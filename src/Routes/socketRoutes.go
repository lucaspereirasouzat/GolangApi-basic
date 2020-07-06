package routes

import (
	chat "docker.go/src/Controllers/WS"
	"github.com/gin-gonic/gin"
)

// AuthRoutes rotas de autenticação
func WebsocketRoutes(route *gin.Engine) {
	hub := chat.NewHub()

	go hub.Run()
	auth := route.Group("auth")
	{
		// auth.Use(middleware.Auth([]string{"user", "adm"}))
		auth.GET("ws", func(c *gin.Context) {
			chat.Wshandler(c, hub)
		})
		// auth.POST("logout", session.Logout)
		// auth.POST("newPassword", session.RequestNewPassword)
		// auth.POST("changePassword", session.ChangePassword)
		// auth.PUT("update", session.UpdateMyUser)
		// auth.GET("myUser", session.ShowMyUser)
		//auth.POST()
	}
}
