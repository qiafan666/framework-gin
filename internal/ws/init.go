package ws

import (
	"context"
	"framework-gin/internal/ws/controllers"
	"framework-gin/internal/ws/internal"
	"framework-gin/internal/ws/proto/pb"
	"github.com/gin-gonic/gin"
)

// Register Start run ws server.
func Register(ctx context.Context, r *gin.Engine) {
	// 注册路由
	handler := internal.GetMsgHandler()

	privateController := controllers.NewPrivateController()
	{
		handler.AddHandler(uint8(pb.Grp_Private), uint8(pb.CmdPrivate_Health),
			&pb.ReqHealth{}, &pb.RspHealth{}, privateController.Health)

	}

	publicController := controllers.NewPublicController()
	{
		handler.AddHandler(uint8(pb.Grp_Public), uint8(pb.CmdPublic_HealthCheck),
			&pb.ReqHealth{}, &pb.RspHealth{}, publicController.Health)
	}

	server := internal.NewWsServer(ctx)
	server.Run(r)
}
