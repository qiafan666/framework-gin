package services

import (
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gerr"
)

type PublicService struct {
}

func NewPublicService() PublicService {
	return PublicService{}
}

func (l PublicService) Health(client *internal.Client, info *pb.ReqHealth) (out *pb.RspHealth, code int) {
	return &pb.RspHealth{
		Msg: "ok",
	}, gerr.OK
}
