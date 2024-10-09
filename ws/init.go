package ws

import (
	"framework-gin/ws/internal"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(r *gin.Engine) {
	longServer := internal.NewWsServer()
	longServer.Run(r)
}
