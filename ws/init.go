package ws

import (
	"framework-gin/ws/internal"
	"framework-gin/ws/routes"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(r *gin.Engine) {

	// 注册路由
	routes.RegisterSysRoutes()

	longServer := internal.NewWsServer()
	longServer.Run(r)
}
