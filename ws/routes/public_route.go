package routes

import (
	"framework-gin/ws/controllers"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
)

type PublicRoute struct {
	publicController controllers.PublicController
}

func NewPublicRoute(
	publicController controllers.PublicController,
) PublicRoute {
	return PublicRoute{
		publicController: publicController,
	}
}

func (p PublicRoute) PublicSetup() {
	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.Grp_Public), uint8(pb.CmdPublic_HealthCheck), &pb.ReqHealth{}, &pb.RspHealth{}, p.publicController.Health)
}
