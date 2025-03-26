package internal

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gerr"
	"google.golang.org/protobuf/proto"
	"sync"
	"sync/atomic"
	"time"
)

type PingPongHandler func(string) error

type Client struct {
	lock       *sync.Mutex
	conn       ConnInterface
	server     ServerInterface
	parseToken *pb.ParseToken
	UserCtx    *UserConnContext
	closed     atomic.Bool
	closedErr  error
	msgHandle  *MsgHandle
	logger     *Logger
}

// ResetClient updates the client's state with new connection and context information.
func (c *Client) ResetClient(ctx *UserConnContext, connInterface ConnInterface, serverInterface ServerInterface) {
	c.lock = new(sync.Mutex)
	c.conn = connInterface
	c.server = serverInterface
	c.parseToken = &pb.ParseToken{}
	c.UserCtx = ctx
	c.closed.Store(false)
	c.closedErr = nil
	c.msgHandle = GetMsgHandler()
}

// GetClientState 获取当前连接是否为私有连接，true为私有连接，false为公共连接
func (c *Client) GetClientState() bool {
	if c.parseToken.UserId != "" {
		return true
	}
	return false
}

func (c *Client) GetClientLiveTime() time.Duration {
	if c.GetClientState() {
		return time.Duration(config.Ws.PrivateLiveTime) * time.Second
	} else {
		return time.Duration(config.Ws.PublicLiveTime) * time.Second
	}
}

func (c *Client) pingHandler(appData string) error {
	if err := c.conn.SetReadDeadline(c.GetClientLiveTime()); err != nil {
		return err
	}

	c.logger.DebugKVs(c.UserCtx.Trace(), "pingHandler", "appData", appData)
	return c.writePongMsg(appData)
}

func (c *Client) pongHandler(_ string) error {
	if err := c.conn.SetReadDeadline(c.GetClientLiveTime()); err != nil {
		return err
	}
	return nil
}

// readMessage 读取消息
func (c *Client) readMessage() {
	defer func() {
		if r := recover(); r != nil {
			c.closedErr = gerr.New("panic error")
			c.logger.PanicKVs(c.UserCtx.Trace(), "readMessage", "err", r)
		}
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(c.GetClientLiveTime())
	c.conn.SetPongHandler(c.pongHandler)
	c.conn.SetPingHandler(c.pingHandler)
	//c.activeHeartbeat()

	for {
		messageType, message, returnErr := c.conn.ReadMessage()
		if returnErr != nil {
			c.logger.WarnKVs(c.UserCtx.Trace(), "readMessage", "err", returnErr, "messageType", messageType)
			c.closedErr = returnErr
			return
		}

		c.logger.DebugKVs(c.UserCtx.Trace(), "readMessage", "messageType", messageType)
		if c.closed.Load() {
			// 连接刚刚关闭，但协程尚未退出的情况
			c.closedErr = gerr.New("conn has closed")
			return
		}

		switch messageType {
		case MessageBinary, MessageText:
			_ = c.conn.SetReadDeadline(c.GetClientLiveTime())
			parseDataErr := c.handleMessage(message)
			if parseDataErr != nil {
				c.closedErr = parseDataErr
				return
			}
		//case MessageText:
		//	c.closedErr = gerr.New("not support message protocol")
		//	return
		case PingMessage:
			err := c.writePongMsg("")
			if err != nil {
				c.logger.ErrorKVs(c.UserCtx.Trace(), "readMessage", "writePongMsg err", err)
			}
		case CloseMessage:
			c.closedErr = gerr.New("client actively close the connection")
			return
		default:
		}
	}
}

// handleMessage 处理消息
func (c *Client) handleMessage(message []byte) error {

	if c.UserCtx.IsCompress && config.Ws.Protocol == constant.ProtocolJson {

		var err error
		message, err = base64.StdEncoding.DecodeString(string(message))
		if err != nil {
			return err
		}

		message, err = c.server.DecompressWithPool(message)
		if err != nil {
			return gerr.Wrap(err)
		}
	}

	var binaryReq = GetReq()
	defer FreeReq(binaryReq)

	if config.Ws.Protocol == constant.ProtocolJson {
		err := json.Unmarshal(message, binaryReq)
		if err != nil {
			c.logger.ErrorKVs(c.UserCtx.Trace(), "handleMessage", "json unmarshal error", err)
			return err
		}
	} else {
		err := c.server.Decode(message, binaryReq)
		if err != nil {
			c.logger.ErrorKVs(c.UserCtx.Trace(), "handleMessage", "decode error", err)
			return err
		}
	}

	if err := c.server.Validate(binaryReq); err != nil {
		c.logger.ErrorKVs(c.UserCtx.Trace(), "handleMessage", "validate error", err)
		return err
	}

	c.UserCtx.BaseReq = BaseReq{
		RequestId: binaryReq.RequestId,
		GrpId:     binaryReq.GrpId,
		CmdId:     binaryReq.CmdId,
	}

	startTime := time.Now()
	data, code := c.msgHandle.DoMsgHandler(c, binaryReq)
	if config.Ws.Protocol == constant.ProtocolJson {
		c.logger.DebugKVs(c.UserCtx.Trace(), "handleMessage DoMsgHandler", "data", string(binaryReq.Data), "code", code, "time", time.Since(startTime).Milliseconds())
	} else {
		c.logger.DebugKVs(c.UserCtx.Trace(), "handleMessage DoMsgHandler", "data", data, "code", code, "time", time.Since(startTime).Milliseconds())
	}

	return c.replyMessage(c.UserCtx.Trace(), binaryReq, data, code)
}

func (c *Client) close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	c.conn.Close()
	c.server.UnRegister(c)
}

func (c *Client) buildData(data proto.Message) ([]byte, error) {
	if config.Ws.Protocol == constant.ProtocolJson {
		return json.Marshal(data)
	} else {
		return c.server.Encode(data)
	}
}

func (c *Client) buildResp(resp *Resp) ([]byte, error) {
	if config.Ws.Protocol == constant.ProtocolJson {
		return json.Marshal(resp)
	} else {
		return c.server.Encode(resp)
	}
}

func (c *Client) replyMessage(ctx context.Context, binaryReq *Req, data proto.Message, code int) error {
	reply := GetResp()
	defer FreeResp(reply)

	reply.GrpId = binaryReq.GrpId
	reply.CmdId = binaryReq.CmdId
	reply.RequestId = binaryReq.RequestId

	if code == 0 {
		marshal, err := c.buildData(data)
		if err != nil {
			c.logger.ErrorKVs(ctx, "replyMessage", "marshal data error", err)
			reply.Code = gerr.UnKnowError
		} else {
			reply.Data = marshal
		}
	} else {
		reply.Code = code
		reply.Msg = gerr.GetLanguageMsg(code, c.UserCtx.Language)
	}
	c.logger.DebugKVs(ctx, "replyMessage", "resp", reply.String())

	err := c.writeRespMsg(reply)
	if err != nil {
		c.logger.WarnKVs(ctx, "replyMessage", "writeBinaryMsg error", err, "resp", reply.String())
	}
	return nil
}

func (c *Client) writeRespMsg(resp *Resp) error {
	if c.closed.Load() {
		return nil
	}

	if c.UserCtx.RequestId == resp.RequestId {
		c.UserCtx.BaseReq = BaseReq{}
	}

	encodedBuf, err := c.buildResp(resp)
	if err != nil {
		return err
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	err = c.conn.SetWriteDeadline(c.GetClientLiveTime())
	if err != nil {
		return err
	}

	if c.UserCtx.IsCompress && config.Ws.Protocol == constant.ProtocolJson {
		resultBuf, compressErr := c.server.CompressWithPool(encodedBuf)
		if compressErr != nil {
			return compressErr
		}
		return c.conn.WriteMessage(config.Ws.Protocol, []byte(base64.StdEncoding.EncodeToString(resultBuf)))
	}

	return c.conn.WriteMessage(config.Ws.Protocol, encodedBuf)
}

// 在Web平台上主动发起心跳
func (c *Client) activeHeartbeat() {
	if c.UserCtx.PlatformID == constant.WebPlatformID {
		go func() {
			c.logger.DebugKVs(c.UserCtx.Trace(), "activeHeartbeat start.")
			ticker := time.NewTicker(pingPeriod)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					if err := c.writePingMsg(); err != nil {
						c.logger.WarnKVs(c.UserCtx.Trace(), "activeHeartbeat", "err", err)
						return
					}
				}
			}
		}()
	}
}

func (c *Client) writePingMsg() error {
	if c.closed.Load() {
		return nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.conn.SetWriteDeadline(c.GetClientLiveTime())
	if err != nil {
		return err
	}
	//c.logger.DebugKVs(c.userCtx.Ctx, "writePingMsg")
	return c.conn.WriteMessage(PingMessage, nil)
}

func (c *Client) writePongMsg(appData string) error {
	c.logger.DebugKVs(c.UserCtx.Trace(), "writePongMsg", "appData", appData)
	if c.closed.Load() {
		c.logger.WarnKVs(c.UserCtx.Trace(), "writePongMsg", "appdata", appData, "closed err", c.closedErr)
		return nil
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.conn.SetWriteDeadline(c.GetClientLiveTime())
	if err != nil {
		c.logger.WarnKVs(c.UserCtx.Trace(), "writePongMsg", "SetWriteDeadline in Server have error", gerr.Wrap(err), "writeWait", c.GetClientLiveTime(), "appData", appData)
		return gerr.Wrap(err)
	}
	err = c.conn.WriteMessage(PongMessage, []byte(appData))
	if err != nil {
		c.logger.WarnKVs(c.UserCtx.Trace(), "writePongMsg", "WriteMessage in Server have error", gerr.Wrap(err), "Pong msg", PongMessage, "appData", appData)
		return gerr.Wrap(err)
	}

	return nil
}

// -------------------- inner function --------------------

// KickOnlineMessage 踢下线 分布式使用
func (c *Client) KickOnlineMessage(reason pb.TypeKickReason) error {

	pbRsp := &pb.RspUserKickOff{
		Reason: reason,
	}
	data, err := c.buildData(pbRsp)
	if err != nil {
		c.logger.ErrorKVs(c.UserCtx.Trace(), "KickOnlineMessage", "marshal data error", err)
		return err
	}
	resp := GetResp()
	defer FreeResp(resp)

	resp.GrpId = uint8(pb.Grp_Sys)
	resp.CmdId = uint8(pb.CmdSys_KickOnlineUser)
	resp.Data = data

	c.logger.DebugKVs(c.UserCtx.Trace(), "KickOnlineMessage", "resp", resp.String())
	err = c.writeRespMsg(resp)
	c.close()
	return err
}

func (c *Client) SendTickerTestMsg() {

	tickerMsg := &pb.RspTickerSubscribe{
		Msg: "hello",
	}
	data, err := c.buildData(tickerMsg)
	if err != nil {
		c.logger.ErrorKVs(c.UserCtx.Trace(), "SendTickerTestMsg", "build data error", err)
		return
	}

	resp := GetResp()
	defer FreeResp(resp)

	resp.GrpId = uint8(pb.Grp_Public)
	resp.CmdId = uint8(pb.CmdPublic_TickerSubscribe)
	resp.Data = data

	c.logger.DebugKVs(c.UserCtx.Trace(), "SendTickerTestMsg", "resp", resp.String())
	_ = c.writeRespMsg(resp)
}
