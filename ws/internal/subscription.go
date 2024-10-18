package internal

import (
	"context"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"sync"
)

// SubscriberUserOnlineStatusChanges 处理用户在线状态的变更。
// 它检查用户的订阅是否发生变化，记录日志并推送用户更新后的在线状态给订阅的客户端。
func (ws *WsServer) SubscriberUserOnlineStatusChanges(ctx context.Context, userID string, platformIDs []int32) {
	// 检查用户的订阅是否发生变化，并记录结果到日志中。
	if ws.clients.RecvSubChange(userID, platformIDs) {
		glog.Slog.DebugKVs(ctx, "gateway receive subscription message and go back online", "userID", userID, "platformIDs", platformIDs)
	} else {
		glog.Slog.DebugKVs(ctx, "gateway ignore user online status changes", "userID", userID, "platformIDs", platformIDs)
	}
	// 推送用户的在线状态变更信息。
	ws.pushUserIDOnlineStatus(ctx, userID, platformIDs)
}

// pushUserIDOnlineStatus 推送指定用户 ID 的在线状态变更给所有订阅该用户的客户端。
func (ws *WsServer) pushUserIDOnlineStatus(ctx context.Context, userID string, platformIDs []int32) {
	// 获取订阅了该用户的所有客户端。
	clients := ws.subscription.GetClient(userID)
	if len(clients) == 0 {
		return
	}
	// 创建并序列化用户在线状态消息。
	onlineStatus, err := proto.Marshal(&pb.RspSubUserOnlineStatus{
		Subscribers: []*pb.SubUserOnlineStatusElem{{UserID: userID, OnlinePlatformIDs: platformIDs}},
	})
	if err != nil {
		glog.Slog.ErrorKVs(ctx, "pushUserIDOnlineStatus json.Marshal", "err", err, "userID", userID, "platformIDs", platformIDs)
		return
	}
	// 将在线状态消息推送给所有订阅该用户的客户端。
	for _, client := range clients {
		if err = client.PushUserOnlineStatus(onlineStatus); err != nil {
			glog.Slog.ErrorKVs(ctx, "UserSubscribeOnlineStatusNotification push failed", "err", err,
				"userID", client.parseToken.UserID, "platformID", client.PlatformID, "changeUserID", userID, "changePlatformID", platformIDs)
		}
	}
}

// SubUserOnlineStatus 处理客户端订阅或取消订阅用户在线状态的请求。
// 它将请求的数据解析后更新订阅状态，并返回当前订阅的用户的在线状态。
func (ws *WsServer) SubUserOnlineStatus(ctx context.Context, client *Client, data *Req) (proto.Message, int) {
	var sub pb.ReqSubUserOnlineStatus
	// 解析请求数据，如果解析失败则返回错误。
	if err := proto.Unmarshal(data.Data, &sub); err != nil {
		glog.Slog.ErrorKVs(ctx, "SubUserOnlineStatus proto.Unmarshal", "err", err, "data", data.Data)
		return nil, gerr.UnKnowError
	}
	// 更新订阅的用户和取消订阅的用户。
	ws.subscription.Sub(client, sub.SubscribeUserID, sub.UnsubscribeUserID)

	var resp *pb.RspSubUserOnlineStatus
	// 如果有订阅用户，查询这些用户的在线状态。
	if len(sub.SubscribeUserID) > 0 {
		resp.Subscribers = make([]*pb.SubUserOnlineStatusElem, 0, len(sub.SubscribeUserID))
		for _, userID := range sub.SubscribeUserID {
			platformIDs, err := ws.localOnlineCache.GetUserOnlinePlatform(ctx, userID)
			if err != nil {
				glog.Slog.ErrorKVs(ctx, "SubUserOnlineStatus GetUserOnlinePlatform failed", "err", err, "userID", userID)
				return nil, gerr.UnKnowError
			}
			// 添加用户的在线状态信息到响应中。
			resp.Subscribers = append(resp.Subscribers, &pb.SubUserOnlineStatusElem{
				UserID:            userID,
				OnlinePlatformIDs: platformIDs,
			})
		}
	}
	// 将响应数据序列化后返回。
	return resp, gerr.OK
}

// newSubscription 创建并返回一个新的订阅管理对象。
func newSubscription() *Subscription {
	return &Subscription{
		userIDs: make(map[string]*subClient),
	}
}

// subClient 表示一个订阅了某个用户的客户端集合。
type subClient struct {
	clients map[string]*Client // 客户端的地址和客户端对象的映射。
}

// Subscription 负责管理用户在线状态的订阅。
// 它存储了用户 ID 与订阅该用户的客户端之间的映射关系。
type Subscription struct {
	lock    sync.RWMutex          // 用于保护对 userIDs 的并发访问。
	userIDs map[string]*subClient // 用户 ID 与订阅该用户的客户端集合的映射。
}

// DelClient 从订阅中移除客户端，并清除该客户端对所有用户的订阅。
func (s *Subscription) DelClient(client *Client) {
	// 获取并锁定客户端的订阅用户列表。
	client.subLock.Lock()
	userIDs := gcommon.MapKeys(client.subUserIDs)
	for _, userID := range userIDs {
		delete(client.subUserIDs, userID)
	}
	client.subLock.Unlock()
	if len(userIDs) == 0 {
		return
	}
	addr := client.userCtx.GetRemoteAddr()
	// 锁定 Subscription 以安全地移除订阅。
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, userID := range userIDs {
		sub, ok := s.userIDs[userID]
		if !ok {
			continue
		}
		delete(sub.clients, addr)
		if len(sub.clients) == 0 {
			delete(s.userIDs, userID)
		}
	}
}

// GetClient 获取所有订阅指定用户 ID 的客户端。
func (s *Subscription) GetClient(userID string) []*Client {
	s.lock.RLock()
	defer s.lock.RUnlock()
	cs, ok := s.userIDs[userID]
	if !ok {
		return nil
	}
	clients := make([]*Client, 0, len(cs.clients))
	for _, client := range cs.clients {
		clients = append(clients, client)
	}
	return clients
}

// Sub 管理客户端对用户在线状态的订阅和取消订阅。
// 它更新客户端的订阅列表，并根据需要更新 Subscription 的映射。
func (s *Subscription) Sub(client *Client, addUserIDs, delUserIDs []string) {
	if len(addUserIDs)+len(delUserIDs) == 0 {
		return
	}
	var (
		del = make(map[string]struct{})
		add = make(map[string]struct{})
	)
	// 锁定客户端的订阅列表进行修改。
	client.subLock.Lock()
	for _, userID := range delUserIDs {
		if _, ok := client.subUserIDs[userID]; !ok {
			continue
		}
		del[userID] = struct{}{}
		delete(client.subUserIDs, userID)
	}
	for _, userID := range addUserIDs {
		delete(del, userID)
		if _, ok := client.subUserIDs[userID]; ok {
			continue
		}
		client.subUserIDs[userID] = struct{}{}
		add[userID] = struct{}{}
	}
	client.subLock.Unlock()
	if len(del)+len(add) == 0 {
		return
	}
	addr := client.userCtx.GetRemoteAddr()
	// 锁定 Subscription 更新订阅映射。
	s.lock.Lock()
	defer s.lock.Unlock()
	for userID := range del {
		sub, ok := s.userIDs[userID]
		if !ok {
			continue
		}
		delete(sub.clients, addr)
		if len(sub.clients) == 0 {
			delete(s.userIDs, userID)
		}
	}
	for userID := range add {
		sub, ok := s.userIDs[userID]
		if !ok {
			sub = &subClient{clients: make(map[string]*Client)}
			s.userIDs[userID] = sub
		}
		sub.clients[addr] = client
	}
}
