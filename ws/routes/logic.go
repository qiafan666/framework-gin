package routes

import (
	"context"
	"framework-gin/common"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/commons/glog"
)

func RegisterLogicRoutes() {
	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.Grp_Logic), uint8(pb.Cmd_Logic_Health), &pb.ReqHealth{}, &pb.RspHealth{}, HealthHandler)
}

// HealthHandler 健康检查
func HealthHandler(ctx context.Context, req proto.Message) (proto.Message, int) {

	// 将 proto.Message 转换为 *pb.ReqHealth
	v, ok := req.(*pb.ReqHealth)
	if !ok {
		return nil, common.InvalidRequestError
	}
	glog.Slog.DebugKVs(ctx, "HealthHandler", "req", v)

	return &pb.RspHealth{Msg: "world"}, commons.OK
}
