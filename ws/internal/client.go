package internal

import (
	"context"
	"framework-gin/ws/constant"
	"framework-gin/ws/mcontext"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin"
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
	PlatformID     int    `json:"platformID"`
	IsCompress     bool   `json:"isCompress"`
	UserID         string `json:"userID"`
	userCtx        *UserConnContext
	longConnServer LongConnServer
	closed         atomic.Bool
	closedErr      error
	token          string
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
	c.UserID = ctx.GetUserID()
	c.userCtx = ctx
	c.longConnServer = longConnServer
	c.closed.Store(false)
	c.closedErr = nil
	c.token = ctx.GetToken()
	c.hbCtx, c.hbCancel = context.WithCancel(c.userCtx)
	c.subLock = new(sync.Mutex)
	if c.subUserIDs != nil {
		clear(c.subUserIDs)
	}
	c.subUserIDs = make(map[string]struct{})
	c.logicHandler = NewMsgHandle()
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(pongWait); err != nil {
		return err
	}

	glog.Slog.DebugKVs(c.userCtx.Ctx, "ping Handler Success.", "appData", appData)
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
			glog.Slog.PanicF(c.userCtx.Ctx, "readMessage panic: %s", r)
		}
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(pongWait)
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	c.activeHeartbeat()

	for {
		glog.Slog.DebugKVs(c.userCtx.Ctx, "readMessage")
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			glog.Slog.WarnKVs(c.userCtx.Ctx, "readMessage", "readMessage", returnErr, "messageType", messageType)
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
			glog.Slog.ErrorKVs(c.userCtx.Ctx, "pingHandler", err)
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
		return err
	}

	if err = c.longConnServer.Validate(binaryReq); err != nil {
		return err
	}

	if binaryReq.SendID != c.UserID {
		return gerr.New("exception conn userID not same to req userID", "binaryReq", binaryReq.String())
	}

	ctx := mcontext.WithMustInfoCtx(
		[]string{binaryReq.RequestID, binaryReq.SendID, constant.PlatformIDToName(c.PlatformID), c.userCtx.GetConnID()},
	)

	glog.Slog.DebugKVs(ctx, "gateway req message", "req", binaryReq.String())

	data, err := c.logicHandler.DoMsgHandler(ctx, binaryReq)
	if err != nil {
		return err
	}

	return c.replyMessage(ctx, binaryReq, err, data)
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

func (c *Client) replyMessage(ctx context.Context, binaryReq *Req, err error, data proto.Message) error {
	errResp := ggin.ParseError(err)
	mReply := Resp{
		GrpID:     binaryReq.GrpID,
		CmdID:     binaryReq.CmdID,
		RequestID: binaryReq.RequestID,
		Code:      errResp.Code,
		Msg:       errResp.Msg,
		Data:      data,
	}
	glog.Slog.DebugKVs(ctx, "gateway reply message", "resp", mReply.String())
	err = c.writeBinaryMsg(mReply)
	if err != nil {
		glog.Slog.WarnKVs(ctx, "wireBinaryMsg replyMessage", err, "resp", mReply.String())
	}
	return nil
}

func (c *Client) PushMessage(ctx context.Context, data proto.Message) error {
	resp := Resp{
		//ReqIdentifier: WSPushMsg,
		//TODO 推送消息格式
		RequestID: mcontext.GetRequestID(ctx),
		Data:      data,
	}
	return c.writeBinaryMsg(resp)
}

func (c *Client) KickOnlineMessage() error {
	//TODO 踢人消息格式
	resp := Resp{
		//ReqIdentifier: WSKickOnlineMsg,
	}
	glog.Slog.DebugKVs(c.userCtx.Ctx, "KickOnlineMessage", "resp", resp.String())
	err := c.writeBinaryMsg(resp)
	c.close()
	return err
}

func (c *Client) PushUserOnlineStatus(data proto.Message) error {
	//TODO 推送用户在线格式
	resp := Resp{
		//ReqIdentifier: WsSubUserOnlineStatus,
		Data: data,
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
						glog.Slog.WarnKVs(c.userCtx.Ctx, "send Ping Message error.", "err", err.Error())
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
	glog.Slog.DebugKVs(c.userCtx.Ctx, "write Ping Msg in Server")
	return c.conn.WriteMessage(PingMessage, nil)
}

func (c *Client) writePongMsg(appData string) error {
	glog.Slog.DebugKVs(c.userCtx.Ctx, "write Pong Msg in Server", "appData", appData)
	if c.closed.Load() {
		glog.Slog.WarnKVs(c.userCtx.Ctx, "is closed in server", nil, "appdata", appData, "closed err", c.closedErr)
		return nil
	}

	c.w.Lock()
	defer c.w.Unlock()

	err := c.conn.SetWriteDeadline(writeWait)
	if err != nil {
		glog.Slog.WarnKVs(c.userCtx.Ctx, "SetWriteDeadline in Server have error", gerr.Wrap(err), "writeWait", writeWait, "appData", appData)
		return gerr.Wrap(err)
	}
	err = c.conn.WriteMessage(PongMessage, []byte(appData))
	if err != nil {
		glog.Slog.WarnKVs(c.userCtx.Ctx, "WriteMessage in Server have error", gerr.Wrap(err), "Pong msg", PongMessage, "appData", appData)
	}

	return gerr.Wrap(err)
}
