package services

import (
	"framework-gin/dao"
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gerr"
	"sync"
)

// SysService service layer interface
type SysService interface {
	SubUserOnlineStatus(client *internal.Client, info *pb.ReqSubUserOnlineStatus) (out *pb.RspSubUserOnlineStatus, code int)
	PushMsgToOther(client *internal.Client, info *pb.ReqPushMsgToOther) (out *pb.RspPushMsgToOther, code int)
}

var sysServiceIns *sysServiceImp
var sysServiceInitOnce sync.Once

func NewSysServiceInstance() SysService {

	sysServiceInitOnce.Do(func() {
		sysServiceIns = &sysServiceImp{
			dao: dao.Instance(),
		}
	})

	return sysServiceIns
}

type sysServiceImp struct {
	dao dao.Dao
}

func (s sysServiceImp) SubUserOnlineStatus(client *internal.Client, info *pb.ReqSubUserOnlineStatus) (out *pb.RspSubUserOnlineStatus, code int) {
	return client.LongConnServer.SubUserOnlineStatus(client, info)
}
func (s sysServiceImp) PushMsgToOther(client *internal.Client, info *pb.ReqPushMsgToOther) (out *pb.RspPushMsgToOther, code int) {

	err := client.PubMessage(client.UserCtx.TraceCtx, info)
	if err != nil {
		return nil, gerr.UnKnowError
	}
	return &pb.RspPushMsgToOther{}, gerr.OK
}
