package internal

import (
	"context"
	"fmt"
	"framework-gin/ws/cache"
	"framework-gin/ws/constant"
	"framework-gin/ws/errs"
	"framework-gin/ws/mcontext"
	"framework-gin/ws/proto/pb"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gcompress"
	"github.com/qiafan666/gotato/commons/glog"
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
	//SetDiscoveryRegistry(client discovery.SvcDiscoveryRegistry, config *Config)
	KickUserConn(client *Client) error
	UnRegister(c *Client)
	SetKickHandlerInfo(i *kickHandler)
	//SubUserOnlineStatus(ctx context.Context, client *Client, data *Req) ([]byte, error)
	gcompress.Compressor
	gcommon.Encoder
	//MessageHandler
}

type WsServer struct {
	port              int
	wsMaxConnNum      int64
	registerChan      chan *Client
	unregisterChan    chan *Client
	kickHandlerChan   chan *kickHandler
	clients           UserMap
	online            *rpccache.OnlineCache
	subscription      *Subscription
	clientPool        sync.Pool
	onlineUserNum     atomic.Int64
	onlineUserConnNum atomic.Int64
	handshakeTimeout  time.Duration
	writeBufferSize   int
	//userClient        *rpcclient.UserRpcClient
	//authClient        *rpcclient.Auth
	//disCov            discovery.SvcDiscoveryRegistry
	gcompress.Compressor
	gcommon.Encoder
	//MessageHandler
}

type kickHandler struct {
	clientOK   bool
	oldClients []*Client
	newClient  *Client
}

//func (ws *WsServer) SetDiscoveryRegistry(disCov discovery.SvcDiscoveryRegistry, config *Config) {
//	ws.MessageHandler = NewGrpcHandler(ws.validate, disCov, &config.Share.RpcRegisterName)
//	u := rpcclient.NewUserRpcClient(disCov, config.Share.RpcRegisterName.User, config.Share.IMAdminUserID)
//	ws.authClient = rpcclient.NewAuth(disCov, config.Share.RpcRegisterName.Auth)
//	ws.userClient = &u
//	ws.disCov = disCov
//}

//func (ws *WsServer) SetUserOnlineStatus(ctx context.Context, client *Client, status int32) {
//	err := ws.userClient.SetUserStatus(ctx, client.UserID, status, client.PlatformID)
//	if err != nil {
//		log.ZWarn(ctx, "SetUserStatus err", err)
//	}
//	switch status {
//	case constant.Online:
//		ws.webhookAfterUserOnline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOnline, client.UserID, client.PlatformID, client.IsBackground, client.ctx.GetConnID())
//	case constant.Offline:
//		ws.webhookAfterUserOffline(ctx, &ws.msgGatewayConfig.WebhooksConfig.AfterUserOffline, client.UserID, client.PlatformID, client.ctx.GetConnID())
//	}
//}

func (ws *WsServer) UnRegister(c *Client) {
	ws.unregisterChan <- c
}

func (ws *WsServer) Validate(s any) error {
	return gcommon.Validate(s)
}

func (ws *WsServer) GetUserAllCons(userID string) ([]*Client, bool) {
	return ws.clients.GetAll(userID)
}

func (ws *WsServer) GetUserPlatformCons(userID string, platform int) ([]*Client, bool, bool) {
	return ws.clients.Get(userID, platform)
}

func NewWsServer() *WsServer {

	return &WsServer{
		wsMaxConnNum:     100000,
		writeBufferSize:  4096,
		handshakeTimeout: 10,
		clientPool: sync.Pool{
			New: func() any {
				return new(Client)
			},
		},
		registerChan:    make(chan *Client, 1000),
		unregisterChan:  make(chan *Client, 1000),
		kickHandlerChan: make(chan *kickHandler, 1000),
		//validate:        v,
		clients:      newUserMap(),
		subscription: newSubscription(),
		Compressor:   gcompress.NewGzipCompressor(),
		Encoder:      gcommon.NewGobEncoder(),
	}
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
	r.GET("/", ws.wsHandler)
}

//var concurrentRequest = 3
//
//func (ws *WsServer) sendUserOnlineInfoToOtherNode(ctx context.Context, client *Client) error {
//	conns, err := ws.disCov.GetConns(ctx, ws.msgGatewayConfig.Share.RpcRegisterName.MessageGateway)
//	if err != nil {
//		return err
//	}
//
//	wg := errgroup.Group{}
//	wg.SetLimit(concurrentRequest)
//
//	// Online push user online message to other node
//	for _, v := range conns {
//		v := v
//		glog.Slog.DebugKVs(ctx, " sendUserOnlineInfoToOtherNode conn ", "target", v.Target())
//		if v.Target() == ws.disCov.GetSelfConnTarget() {
//			glog.Slog.DebugKVs(ctx, "Filter out this node", "node", v.Target())
//			continue
//		}
//
//		wg.Go(func() error {
//			msgClient := ws1.NewMsgGatewayClient(v)
//			_, err := msgClient.MultiTerminalLoginCheck(ctx, &ws1.MultiTerminalLoginCheckReq{
//				UserID:     client.UserID,
//				PlatformID: int32(client.PlatformID), Token: client.token,
//			})
//			if err != nil {
//				glog.Slog.WarnKVs(ctx, "MultiTerminalLoginCheck err", err, "node", v.Target())
//			}
//			return nil
//		})
//	}
//
//	_ = wg.Wait()
//	return nil
//}

func (ws *WsServer) SetKickHandlerInfo(i *kickHandler) {
	ws.kickHandlerChan <- i
}

func (ws *WsServer) registerClient(client *Client) {
	var (
		userOK     bool
		clientOK   bool
		oldClients []*Client
	)
	oldClients, userOK, clientOK = ws.clients.Get(client.UserID, client.PlatformID)
	if !userOK {
		ws.clients.Set(client.UserID, client)
		glog.Slog.DebugKVs(client.ctx, "user not exist", "userID", client.UserID, "platformID", client.PlatformID)
		ws.onlineUserNum.Add(1)
		ws.onlineUserConnNum.Add(1)
	} else {
		ws.multiTerminalLoginChecker(clientOK, oldClients, client)
		glog.Slog.DebugKVs(client.ctx, "user exist", "userID", client.UserID, "platformID", client.PlatformID)
		if clientOK {
			ws.clients.Set(client.UserID, client)
			// 当前平台连接已经存在，增加连接数
			glog.Slog.InfoKVs(client.ctx, "repeat login", "userID", client.UserID, "platformID",
				client.PlatformID, "old remote addr", getRemoteAdders(oldClients))
			ws.onlineUserConnNum.Add(1)
		} else {
			ws.clients.Set(client.UserID, client)
			ws.onlineUserConnNum.Add(1)
		}
	}

	//TODO 多节点通知在线状态，如何处理
	//wg := sync.WaitGroup{}
	//
	//if ws.msgGatewayConfig.Discovery.Enable != "k8s" {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		_ = ws.sendUserOnlineInfoToOtherNode(client.ctx, client)
	//	}()
	//}

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	ws.SetUserOnlineStatus(client.ctx, client, constant.Online)
	//}()
	//wg.Wait()

	glog.Slog.InfoKVs(
		client.ctx,
		"user online",
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
			ret = c.ctx.GetRemoteAddr()
		} else {
			ret += "@" + c.ctx.GetRemoteAddr()
		}
	}
	return ret
}

func (ws *WsServer) KickUserConn(client *Client) error {
	ws.clients.DeleteClients(client.UserID, []*Client{client})
	return client.KickOnlineMessage()
}

func (ws *WsServer) multiTerminalLoginChecker(clientOK bool, oldClients []*Client, newClient *Client) {
	switch wsConfig.Ws.MultiLoginPolicy {
	case constant.DefaultNotKick:
	case constant.PCAndOther:
		if constant.PlatformIDToClass(newClient.PlatformID) == constant.TerminalPC {
			return
		}
		fallthrough
	case constant.AllLoginButSameTermKick:
		if !clientOK {
			return
		}
		ws.clients.DeleteClients(newClient.UserID, oldClients)
		for _, c := range oldClients {
			err := c.KickOnlineMessage()
			if err != nil {
				glog.Slog.WarnKVs(c.ctx, "KickOnlineMessage", "err", err)
			}
		}
		ctx := mcontext.WithMustInfoCtx(
			[]string{newClient.ctx.GetOperationID(), newClient.ctx.GetUserID(),
				constant.PlatformIDToName(newClient.PlatformID), newClient.ctx.GetConnID()},
		)
		//TODO 设置Token过期
		if _, err := ws.InvalidateToken(ctx, newClient.token, newClient.UserID, newClient.PlatformID); err != nil {
			glog.Slog.WarnKVs(ctx, "InvalidateToken err", "err", err, "userID", newClient.UserID, "platformID", newClient.PlatformID)
		}
	}
}

func (ws *WsServer) unregisterClient(client *Client) {
	defer ws.clientPool.Put(client)
	isDeleteUser := ws.clients.DeleteClients(client.UserID, []*Client{client})
	if isDeleteUser {
		ws.onlineUserNum.Add(-1)
	}
	ws.onlineUserConnNum.Add(-1)
	ws.subscription.DelClient(client)
	//ws.SetUserOnlineStatus(client.ctx, client, constant.Offline)
	glog.Slog.InfoKVs(client.ctx, "user offline", "close reason", client.closedErr, "online user Num",
		ws.onlineUserNum.Load(), "online user conn Num", ws.onlineUserConnNum.Load())
}

// validateRespWithRequest 检查token中的用户ID和平台ID是否与请求中的匹配
func (ws *WsServer) validateRespWithRequest(ctx *UserConnContext, resp *pb.ParseToken) error {
	userID := ctx.GetUserID()
	platformID := gcast.ToInt32(ctx.GetPlatformID())
	if resp.UserID != userID {
		return errs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token uid %s != userID %s", resp.UserID, userID))
	}

	if resp.PlatformID != platformID {
		return errs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token platform %d != platformID %d", resp.PlatformID, platformID))
	}
	return nil
}

func (ws *WsServer) wsHandler(c *gin.Context) {

	// 创建一个新的连接上下文
	connContext := newContext(c.Writer, c.Request)

	// 检查当前在线用户连接数是否超过最大限制
	if ws.onlineUserConnNum.Load() >= ws.wsMaxConnNum {
		// 如果超过最大连接数限制，通过HTTP返回错误并停止处理
		httpError(connContext, errs.ErrConnOverMaxNumLimit.WrapMsg("超过最大连接数限制"))
		return
	}

	// 解析必要的参数（例如用户ID、Token）
	err := connContext.ParseEssentialArgs()
	if err != nil {
		// 如果解析过程中出错，通过HTTP返回错误并停止处理
		httpError(connContext, err)
		return
	}

	// 调用认证客户端，解析从上下文中获取的Token
	resp, err := ws.ParseToken(connContext, connContext.GetToken())
	if err != nil {
		// 如果Token解析失败，根据上下文标志决定是否通过WebSocket发送错误信息
		shouldSendError := connContext.ShouldSendResp()
		if shouldSendError {
			// 创建一个WebSocket连接对象并尝试通过WebSocket发送错误信息
			wsLongConn := newGWebSocket(ws.handshakeTimeout, ws.writeBufferSize)
			if err = wsLongConn.RespondWithError(err, c.Writer, c.Request); err == nil {
				// 如果错误信息成功通过WebSocket发送，则停止处理
				return
			}
		}
		// 如果不需要或无法通过WebSocket发送，返回HTTP错误并停止处理
		httpError(connContext, err)
		return
	}

	// 验证认证响应是否与请求匹配（例如用户ID和平台ID）
	err = ws.validateRespWithRequest(connContext, &resp)
	if err != nil {
		// 如果验证失败，通过HTTP返回错误并停止处理
		httpError(connContext, err)
		return
	}

	// 创建WebSocket长连接对象
	wsLongConn := newGWebSocket(ws.handshakeTimeout, ws.writeBufferSize)
	if err = wsLongConn.GenerateLongConn(c.Writer, c.Request); err != nil {
		// 如果长连接创建失败，握手过程中已处理错误
		glog.Slog.WarnKVs(connContext, "长连接创建失败: %v", err)
		return
	} else {
		// 检查是否需要通过WebSocket发送正常响应
		shouldSendSuccessResp := connContext.ShouldSendResp()
		if shouldSendSuccessResp {
			// 尝试通过WebSocket发送成功信息
			if err = wsLongConn.RespondWithSuccess(); err != nil {
				// 如果成功信息发送成功，则结束处理
				return
			}
		}
	}

	// 从客户端池中获取一个客户端对象，重置其状态，并将其与当前的WebSocket长连接关联
	client := ws.clientPool.Get().(*Client)
	client.ResetClient(connContext, wsLongConn, ws)

	// 将客户端注册到服务器并开始消息处理
	ws.registerChan <- client
	go client.readMessage()
}

func (ws *WsServer) ParseToken(connContext *UserConnContext, token string) (pb.ParseToken, error) {
	//TODO
	return pb.ParseToken{}, nil
}

func (ws *WsServer) InvalidateToken(ctx context.Context, token string, userId string, platformId int) (pb.InvalidateToken, error) {
	//TODO 过期token处理
	return pb.InvalidateToken{}, nil
}
