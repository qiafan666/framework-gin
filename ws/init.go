package ws

import (
	"context"
	"framework-gin/ws/controllers"
	"framework-gin/ws/internal"
	"framework-gin/ws/routes"
	"framework-gin/ws/services"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(ctx context.Context, r *gin.Engine) {
	// 注册路由
	routes.NewPrivateRoute(controllers.NewPrivateController(services.NewPrivateService())).PrivateSetup()
	routes.NewPublicRoute(controllers.NewPublicController(services.NewPublicService())).PublicSetup()

	server := internal.NewWsServer(ctx)
	server.Run(r)
}
