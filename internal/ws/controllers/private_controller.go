package controllers

import (
	"framework-gin/internal/ws/internal"
	"framework-gin/internal/ws/proto/pb"
	"framework-gin/internal/ws/services"
	"github.com/qiafan666/gotato/commons/gerr"
	"google.golang.org/protobuf/proto"
)

type PrivateController struct {
	privateService services.PrivateService
}

func NewPrivateController() PrivateController {
	return PrivateController{
		privateService: services.NewPrivateService(),
	}
}

// Health 健康检查
func (p PrivateController) Health(client *internal.Client, req proto.Message) (proto.Message, int) {
	// 将 proto.Message 转换为 *pb.ReqHealth
	pbReq, ok := req.(*pb.ReqHealth)
	if !ok {
		//glog.Slog.ErrorKVs(client.UserCtx.TraceCtx, "req type error", "req", req)
		return nil, gerr.ParameterError
	}
	return p.privateService.Health(client, pbReq)
}
