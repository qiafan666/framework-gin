package routes

import (
	"framework-gin/ws/controllers"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
)

type PrivateRoute struct {
	privateController controllers.PrivateController
}

func NewPrivateRoute(
	privateController controllers.PrivateController,
) PrivateRoute {
	return PrivateRoute{
		privateController: privateController,
	}
}

func (p PrivateRoute) PrivateSetup() {
	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.Grp_Private), uint8(pb.CmdPrivate_Health), &pb.ReqHealth{}, &pb.RspHealth{}, p.privateController.Health)
}
