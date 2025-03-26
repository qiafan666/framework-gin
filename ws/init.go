package ws

import (
	"framework-gin/ws/controllers"
	"framework-gin/ws/internal"
	"framework-gin/ws/routes"
	"framework-gin/ws/services"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(r *gin.Engine) {
	// 注册路由
	routes.NewPrivateRoute(controllers.NewPrivateController(services.NewPrivateService()))
	routes.NewPublicRoute(controllers.NewPublicController(services.NewPublicService()))

	longServer := internal.NewWsServer()
	longServer.Run(r)
}
