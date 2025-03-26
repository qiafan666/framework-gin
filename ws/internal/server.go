package internal

import (
	"framework-gin/common"
	"framework-gin/common/errs"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gcompress"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin"
	"github.com/qiafan666/gotato/commons/glog"
	"sync"
	"sync/atomic"
	"time"
)

type ServerInterface interface {
	Validate(s any) error
	UnRegister(c *Client)
	gcompress.ICompressor
	gcommon.IEncoder
}

type Server struct {
	port             int
	wsMaxConnNum     int64
	registerChan     chan *Client
	unregisterChan   chan *Client
	clients          clientManagerInterface
	clientPool       sync.Pool
	statelessConnNum atomic.Int64
	statefulConnNum  atomic.Int64
	handshakeTimeout time.Duration
	writeBufferSize  int
	logger           *Logger
	gcompress.ICompressor
	gcommon.IEncoder
}

func NewWsServer() *Server {

	server := &Server{
		wsMaxConnNum:     100000,
		writeBufferSize:  4096,
		handshakeTimeout: 10 * time.Second,
		clientPool: sync.Pool{
			New: func() any {
				return new(Client)
			},
		},
		registerChan:   make(chan *Client, 1000),
		unregisterChan: make(chan *Client, 1000),
		clients:        newClientManager(32),
		ICompressor:    gcompress.NewGzipCompressor(),
		IEncoder:       gcommon.NewGobEncoder(),
		logger:         NewLogger("ws", glog.ZapLog),
	}
	return server
}

func (s *Server) Validate(data any) error {
	return gcommon.Validate(data)
}

func (s *Server) UnRegister(c *Client) {
	s.unregisterChan <- c
}

// ------------------------ inner function ------------------------

func (s *Server) Run(r *gin.Engine) {
	var (
		client *Client
	)

	go func() {
		for {
			select {
			case client = <-s.registerChan:
				s.registerClient(client)
			case client = <-s.unregisterChan:
				s.unregisterClient(client)
			}
		}
	}()
	r.GET("/stream", s.wsHandler)
}

func (s *Server) wsHandler(c *gin.Context) {

	// 创建一个新的连接上下文
	connCtx := newContext(c.Writer, c.Request)

	// 检查当前在线用户连接数是否超过最大限制
	if s.statefulConnNum.Load()+s.statelessConnNum.Load() >= s.wsMaxConnNum {
		// 如果超过最大连接数限制，通过HTTP返回错误并停止处理
		ggin.HttpError(connCtx.RespWriter, gerr.NewLang(errs.ConnOverMaxNumLimit, connCtx.Language))
		return
	}

	// 解析必要的参数（例如用户ID、Token）
	err := connCtx.ParseEssentialArgs()
	if err != nil {
		// 如果解析过程中出错，通过HTTP返回错误并停止处理
		ggin.HttpError(connCtx.RespWriter, err)
		return
	}

	// 如果Token解析失败，根据上下文标志决定是否通过WebSocket发送错误信息
	shouldSendResp := connCtx.ShouldSendResp()
	// 调用认证客户端，解析从上下文中获取的Token
	parseToken, err := s.ParseToken(connCtx)
	if err != nil {
		if !shouldSendResp {
			// 创建一个WebSocket连接对象并尝试通过WebSocket发送错误信息
			wsLongConn := newGWebSocket(s.handshakeTimeout, s.writeBufferSize)
			if err = wsLongConn.RespondWithError(err, c.Writer, c.Request); err == nil {
				// 如果错误信息成功通过WebSocket发送，则停止处理
				return
			}
		}
		// 如果不需要或无法通过WebSocket发送，返回HTTP错误并停止处理
		ggin.HttpError(connCtx.RespWriter, err)
		return
	}

	// 创建WebSocket长连接对象
	wsLongConn := newGWebSocket(s.handshakeTimeout, s.writeBufferSize)
	if err = wsLongConn.GenerateLongConn(c.Writer, c.Request); err != nil {
		// 如果长连接创建失败，握手过程中已处理错误
		s.logger.WarnF(connCtx.TraceCtx, "long connection generate failed, err : %v", err)
		return
	} else {
		// 检查是否需要通过WebSocket发送正常响应
		if !shouldSendResp {
			// 尝试通过WebSocket发送成功信息
			if err = wsLongConn.RespondWithSuccess(); err != nil {
				return
			}
		}
	}

	// 从客户端池中获取一个客户端对象，重置其状态，并将其与当前的WebSocket长连接关联
	client := s.clientPool.Get().(*Client)
	client.ResetClient(connCtx, wsLongConn, s)
	client.parseToken = parseToken
	client.logger = s.logger
	if client.GetClientState() {
		client.UserCtx.TraceCtx = SetTraceCtx(
			[]any{constant.PlatformIDToName(gcast.ToInt(client.UserCtx.Req.Header.Get(common.HeaderPlatformID))),
				connCtx.ConnID, connCtx.RemoteAddr, parseToken.UserId})
	} else {
		client.UserCtx.TraceCtx = SetTraceCtx(
			[]any{constant.PlatformIDToName(gcast.ToInt(client.UserCtx.Req.Header.Get(common.HeaderPlatformID))),
				connCtx.ConnID, connCtx.RemoteAddr, client.parseToken.Uuid})
	}

	// 将客户端注册到服务器并开始消息处理
	s.registerChan <- client
	go client.readMessage()
}

func (s *Server) ParseToken(userConCtx *UserConnContext) (*pb.ParseToken, error) {

	return &pb.ParseToken{
		UserId: "",
		Uuid:   userConCtx.Uuid,
	}, nil
}

// registerClient 注册客户端
func (s *Server) registerClient(client *Client) {
	var (
		clientOK   bool // 当前client是否已存在，只存在有状态连接
		oldClients []*Client
		state      bool // 是否状态连接 true:有状态连接 false:无状态连接
	)
	oldClients, clientOK, state = s.clients.GetOldClients(client)
	if len(oldClients) == 0 {
		s.logger.DebugKVs(client.UserCtx.TraceCtx, "registerClient,user not exist", "userID", client.parseToken.UserId, "uuid", client.parseToken.Uuid, "platformID", client.UserCtx.PlatformID)
	} else {
		s.multiTerminalLoginChecker(clientOK, oldClients, client)
		s.logger.DebugKVs(client.UserCtx.TraceCtx, "registerClient,user exist", "userID", client.parseToken.UserId, "uuid", client.parseToken.Uuid, "platformID", client.UserCtx.PlatformID)
		if clientOK {
			// 当前平台连接已经存在，增加连接数
			s.logger.InfoKVs(client.UserCtx.TraceCtx, "registerClient,repeat login", "userID", client.parseToken.UserId, "uuid", client.parseToken.Uuid, "platformID",
				client.UserCtx.PlatformID, "old remote addr", s.getRemoteAdders(oldClients))
		}
	}
	s.clients.Set(client)
	if state {
		s.statefulConnNum.Add(1)
	} else {
		s.statelessConnNum.Add(1)
	}

	s.logger.InfoKVs(
		client.UserCtx.TraceCtx,
		"registerClient,user online",
		"statefulConnNum",
		s.statefulConnNum.Load(),
		"statelessConnNum",
		s.statelessConnNum.Load(),
	)
}

// multiTerminalLoginChecker 多端登录检查
func (s *Server) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	switch config.Ws.MulitLoginPolicy {
	case constant.DefaultNotKick:
	case constant.PCAndOther:
		if constant.PlatformIDToClass(newClient.UserCtx.PlatformID) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if !clientOK {
			return
		}
		for _, c := range oldClients {
			err := c.KickOnlineMessage(pb.TypeKickReason_OnlyOneClient)
			if err != nil {
				s.logger.WarnKVs(c.UserCtx.Trace(), "multiTerminalLoginChecker,KickOnlineMessage", "err", err)
			}
		}
	}
}

func (s *Server) getRemoteAdders(client []*Client) string {
	var ret string
	for i, c := range client {
		if i == 0 {
			ret = c.UserCtx.GetRemoteAddr()
		} else {
			ret += "@" + c.UserCtx.GetRemoteAddr()
		}
	}
	return ret
}

func (s *Server) unregisterClient(client *Client) {
	defer s.clientPool.Put(client)
	isDeleteUser := s.clients.DeleteClients([]*Client{client})

	if isDeleteUser {
		if client.GetClientState() {
			s.statefulConnNum.Add(-1)
		} else {
			s.statelessConnNum.Add(-1)
		}
	}
	s.logger.WarnKVs(client.UserCtx.TraceCtx, "unregisterClient user offline", "close reason", client.closedErr, "statefulConnNum",
		s.statefulConnNum.Load(), "statelessConnNum", s.statelessConnNum.Load())
}
