package services

import (
	"context"
	"framework-gin/dao"
	"framework-gin/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gerr"
	"sync"
)

// LogicService service layer interface
type LogicService interface {
	Health(ctx context.Context, info *pb.ReqHealth) (out *pb.RspHealth, code int)
}

var logicServiceIns *logicServiceImp
var logicServiceInitOnce sync.Once

func NewLogicServiceInstance() LogicService {

	logicServiceInitOnce.Do(func() {
		logicServiceIns = &logicServiceImp{
			dao: dao.Instance(),
		}
	})

	return logicServiceIns
}

type logicServiceImp struct {
	dao dao.Dao
}

func (g logicServiceImp) Health(ctx context.Context, info *pb.ReqHealth) (out *pb.RspHealth, code int) {
	return &pb.RspHealth{
		Msg: "ok",
	}, gerr.OK
}
