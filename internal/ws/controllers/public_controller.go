package controllers

import (
	"framework-gin/internal/ws/internal"
	"framework-gin/internal/ws/proto/pb"
	"framework-gin/internal/ws/services"
	"github.com/qiafan666/gotato/commons/gerr"
	"google.golang.org/protobuf/proto"
)

type PublicController struct {
	publicService services.PublicService
}

func NewPublicController() PublicController {
	return PublicController{
		publicService: services.NewPublicService(),
	}
}

// Health 健康检查
func (p PublicController) Health(client *internal.Client, req proto.Message) (proto.Message, int) {
	// 将 proto.Message 转换为 *pb.ReqHealth
	pbReq, ok := req.(*pb.ReqHealth)
	if !ok {
		//glog.Slog.ErrorKVs(client.UserCtx.TraceCtx, "req type error", "req", req)
		return nil, gerr.ParameterError
	}
	return p.publicService.Health(client, pbReq)
}
