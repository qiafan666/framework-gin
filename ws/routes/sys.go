package routes

import (
	"framework-gin/ws/errs"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
)

func RegisterSysRoutes() {
	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.GRP_SYS), uint8(pb.SYS_CMD_HEALTH), &pb.ReqHealth{}, &pb.RspHealth{}, HealthHandler)
}

func HealthHandler(req proto.Message) (proto.Message, error) {
	// 将 proto.Message 转换为 *pb.ReqHealth
	_, ok := req.(*pb.ReqHealth)
	if !ok {
		return nil, errs.ErrInvalidRequest.WrapMsg("invalid request type")
	}

	return &pb.RspHealth{
		Msg: "ok",
	}, nil
}
