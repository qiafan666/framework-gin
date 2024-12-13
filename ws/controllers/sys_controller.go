package controllers

import (
	"framework-gin/ws/internal"
	"framework-gin/ws/proto/pb"
	"framework-gin/ws/services"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
)

var sysController *SysControllerImp

type SysControllerImp struct {
	sysService services.SysService
}

func NewSysController() *SysControllerImp {
	sysController = &SysControllerImp{
		sysService: services.NewSysServiceInstance(),
	}
	return sysController
}

func InitSysController() {
	sysController = NewSysController()

	handler := internal.GetMsgHandler()
	handler.AddHandler(uint8(pb.Grp_Sys), uint8(pb.CmdSys_SubscribeOnlineUser), &pb.ReqSubUserOnlineStatus{}, &pb.RspSubUserOnlineStatus{}, SubUserOnlineStatus)
	handler.AddHandler(uint8(pb.Grp_Sys), uint8(pb.CmdSys_PushMessage), &pb.ReqPushMsgToOther{}, &pb.RspPushMsgToOther{}, PushMsgToOther)
}

// SubUserOnlineStatus 订阅在线用户状态
func SubUserOnlineStatus(client *internal.Client, req proto.Message) (proto.Message, int) {
	pbReq, ok := req.(*pb.ReqSubUserOnlineStatus)
	if !ok {
		glog.Slog.ErrorKVs(client.UserCtx.TraceCtx, "req type error", "req", req)
		return nil, gerr.ParameterError
	}

	return sysController.sysService.SubUserOnlineStatus(client, pbReq)
}

// PushMsgToOther 推送消息到其他人
func PushMsgToOther(client *internal.Client, req proto.Message) (proto.Message, int) {
	pbReq, ok := req.(*pb.ReqPushMsgToOther)
	if !ok {
		glog.Slog.ErrorKVs(client.UserCtx.TraceCtx, "req type error", "req", req)
		return nil, gerr.ParameterError
	}

	return sysController.sysService.PushMsgToOther(client, pbReq)
}
