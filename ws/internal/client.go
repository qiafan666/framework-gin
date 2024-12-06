package internal

import (
	"context"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrConnClosed                = gerr.New("conn has closed")
	ErrNotSupportMessageProtocol = gerr.New("not support message protocol")
	ErrClientClosed              = gerr.New("client actively close the connection")
	ErrPanic                     = gerr.New("panic error")
)

const (
	// MessageText 表示 UTF-8 编码的文本消息，例如 JSON。
	MessageText = iota + 1
	// MessageBinary 表示二进制消息，例如 protobufs。
	MessageBinary
	// CloseMessage 表示关闭控制消息。可选消息负载包含一个数字代码和文本。
	// 使用 FormatCloseMessage 函数格式化关闭消息负载。
	CloseMessage = 8

	// PingMessage 表示 ping 控制消息。可选消息负载为 UTF-8 编码文本。
	PingMessage = 9

	// PongMessage 表示 pong 控制消息。可选消息负载为 UTF-8 编码文本。
	PongMessage = 10
)

type PingPongHandler func(string) error

type Client struct {
	w              *sync.Mutex
	conn           LongConn
	parseToken     *pb.ParseToken
	UserCtx        *UserConnContext
	LongConnServer LongConnServer
	closed         atomic.Bool
	closedErr      error
	hbCtx          context.Context
	hbCancel       context.CancelFunc
	subLock        *sync.Mutex
	subUserIDs     map[string]struct{} // 订阅用户ID集合
	MsgHandle      *MsgHandle
}

// ResetClient updates the client's state with new connection and context information.
func (c *Client) ResetClient(ctx *UserConnContext, conn LongConn, longConnServer LongConnServer) {
	c.w = new(sync.Mutex)
	c.conn = conn
	c.parseToken = &pb.ParseToken{}
	c.UserCtx = ctx
	c.LongConnServer = longConnServer
	c.closed.Store(false)
	c.closedErr = nil
	c.hbCtx, c.hbCancel = context.WithCancel(c.UserCtx)
	c.subLock = new(sync.Mutex)
	if c.subUserIDs != nil {
		clear(c.subUserIDs)
	}
	c.subUserIDs = make(map[string]struct{})
	c.MsgHandle = GetMsgHandler()
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}

	glog.Slog.DebugKVs(c.UserCtx.Ctx, "pingHandler", "appData", appData)
	return c.writePongMsg(appData)
}

func (c *Client) pongHandler(_ string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}
	return nil
}

// readMessage 读取消息
func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			c.closedErr = ErrPanic
			glog.Slog.PanicKVs(c.UserCtx.Ctx, "readMessage", "err", r)
		}
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	c.activeHeartbeat()

	for {
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			glog.Slog.WarnKVs(c.UserCtx.Ctx, "readMessage", "err", returnErr, "messageType", messageType)
			c.closedErr = returnErr
			return
		}

		glog.Slog.DebugKVs(c.UserCtx.Ctx, "readMessage", "messageType", messageType)
		if c.closed.Load() {
			// 连接刚刚关闭，但协程尚未退出的情况
			c.closedErr = ErrConnClosed
			return
		}

		switch messageType {
		case MessageBinary:
			_ = c.conn.SetReadDeadline(pongWait)
			parseDataErr := c.handleMessage(message)
			if parseDataErr != nil {
				c.closedErr = parseDataErr
				return
			}
		case MessageText:
			c.closedErr = ErrNotSupportMessageProtocol
			return
		case PingMessage:
			err := c.writePongMsg("")
			glog.Slog.ErrorKVs(c.UserCtx.Ctx, "readMessage", "writePongMsg err", err)
		case CloseMessage:
			c.closedErr = ErrClientClosed
			return
		default:
		}
	}
}

// handleMessage 处理消息
func (c *Client) handleMessage(message []byte) error {

	if c.UserCtx.IsCompress {
		var err error
		message, err = c.LongConnServer.DecompressWithPool(message)
		if err != nil {
			return gerr.Wrap(err)
		}
	}

	var binaryReq = GetReq()
	defer FreeReq(binaryReq)

	err := c.LongConnServer.Decode(message, binaryReq)
	if err != nil {
		glog.Slog.ErrorKVs(c.UserCtx.Ctx, "handleMessage", "decode error", err)
		return err
	}
	if err = c.LongConnServer.Validate(binaryReq); err != nil {
		glog.Slog.ErrorKVs(c.UserCtx.Ctx, "handleMessage", "validate error", err)
		return err
	}

	c.UserCtx.Ctx = WithMustInfoCtx(
		[]any{constant.PlatformIDToName(c.UserCtx.PlatformID), c.UserCtx.GetConnID(), c.parseToken.UserId, binaryReq.RequestId, c.UserCtx.RemoteAddr},
	)

	glog.Slog.DebugKVs(c.UserCtx.Ctx, "handleMessage", "req", binaryReq.String())
	startTime := time.Now()
	data, code := c.MsgHandle.DoMsgHandler(c, binaryReq)
	glog.Slog.DebugKVs(c.UserCtx.Ctx, "handleMessage DoMsgHandler", "data", data, "code", code, "time", time.Since(startTime))

	return c.replyMessage(c.UserCtx.Ctx, binaryReq, data, code)
}

func (c *Client) close() {
	c.w.Lock()
	defer c.w.Unlock()
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	c.conn.Close()
	c.hbCancel()
	c.LongConnServer.UnRegister(c)
}

func (c *Client) replyMessage(ctx context.Context, binaryReq *Req, data proto.Message, code int) error {
	mReply := Resp{
		GrpID:     binaryReq.GrpId,
		CmdID:     binaryReq.CmdId,
		RequestID: binaryReq.RequestId,
	}

	if code == 0 {
		marshal, err := proto.Marshal(data)
		if err != nil {
			glog.Slog.ErrorKVs(ctx, "replyMessage", "marshal data error", err)
			mReply.Code = gerr.UnKnowError
		} else {
			mReply.Data = marshal
		}
	} else {
		mReply.Code = code
		mReply.Msg = gerr.GetLanguageMsg(code, c.UserCtx.Language)
	}
	glog.Slog.DebugKVs(ctx, "replyMessage", "resp", mReply.String())

	err := c.writeBinaryMsg(mReply)
	if err != nil {
		glog.Slog.WarnKVs(ctx, "replyMessage", "writeBinaryMsg error", err, "resp", mReply.String())
	}
	return nil
}

func (c *Client) writeBinaryMsg(resp Resp) error {
	if c.closed.Load() {
		return nil
	}

	encodedBuf, err := c.LongConnServer.Encode(resp)
	if err != nil {
		return err
	}

	c.w.Lock()
	defer c.w.Unlock()

	err = c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}

	if c.UserCtx.IsCompress {
		resultBuf, compressErr := c.LongConnServer.CompressWithPool(encodedBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	}

	return c.conn.WriteMessage(MessageBinary, encodedBuf)
}

// 在Web平台上主动发起心跳
func (c *Client) activeHeartbeat() {
	if c.UserCtx.PlatformID == constant.WebPlatformID {
		go func() {
			glog.Slog.DebugKVs(c.UserCtx.Ctx, "activeHeartbeat start.")
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := c.writePingMsg(); err != nil {
						glog.Slog.WarnKVs(c.UserCtx.Ctx, "activeHeartbeat", "err", err)
						return
					}
				case <-c.hbCtx.Done():
					return
				}
			}
		}()
	}
}

func (c *Client) writePingMsg() error {
	if c.closed.Load() {
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}
	//glog.Slog.DebugKVs(c.userCtx.Ctx, "writePingMsg")
	return c.conn.WriteMessage(PingMessage, nil)
}

func (c *Client) writePongMsg(appData string) error {
	//glog.Slog.DebugKVs(c.userCtx.Ctx, "writePongMsg", "appData", appData)
	if c.closed.Load() {
		glog.Slog.WarnKVs(c.UserCtx.Ctx, "writePongMsg", "appdata", appData, "closed err", c.closedErr)
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		glog.Slog.WarnKVs(c.UserCtx.Ctx, "writePongMsg", "SetWriteDeadline in Server have error", gerr.Wrap(err), "writeWait", writeWait, "appData", appData)
		return gerr.Wrap(err)
	}
	err = c.conn.WriteMessage(PongMessage, []byte(appData))
	if err != nil {
		glog.Slog.WarnKVs(c.UserCtx.Ctx, "writePongMsg", "WriteMessage in Server have error", gerr.Wrap(err), "Pong msg", PongMessage, "appData", appData)
	}

	return gerr.Wrap(err)
}

// ------------------------ inner function ------------------------

// PubMessage 推送消息,消息进入redis pub/sub,由redis订阅者进行处理
func (c *Client) PubMessage(ctx context.Context, pubsub *pb.ReqPushMsgToOther) error {
	return c.LongConnServer.PubMsgChannel(ctx, pubsub)
}

// PushMessage 服务器主动向客户端推送消息
func (c *Client) PushMessage(req *pb.ReqPushMsgToOther) error {
	resp := Resp{
		GrpID: uint8(req.GrpId),
		CmdID: uint8(req.CmdId),
		Data:  req.Data,
	}
	glog.Slog.DebugKVs(c.UserCtx.Ctx, "PushMessage", "resp", resp.String())
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

// KickOnlineMessage 踢下线 分布式使用
func (c *Client) KickOnlineMessage(reason pb.KickReason) error {

	pbRsp := &pb.RpcUserKickOff{
		Reason: reason,
	}
	protoData, err := proto.Marshal(pbRsp)
	if err != nil {
		glog.Slog.ErrorKVs(c.UserCtx.Ctx, "KickOnlineMessage", "marshal data error", err)
		return err
	}
	resp := Resp{
		GrpID: uint8(pb.Grp_Sys),
		CmdID: uint8(pb.CmdSys_KickOnlineUser),
		Data:  protoData,
	}
	glog.Slog.DebugKVs(c.UserCtx.Ctx, "KickOnlineMessage", "resp", resp.String())
	err = c.writeBinaryMsg(resp)
	c.close()
	return err
}

// PushUserOnlineStatus 推送用户在线状态
func (c *Client) PushUserOnlineStatus(data []byte) error {
	resp := Resp{
		GrpID: uint8(pb.Grp_Sys),
		CmdID: uint8(pb.CmdSys_SubscribeOnlineUser),
		Data:  data,
	}
	return c.writeBinaryMsg(resp)
}
