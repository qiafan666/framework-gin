package internal

import (
	"context"
	"framework-gin/common/function"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin/jsonutil"
	"github.com/qiafan666/gotato/commons/glog"
	"sync"
)

type Req struct {
	//ReqIdentifier int32  `json:"reqIdentifier" validate:"required"`
	//Token       string `json:"token"`
	SendID      string `json:"sendID"        validate:"required"`
	OperationID string `json:"operationID"   validate:"required"`
	//MsgIncr       string `json:"msgIncr"       validate:"required"`
	GrpID uint8  `json:"grpID" validate:"required"` // 消息组id
	CmdID uint8  `json:"cmdID" validate:"required"` // 消息的ID
	Data  []byte `json:"data"`
}

func (r *Req) String() string {
	var tReq Req
	//tReq.ReqIdentifier = r.ReqIdentifier
	//tReq.Token = r.Token
	tReq.SendID = r.SendID
	tReq.OperationID = r.OperationID
	//tReq.MsgIncr = r.MsgIncr
	tReq.GrpID = r.GrpID
	tReq.CmdID = r.CmdID
	tReq.Data = r.Data
	return jsonutil.StructToJsonString(tReq)
}

var reqPool = sync.Pool{
	New: func() any {
		return new(Req)
	},
}

func getReq() *Req {
	req := reqPool.Get().(*Req)
	req.Data = nil
	//req.MsgIncr = ""
	req.OperationID = ""
	//req.ReqIdentifier = 0
	req.SendID = ""
	//req.Token = ""
	req.GrpID = 0
	req.CmdID = 0
	return req
}

func freeReq(req *Req) {
	reqPool.Put(req)
}

type Resp struct {
	//ReqIdentifier int32  `json:"reqIdentifier"`
	//MsgIncr       string `json:"msgIncr"`
	GrpID       uint8         `json:"grp_id"`
	CmdID       uint8         `json:"cmd_id"`
	OperationID string        `json:"operation_id"`
	Code        int           `json:"code"`
	Msg         string        `json:"msg"`
	Data        proto.Message `json:"data"`
}

func (r *Resp) String() string {
	var tResp Resp
	//tResp.ReqIdentifier = r.ReqIdentifier
	//tResp.MsgIncr = r.MsgIncr
	tResp.OperationID = r.OperationID
	tResp.Code = r.Code
	tResp.Msg = r.Msg
	tResp.GrpID = r.GrpID
	tResp.CmdID = r.CmdID
	tResp.Data = r.Data
	return jsonutil.StructToJsonString(tResp)
}

var handler *MsgHandle

type IRouter interface {
	Do(req *Request)
}

// HandlerFunc 消息处理函数
type HandlerFunc func(req proto.Message) (proto.Message, error)

type Handler struct {
	f   HandlerFunc   // 业务处理函数
	req proto.Message // 请求pb
	rsp proto.Message // 返回pb
}

// MsgHandle 消息管理
type MsgHandle struct {
	Apis map[uint32]*Handler // 存放每个MsgID 所对应的处理方法的map属性
}

func init() {
	handler = NewMsgHandle()
}

// NewMsgHandle 创建MsgHandle
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]*Handler),
	}
}

func GetMsgHandler() *MsgHandle {
	return handler
}

// DoMsgHandler 处理业务
func (m *MsgHandle) DoMsgHandler(ctx context.Context, req *Req) (proto.Message, error) {
	msgID := genMsgID(req.GrpID, req.CmdID)
	h, ok := m.Apis[msgID]
	if !ok {
		glog.Slog.ErrorF(ctx, "DoMsgHandler msgID: %v not found", msgID)
		return nil, gerr.New("msgID not found")
	}

	// 解析pb消息
	var dataReq proto.Message
	if h.req != nil {
		dataReq = proto.Clone(h.req)
		if err := proto.Unmarshal(req.Data, dataReq); err != nil {
			return nil, gerr.New("unmarshal req pb msg err: %v", err)
		}
	}

	// 执行业务
	return h.f(dataReq)
}

// AddHandler 为消息添加具体的处理逻辑
func (m *MsgHandle) AddHandler(grp, cmd uint8, req, rsp proto.Message, f HandlerFunc) {
	msgID := genMsgID(grp, cmd)
	// 1 判断当前msgID绑定的处理方法是否已经存在
	if _, ok := m.Apis[msgID]; ok {
		glog.Slog.PanicF(function.WsCtx, "repeated handler", "grp", grp, "cmd", cmd)
	}
	// 2 添加msg与api的绑定关系
	m.Apis[msgID] = &Handler{
		f:   f,
		req: req,
		rsp: rsp,
	}
	// 3 添加返回pb route
	//if rsp != nil {
	//	AddRoute(grp, cmd, rsp)
	//}
}

// Close 关闭
func (m *MsgHandle) Close() error {
	return nil
}

// genMsgID 通过grp cmd生成msgID key
func genMsgID(grp, cmd uint8) uint32 {
	return uint32(grp)*1000 + uint32(cmd)
}

// ------------------------ route ------------------------

//type MessageID struct {
//	grp uint8
//	cmd uint8
//}
//
//var route map[reflect.Type]*MessageID // key:msg type, val: grp+cmd
//
//func init() {
//	route = make(map[reflect.Type]*MessageID)
//}
//
//func AddRoute(grp, cmd uint8, msg proto.Message) {
//	Type := reflect.TypeOf(msg)
//	if _, ok := route[Type]; ok {
//		glog.Slog.Printf("add route had add")
//	}
//	route[Type] = &MessageID{grp: grp, cmd: cmd}
//}
//
//func MsgID(msg proto.Message) (grp uint8, cmd uint8) {
//	Type := reflect.TypeOf(msg)
//	if route[Type] != nil {
//		grp, cmd = route[Type].grp, route[Type].cmd
//	}
//	return
//}

//type MessageHandler interface {
//	GetSeq(context context.Context, data *Req) ([]byte, error)
//	SendMessage(context context.Context, data *Req) ([]byte, error)
//	SendSignalMessage(context context.Context, data *Req) ([]byte, error)
//	PullMessageBySeqList(context context.Context, data *Req) ([]byte, error)
//	UserLogout(context context.Context, data *Req) ([]byte, error)
//	SetUserDeviceBackground(context context.Context, data *Req) ([]byte, bool, error)
//}
//
//var _ MessageHandler = (*GrpcHandler)(nil)
//
//type GrpcHandler struct {
//	msgRpcClient *rpcclient.MessageRpcClient
//	pushClient   *rpcclient.PushRpcClient
//	validate     *validator.Validate
//}
//
//func NewGrpcHandler(validate *validator.Validate, client discovery.SvcDiscoveryRegistry, rpcRegisterName *config.RpcRegisterName) *GrpcHandler {
//	msgRpcClient := rpcclient.NewMessageRpcClient(client, rpcRegisterName.Msg)
//	pushRpcClient := rpcclient.NewPushRpcClient(client, rpcRegisterName.Push)
//	return &GrpcHandler{
//		msgRpcClient: &msgRpcClient,
//		pushClient:   &pushRpcClient, validate: validate,
//	}
//}
//
//func (g GrpcHandler) GetSeq(ctx context.Context, data *Req) ([]byte, error) {
//	req := sdkws.GetMaxSeqReq{}
//	if err := proto.Unmarshal(data.Data, &req); err != nil {
//		return nil, gerr.WrapMsg(err, "GetSeq: error unmarshaling request", "action", "unmarshal", "dataType", "GetMaxSeqReq")
//	}
//	if err := g.validate.Struct(&req); err != nil {
//		return nil, gerr.WrapMsg(err, "GetSeq: validation failed", "action", "validate", "dataType", "GetMaxSeqReq")
//	}
//	resp, err := g.msgRpcClient.GetMaxSeq(ctx, &req)
//	if err != nil {
//		return nil, err
//	}
//	c, err := proto.Marshal(resp)
//	if err != nil {
//		return nil, gerr.WrapMsg(err, "GetSeq: error marshaling response", "action", "marshal", "dataType", "GetMaxSeqResp")
//	}
//	return c, nil
//}
//
//// SendMessage handles the sending of messages through gRPC. It unmarshals the request data,
//// validates the message, and then sends it using the message RPC client.
//func (g GrpcHandler) SendMessage(ctx context.Context, data *Req) ([]byte, error) {
//	var msgData sdkws.MsgData
//	if err := proto.Unmarshal(data.Data, &msgData); err != nil {
//		return nil, gerr.WrapMsg(err, "SendMessage: error unmarshaling message data", "action", "unmarshal", "dataType", "MsgData")
//	}
//
//	if err := g.validate.Struct(&msgData); err != nil {
//		return nil, gerr.WrapMsg(err, "SendMessage: message data validation failed", "action", "validate", "dataType", "MsgData")
//	}
//
//	req := msg.SendMsgReq{MsgData: &msgData}
//	resp, err := g.msgRpcClient.SendMsg(ctx, &req)
//	if err != nil {
//		return nil, err
//	}
//
//	c, err := proto.Marshal(resp)
//	if err != nil {
//		return nil, gerr.WrapMsg(err, "SendMessage: error marshaling response", "action", "marshal", "dataType", "SendMsgResp")
//	}
//
//	return c, nil
//}
//
//func (g GrpcHandler) SendSignalMessage(context context.Context, data *Req) ([]byte, error) {
//	resp, err := g.msgRpcClient.SendMsg(context, nil)
//	if err != nil {
//		return nil, err
//	}
//	c, err := proto.Marshal(resp)
//	if err != nil {
//		return nil, gerr.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "SendMsgResp")
//	}
//	return c, nil
//}
//
//func (g GrpcHandler) PullMessageBySeqList(context context.Context, data *Req) ([]byte, error) {
//	req := sdkws.PullMessageBySeqsReq{}
//	if err := proto.Unmarshal(data.Data, &req); err != nil {
//		return nil, gerr.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "PullMessageBySeqsReq")
//	}
//	if err := g.validate.Struct(data); err != nil {
//		return nil, gerr.WrapMsg(err, "validation failed", "action", "validate", "dataType", "PullMessageBySeqsReq")
//	}
//	resp, err := g.msgRpcClient.PullMessageBySeqList(context, &req)
//	if err != nil {
//		return nil, err
//	}
//	c, err := proto.Marshal(resp)
//	if err != nil {
//		return nil, gerr.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "PullMessageBySeqsResp")
//	}
//	return c, nil
//}
//
//func (g GrpcHandler) UserLogout(context context.Context, data *Req) ([]byte, error) {
//	req := push.DelUserPushTokenReq{}
//	if err := proto.Unmarshal(data.Data, &req); err != nil {
//		return nil, gerr.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "DelUserPushTokenReq")
//	}
//	resp, err := g.pushClient.DelUserPushToken(context, &req)
//	if err != nil {
//		return nil, err
//	}
//	c, err := proto.Marshal(resp)
//	if err != nil {
//		return nil, gerr.WrapMsg(err, "error marshaling response", "action", "marshal", "dataType", "DelUserPushTokenResp")
//	}
//	return c, nil
//}
//
//func (g GrpcHandler) SetUserDeviceBackground(_ context.Context, data *Req) ([]byte, bool, error) {
//	req := sdkws.SetAppBackgroundStatusReq{}
//	if err := proto.Unmarshal(data.Data, &req); err != nil {
//		return nil, false, gerr.WrapMsg(err, "error unmarshaling request", "action", "unmarshal", "dataType", "SetAppBackgroundStatusReq")
//	}
//	if err := g.validate.Struct(data); err != nil {
//		return nil, false, gerr.WrapMsg(err, "validation failed", "action", "validate", "dataType", "SetAppBackgroundStatusReq")
//	}
//	return nil, req.IsBackground, nil
//}
