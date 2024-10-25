package internal

import (
	"context"
	"framework-gin/common"
	"framework-gin/common/errs"
	"framework-gin/common/function"
	"framework-gin/middleware"
	"framework-gin/ws/constant"
	"framework-gin/ws/localcache"
	"framework-gin/ws/proto/pb"
	"framework-gin/ws/redis"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gcompress"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin"
	"github.com/qiafan666/gotato/commons/glog"
	"github.com/qiafan666/gotato/commons/gredis"
	v2 "github.com/qiafan666/gotato/v2"
	"sync"
	"sync/atomic"
	"time"
)

var wsConfig struct {
	Ws struct {
		MultiLoginPolicy int `yaml:"multi_login_policy"`
	} `yaml:"ws"`
}

func init() {
	v2.GetGotatoInstance().LoadCustomizeConfig(&wsConfig)
}

type LongConnServer interface {
	Run(r *gin.Engine)
	wsHandler(c *gin.Context)
	GetUserAllCons(userID string) ([]*Client, bool)
	GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool)
	Validate(s any) error
	UnRegister(c *Client)
	SetKickHandlerInfo(i *kickHandler)
	SubUserOnlineStatus(client *Client, data *pb.ReqSubUserOnlineStatus) (*pb.RspSubUserOnlineStatus, int)
	PubMsgChannel(ctx context.Context, data *pb.ReqPushMsgToOther) error
	gcompress.Compressor
	gcommon.Encoder
}

type WsServer struct {
	port              int
	wsMaxConnNum      int64
	registerChan      chan *Client
	unregisterChan    chan *Client
	kickHandlerChan   chan *kickHandler
	clients           UserMap
	localOnlineCache  *localcache.OnlineCache
	rdbOnline         redis.OnlineCache
	subscription      *Subscription
	clientPool        sync.Pool
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	handshakeTimeout  time.Duration
	writeBufferSize   int
	gcompress.Compressor
	gcommon.Encoder
	msgChannel redis.ChannelInterface
}

type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

func NewWsServer() *WsServer {

	rdb, err := gredis.NewRedisClient(function.WsCtx, &gredis.Config{
		ClusterMode: false,
		Address:     []string{"10.80.50.90:6379"},
		Username:    "",
		Password:    "",
		DB:          10,
	})
	if err != nil {
		glog.Slog.ErrorKVs(function.WsCtx, "redis connect error", "err", err.Error())
		panic(err)
	}

	wsServer := &WsServer{
		wsMaxConnNum:     100000,
		writeBufferSize:  4096,
		handshakeTimeout: 10 * time.Second,
		clientPool: sync.Pool{
			New: func() any {
				return new(Client)
			},
		},
		registerChan:    make(chan *Client, 1000),
		unregisterChan:  make(chan *Client, 1000),
		kickHandlerChan: make(chan *kickHandler, 1000),
		clients:         newUserMap(),
		subscription:    newSubscription(),
		Compressor:      gcompress.NewGzipCompressor(),
		Encoder:         gcommon.NewGobEncoder(),
	}
	wsServer.localOnlineCache = localcache.NewOnlineCache(rdb, wsServer.SubscriberUserOnlineStatusChanges)
	wsServer.rdbOnline = redis.NewUserOnline(rdb)
	wsServer.msgChannel = redis.NewMsgChannel(rdb)
	return wsServer
}

func (ws *WsServer) Run(r *gin.Engine) {
	var (
		client *Client
	)

	go func() {
		for {
			select {
			case client = <-ws.registerChan:
				ws.registerClient(client)
			case client = <-ws.unregisterChan:
				ws.unregisterClient(client)
			case onlineInfo := <-ws.kickHandlerChan:
				ws.multiTerminalLoginChecker(onlineInfo.clientOK, onlineInfo.oldClients, onlineInfo.newClient)
			}
		}
	}()
	r.GET("/ws", ws.wsHandler)
}

func (ws *WsServer) wsHandler(c *gin.Context) {

	// 创建一个新的连接上下文
	connCtx := newContext(c.Writer, c.Request)

	// 检查当前在线用户连接数是否超过最大限制
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
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
	parseToken, err := ws.ParseToken(connCtx, connCtx.GetToken())
	if err != nil {
		if !shouldSendResp {
			// 创建一个WebSocket连接对象并尝试通过WebSocket发送错误信息
			wsLongConn := newGWebSocket(ws.handshakeTimeout, ws.writeBufferSize)
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
	wsLongConn := newGWebSocket(ws.handshakeTimeout, ws.writeBufferSize)
	if err = wsLongConn.GenerateLongConn(c.Writer, c.Request); err != nil {
		// 如果长连接创建失败，握手过程中已处理错误
		glog.Slog.WarnF(connCtx.Ctx, "长连接创建失败: %v", err)
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
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(connCtx, wsLongConn, ws)
	client.parseToken = parseToken

	// 将客户端注册到服务器并开始消息处理
	ws.registerChan <- client
	go client.readMessage()
}

func (ws *WsServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

func (ws *WsServer) GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool) {
	return ws.clients.Get(userID, platform)
}

func (ws *WsServer) Validate(s any) error {
	return gcommon.Validate(s)
}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

// SetKickHandlerInfo 暴露出去的踢人检查，分布式环境下使用
func (ws *WsServer) SetKickHandlerInfo(i *kickHandler) {
	ws.kickHandlerChan <- i
}

// ------------------------ inner function ------------------------

// multiTerminalLoginChecker 多端登录检查
func (ws *WsServer) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	switch wsConfig.Ws.MultiLoginPolicy {
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
		ws.clients.DeleteClients(newClient.parseToken.UserId, oldClients)
		for _, c := range oldClients {
			err := c.KickOnlineMessage(pb.KickReason_OnlyOneClient)
			if err != nil {
				glog.Slog.WarnKVs(c.UserCtx.Ctx, "multiTerminalLoginChecker,KickOnlineMessage", "err", err)
			}
		}
	}
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK     bool
		clientOK   bool
		oldClients []*Client
	)
	oldClients, userOK, clientOK = ws.clients.Get(client.parseToken.UserId, client.UserCtx.PlatformID)
	if !userOK {
		ws.clients.Set(client.parseToken.UserId, client)
		glog.Slog.DebugKVs(client.UserCtx.Ctx, "registerClient,user not exist", "userID", client.parseToken.UserId, "platformID", client.UserCtx.PlatformID)
		ws.onlineUserNum.Add(1)
		ws.onlineUserConnNum.Add(1)
	} else {
		ws.multiTerminalLoginChecker(clientOK, oldClients, client)
		glog.Slog.DebugKVs(client.UserCtx.Ctx, "registerClient,user exist", "userID", client.parseToken.UserId, "platformID", client.UserCtx.PlatformID)
		if clientOK {
			ws.clients.Set(client.parseToken.UserId, client)
			// 当前平台连接已经存在，增加连接数
			glog.Slog.InfoKVs(client.UserCtx.Ctx, "registerClient,repeat login", "userID", client.parseToken.UserId, "platformID",
				client.UserCtx.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			ws.onlineUserConnNum.Add(1)
		} else {
			ws.clients.Set(client.parseToken.UserId, client)
			ws.onlineUserConnNum.Add(1)
		}
	}

	//TODO 多节点通知在线状态，如何处理

	glog.Slog.InfoKVs(
		client.UserCtx.Ctx,
		"registerClient,user online",
		"online user Num",
		ws.onlineUserNum.Load(),
		"online user conn Num",
		ws.onlineUserConnNum.Load(),
	)
}

func getRemoteAdders(client []*Client) string {
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

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.DeleteClients(client.parseToken.UserId, []*Client{client})
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
	}
	ws.onlineUserConnNum.Add(-1)
	ws.subscription.DelClient(client)
	glog.Slog.WarnKVs(client.UserCtx.Ctx, "unregisterClient user offline", "close reason", client.closedErr, "online user Num",
		ws.onlineUserNum.Load(), "online user conn Num", ws.onlineUserConnNum.Load())
}

func (ws *WsServer) ParseToken(connContext *UserConnContext, token string) (*pb.ParseToken, error) {
	if token == "test" {
		return &pb.ParseToken{UserId: "ning"}, nil
	}
	resultMap, err := middleware.ParseToken(connContext.Req.Header.Get(common.HeaderAuthorization))
	if err != nil {
		return &pb.ParseToken{}, err
	}
	return &pb.ParseToken{
		UserId: gcast.ToString(resultMap["user_id"]),
	}, nil
}

func (ws *WsServer) subscribe() {
	subscribe := ws.msgChannel.MsgSubscribe(function.WsCtx, constant.MsgChannel)

	// 启动一个 goroutine 来接收消息
	gcommon.Go(func() {
		for msg := range subscribe {
			if msg == nil {
				continue
			}
			if msg.All {
				for _, client := range ws.clients.AllCons() {
					if client.closed.Load() {
						continue
					}
					client.w.Lock()
					glog.Slog.DebugKVs(client.UserCtx.Ctx, "subscribe message", "msg", msg.String())
					err := client.PushMessage(msg)
					if err != nil {
						glog.Slog.ErrorKVs(client.UserCtx.Ctx, "subscribe message,PushMessage error", "err", err)
						continue
					}
					client.w.Unlock()
				}
			} else {
				oldClients, userOK, clientOK := ws.GetUserPlatformCons(msg.ToUserID, int(msg.PlatformID))
				if userOK && clientOK {
					for _, oldClient := range oldClients {
						if oldClient.closed.Load() {
							continue
						}
						oldClient.w.Lock()
						for _, client := range oldClients {
							glog.Slog.DebugKVs(client.UserCtx.Ctx, "subscribe message", "msg", msg.String())
							err := client.PushMessage(msg)
							if err != nil {
								glog.Slog.ErrorKVs(client.UserCtx.Ctx, "subscribe message,PushMessage error", "err", err)
								continue
							}
						}
						oldClient.w.Unlock()
					}
				}
			}
		}
	})
}

func (ws *WsServer) PubMsgChannel(ctx context.Context, pubMsg *pb.ReqPushMsgToOther) error {
	return ws.msgChannel.MsgPublish(ctx, pubMsg)
}
