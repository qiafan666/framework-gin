package ws

import (
	"framework-gin/ws/controllers"
	"framework-gin/ws/internal"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(r *gin.Engine) {
	// 注册路由
	controllers.InitSysController()
	controllers.InitLogicController()

	longServer := internal.NewWsServer()
	go longServer.ChangeOnlineStatus(4)
	longServer.Run(r)
}
