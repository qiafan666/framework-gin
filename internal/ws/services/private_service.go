package services

import (
	"framework-gin/internal/ws/internal"
	"framework-gin/internal/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gerr"
)

type PrivateService struct {
}

func NewPrivateService() PrivateService {
	return PrivateService{}
}

func (l PrivateService) Health(client *internal.Client, info *pb.ReqHealth) (out *pb.RspHealth, code int) {
	return &pb.RspHealth{
		Msg: "ok",
	}, gerr.OK
}
