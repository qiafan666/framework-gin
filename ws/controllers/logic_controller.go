package controllers

import (
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"framework-gin/ws/services"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gerr"
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
func Health(client *internal.Client, req proto.Message) (proto.Message, int) {

	// 将 proto.Message 转换为 *pb.ReqHealth
	pbReq, ok := req.(*pb.ReqHealth)
	if !ok {
		glog.Slog.ErrorKVs(client.UserCtx.Ctx, "req type error", "req", req)
		return nil, gerr.ParameterError
	}
	return logicController.logicService.Health(client, pbReq)
}
