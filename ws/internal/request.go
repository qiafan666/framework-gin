package internal

import (
	"github.com/golang/protobuf/proto"
)

type IRequest interface {
	GetConnection() LongConn
	ProtoMessage() proto.Message
	SetProtoMessage(proto.Message)
	GetGrp() uint8
	GetCmd() uint8
	GetRequestID() string
	GetData() []byte
}

// Request 请求
type Request struct {
	conn     *GWebSocket   // 已经和客户端建立好的链接
	msg      *Req          // 客户端请求的数据
	protoMsg proto.Message // 解析的pb消息
}

// GetConnection 获取conn连接
func (r *Request) GetConnection() LongConn {
	return r.conn
}

// GetGrp 获取Grp
func (r *Request) GetGrp() uint8 {
	return r.msg.GrpID
}

// GetCmd 获取Cmd
func (r *Request) GetCmd() uint8 {
	return r.msg.CmdID
}

// GetData 获取tcp消息数据
func (r *Request) GetData() []byte {
	return r.msg.Data
}

// GetOperateID 获取操作ID
func (r *Request) GetRequestID() string {
	return r.msg.RequestID
}

// ProtoMessage 获取pb协议消息
func (r *Request) ProtoMessage() proto.Message {
	return r.protoMsg
}

// ProtoMessage 设置消息
func (r *Request) SetProtoMessage(pbMsg proto.Message) {
	r.protoMsg = pbMsg
}
