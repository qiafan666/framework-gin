package controllers

import (
	"context"
	"framework-gin/common/errs"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"framework-gin/ws/services"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/glog"
)

var logicController *LogicControllerImp

type LogicControllerImp struct {
	logicService services.LogicService
}

func NewLogicController() *LogicControllerImp {
	logicController = &LogicControllerImp{
		logicService: services.NewLogicServiceInstance(),
	}
	return logicController
}

func InitLogicController() {
	logicController = NewLogicController()

	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.Grp_Logic), uint8(pb.Cmd_Logic_Health), &pb.ReqHealth{}, &pb.RspHealth{}, Health)
}

// Health 健康检查
func Health(ctx context.Context, req proto.Message) (proto.Message, int) {

	// 将 proto.Message 转换为 *pb.ReqHealth
	v, ok := req.(*pb.ReqHealth)
	if !ok {
		return nil, errs.InvalidRequestError
	}
	glog.Slog.DebugKVs(ctx, "HealthHandler", "req", v)

	return logicController.logicService.Health(ctx, v)
}
