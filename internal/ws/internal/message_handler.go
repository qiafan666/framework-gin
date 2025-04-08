package internal

import (
	"encoding/json"
	"framework-gin/internal/ws/constant"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gerr"
	"google.golang.org/protobuf/proto"
	"sync"
)

// -------------------- req and resp --------------------

type Req struct {
	RequestId string `json:"request_id"   validate:"required"`
	GrpId     uint8  `json:"grp_id" validate:"required"` // 消息组id
	CmdId     uint8  `json:"cmd_id" validate:"required"` // 消息的ID
	Data      []byte `json:"data"`
}

func (r *Req) String() string {
	var tReq Req
	tReq.RequestId = r.RequestId
	tReq.GrpId = r.GrpId
	tReq.CmdId = r.CmdId
	tReq.Data = r.Data
	return gcast.ToString(tReq)
}

var reqPool = sync.Pool{
	New: func() any {
		return new(Req)
	},
}

func GetReq() *Req {
	req := reqPool.Get().(*Req)
	req.Data = nil
	req.RequestId = ""
	req.GrpId = 0
	req.CmdId = 0
	return req
}

func FreeReq(req *Req) {
	reqPool.Put(req)
}

type Resp struct {
	GrpId     uint8  `json:"grp_id"`
	CmdId     uint8  `json:"cmd_id"`
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	Data      []byte `json:"data"`
}

func (r *Resp) String() string {
	var tResp Resp
	tResp.RequestId = r.RequestId
	tResp.Code = r.Code
	tResp.Msg = r.Msg
	tResp.GrpId = r.GrpId
	tResp.CmdId = r.CmdId
	tResp.Data = r.Data
	return gcast.ToString(tResp)
}

var respPool = sync.Pool{
	New: func() any {
		return new(Resp)
	},
}

func GetResp() *Resp {
	resp := respPool.Get().(*Resp)
	resp.GrpId = 0
	resp.CmdId = 0
	resp.RequestId = ""
	resp.Code = 0
	resp.Msg = ""
	resp.Data = nil
	return resp
}

func FreeResp(resp *Resp) {
	respPool.Put(resp)
}

// -------------------- msg handler --------------------

var handler *MsgHandle

// HandlerFunc 消息处理函数
type HandlerFunc func(client *Client, req proto.Message) (proto.Message, int)

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
func (m *MsgHandle) DoMsgHandler(client *Client, req *Req) (proto.Message, int) {
	msgID := genMsgID(req.GrpId, req.CmdId)
	h, ok := m.Apis[msgID]
	if !ok {
		client.logger.ErrorKVs(client.UserCtx.TraceCtx, "DoMsgHandler", "msgID not found,MsgID", msgID)
		return nil, gerr.UnKnowError
	}

	// 解析pb消息
	var dataReq proto.Message
	if config.Ws.Protocol == constant.ProtocolJson {
		if h.req != nil {
			dataReq = proto.Clone(h.req)
			if err := json.Unmarshal(req.Data, dataReq); err != nil {
				client.logger.ErrorKVs(client.UserCtx.TraceCtx, "DoMsgHandler", "unmarshal req pb msg err", err)
				return nil, gerr.UnKnowError
			}
		}
	} else {
		if h.req != nil {
			dataReq = proto.Clone(h.req)
			if err := proto.Unmarshal(req.Data, dataReq); err != nil {
				client.logger.ErrorKVs(client.UserCtx.TraceCtx, "DoMsgHandler", "unmarshal req pb msg err", err)
				return nil, gerr.UnKnowError
			}
		}
	}

	// 执行业务
	return h.f(client, dataReq)
}

// AddHandler 为消息添加具体的处理逻辑
func (m *MsgHandle) AddHandler(grp, cmd uint8, req, rsp proto.Message, f HandlerFunc) {
	msgID := genMsgID(grp, cmd)
	if _, ok := m.Apis[msgID]; !ok {
		m.Apis[msgID] = &Handler{
			f:   f,
			req: req,
			rsp: rsp,
		}
	}
}

// GetPbReq 获取请求pb
func (m *MsgHandle) GetPbReq(grpID, cmdID uint8) (proto.Message, int) {
	genID := genMsgID(grpID, cmdID)
	h, ok := m.Apis[genID]
	if ok {
		return h.req, gerr.OK
	}
	return nil, gerr.ParameterError
}

// GetPbRsp 获取返回pb
func (m *MsgHandle) GetPbRsp(msgID, grpID uint8) (proto.Message, int) {
	genID := genMsgID(grpID, msgID)
	h, ok := m.Apis[genID]
	if ok {
		return h.rsp, gerr.OK
	}
	return nil, gerr.ParameterError
}

// 根据
// genMsgID 通过grp cmd生成msgID key
func genMsgID(grp, cmd uint8) uint32 {
	return uint32(grp)*1000 + uint32(cmd)
}
