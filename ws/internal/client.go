package internal

import (
	"context"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gcast"
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
	PlatformID     int  `json:"platform_id"`
	IsCompress     bool `json:"is_compress"`
	parseToken     *pb.ParseToken
	userCtx        *UserConnContext
	longConnServer LongConnServer
	closed         atomic.Bool
	closedErr      error
	hbCtx          context.Context
	hbCancel       context.CancelFunc
	subLock        *sync.Mutex
	subUserIDs     map[string]struct{} // 订阅用户ID集合
	logicHandler   *MsgHandle
}

// ResetClient updates the client's state with new connection and context information.
func (c *Client) ResetClient(ctx *UserConnContext, conn LongConn, longConnServer LongConnServer) {
	c.w = new(sync.Mutex)
	c.conn = conn
	c.PlatformID = gcast.ToInt(ctx.GetPlatformID())
	c.IsCompress = ctx.GetCompression()
	c.parseToken = &pb.ParseToken{}
	c.userCtx = ctx
	c.longConnServer = longConnServer
	c.closed.Store(false)
	c.closedErr = nil
	c.hbCtx, c.hbCancel = context.WithCancel(c.userCtx)
	c.subLock = new(sync.Mutex)
	if c.subUserIDs != nil {
		clear(c.subUserIDs)
	}
	c.subUserIDs = make(map[string]struct{})
	c.logicHandler = GetMsgHandler()
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}

	glog.Slog.DebugKVs(c.userCtx.Ctx, "pingHandler", "appData", appData)
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
			glog.Slog.PanicKVs(c.userCtx.Ctx, "readMessage", "err", r)
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
			glog.Slog.WarnKVs(c.userCtx.Ctx, "readMessage", "err", returnErr, "messageType", messageType)
			c.closedErr = returnErr
			return
		}

		glog.Slog.DebugKVs(c.userCtx.Ctx, "readMessage", "messageType", messageType)
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
			glog.Slog.ErrorKVs(c.userCtx.Ctx, "readMessage", "writePongMsg err", err)
		case CloseMessage:
			c.closedErr = ErrClientClosed
			return
		default:
		}
	}
}

// handleMessage 处理消息
func (c *Client) handleMessage(message []byte) error {

	if c.IsCompress {
		var err error
		message, err = c.longConnServer.DecompressWithPool(message)
		if err != nil {
			return gerr.Wrap(err)
		}
	}

	var binaryReq = getReq()
	defer freeReq(binaryReq)

	err := c.longConnServer.Decode(message, binaryReq)
	if err != nil {
		glog.Slog.ErrorKVs(c.userCtx.Ctx, "handleMessage", "decode error", err)
		return err
	}
	if err = c.longConnServer.Validate(binaryReq); err != nil {
		glog.Slog.ErrorKVs(c.userCtx.Ctx, "handleMessage", "validate error", err)
		return err
	}

	c.userCtx.Ctx = WithMustInfoCtx(
		[]any{constant.PlatformIDToName(c.PlatformID), c.userCtx.GetConnID(), c.parseToken.UserId, binaryReq.RequestID, c.userCtx.RemoteAddr},
	)

	glog.Slog.DebugKVs(c.userCtx.Ctx, "handleMessage", "req", binaryReq.String())

	var data proto.Message
	var code int
	if binaryReq.GrpID == uint8(pb.Grp_Sys) {
		switch binaryReq.CmdID {
		case uint8(pb.CmdSys_Subscribe_Online_User):
			data, code = c.longConnServer.SubUserOnlineStatus(c.userCtx.Ctx, c, binaryReq)
		default:
			return gerr.New("not support groupID and cmdID", "groupID", binaryReq.GrpID, "cmdID", binaryReq.CmdID)
		}
	} else {
		data, code = c.logicHandler.DoMsgHandler(c.userCtx.Ctx, binaryReq)
		glog.Slog.DebugKVs(c.userCtx.Ctx, "handleMessage DoMsgHandler", "err",
			gerr.GetCodeAndMsg(code, c.userCtx.Language))
	}

	return c.replyMessage(c.userCtx.Ctx, binaryReq, data, code)
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
	c.longConnServer.UnRegister(c)
}

func (c *Client) replyMessage(ctx context.Context, binaryReq *Req, data proto.Message, code int) error {
	glog.Slog.DebugKVs(ctx, "replyMessage", "data", data, "code", code)

	mReply := Resp{
		GrpID:     binaryReq.GrpID,
		CmdID:     binaryReq.CmdID,
		RequestID: binaryReq.RequestID,
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
		mReply.Msg = gerr.GetCodeAndMsg(code, c.userCtx.Language)
	}
	glog.Slog.DebugKVs(ctx, "replyMessage", "resp", mReply.String())

	err := c.writeBinaryMsg(mReply)
	if err != nil {
		glog.Slog.WarnKVs(ctx, "replyMessage", "writeBinaryMsg error", err, "resp", mReply.String())
	}
	return nil
}

func (c *Client) PushMessage(ctx context.Context, data []byte) error {
	resp := Resp{
		//ReqIdentifier: WSPushMsg,
		//TODO 推送消息格式
		Data: data,
	}
	return c.writeBinaryMsg(resp)
}

// KickOnlineMessage 踢下线 分布式使用
func (c *Client) KickOnlineMessage() error {
	resp := Resp{
		GrpID: uint8(pb.Grp_Sys),
		CmdID: uint8(pb.CmdSys_Kick_Online_User),
	}
	glog.Slog.DebugKVs(c.userCtx.Ctx, "KickOnlineMessage", "resp", resp.String())
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

func (c *Client) PushUserOnlineStatus(data []byte) error {
	resp := Resp{
		GrpID: uint8(pb.Grp_Sys),
		CmdID: uint8(pb.CmdSys_Subscribe_Online_User),
		Data:  data,
	}
	return c.writeBinaryMsg(resp)
}

func (c *Client) writeBinaryMsg(resp Resp) error {
	if c.closed.Load() {
		return nil
	}

	encodedBuf, err := c.longConnServer.Encode(resp)
	if err != nil {
		return err
	}

	c.w.Lock()
	defer c.w.Unlock()

	err = c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		return err
	}

	if c.IsCompress {
		resultBuf, compressErr := c.longConnServer.CompressWithPool(encodedBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(MessageBinary, resultBuf)
	}

	return c.conn.WriteMessage(MessageBinary, encodedBuf)
}

// 在Web平台上主动发起心跳
func (c *Client) activeHeartbeat() {
	if c.PlatformID == constant.WebPlatformID {
		go func() {
			glog.Slog.DebugKVs(c.userCtx.Ctx, "activeHeartbeat start.")
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := c.writePingMsg(); err != nil {
						glog.Slog.WarnKVs(c.userCtx.Ctx, "activeHeartbeat", "err", err)
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
		glog.Slog.WarnKVs(c.userCtx.Ctx, "writePongMsg", "appdata", appData, "closed err", c.closedErr)
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		glog.Slog.WarnKVs(c.userCtx.Ctx, "writePongMsg", "SetWriteDeadline in Server have error", gerr.Wrap(err), "writeWait", writeWait, "appData", appData)
		return gerr.Wrap(err)
	}
	err = c.conn.WriteMessage(PongMessage, []byte(appData))
	if err != nil {
		glog.Slog.WarnKVs(c.userCtx.Ctx, "writePongMsg", "WriteMessage in Server have error", gerr.Wrap(err), "Pong msg", PongMessage, "appData", appData)
	}

	return gerr.Wrap(err)
}
