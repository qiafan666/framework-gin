package localcache

import (
	"context"
	"errors"
	"framework-gin/common/function"
	"framework-gin/ws/constant"
	"framework-gin/ws/localcache/localcache"
	"framework-gin/ws/localcache/localcache/lru"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"time"
)

// NewOnlineCache 创建一个新的 OnlineCache 实例，并启动一个 Goroutine 监听 Redis 消息。
func NewOnlineCache(rdb redis.UniversalClient, fn func(ctx context.Context, userID string, platformIDs []int32)) *OnlineCache {
	x := &OnlineCache{
		// 初始化本地 LRU 缓存，使用 LRU 算法存储用户平台 ID 列表。
		local: lru.NewSlotLRU(1024, localcache.LRUStringHash, func() lru.LRU[string, []int32] {
			return lru.NewLayLRU[string, []int32](2048, constant.OnlineExpire/2, time.Second*3, localcache.EmptyTarget{}, func(key string, value []int32) {})
		}),
		rdb: rdb,
	}

	// 启动一个 Goroutine 订阅 Redis 中的在线状态更新频道。
	go func() {
		// 订阅 Redis 消息，处理来自指定频道的消息。
		for message := range rdb.Subscribe(function.WsCtx, constant.OnlineChannel).Channel() {
			// 解析 Redis 消息的负载内容，获取 userID 和 platformIDs。
			userID, platformIDs, err := ParseUserOnlineStatus(message.Payload)
			if err != nil {
				// 如果解析失败，记录错误并继续下一条消息。
				glog.Slog.ErrorKVs(function.WsCtx, "OnlineCache setUserOnline redis subscribe parseUserOnlineStatus", err, "payload", message.Payload, "channel", message.Channel)
				continue
			}

			// 更新本地缓存中的用户在线状态。
			storageCache := x.setUserOnline(userID, platformIDs)
			glog.Slog.DebugKVs(function.WsCtx, "OnlineCache setUserOnline", "userID", userID, "platformIDs", platformIDs, "payload", message.Payload, "storageCache", storageCache)

			// 如果回调函数不为 nil，则调用回调函数进行进一步处理。
			if fn != nil {
				fn(function.WsCtx, userID, platformIDs)
			}
		}
	}()

	return x
}

// OnlineCache 是一个用于存储和管理用户在线状态的缓存。
type OnlineCache struct {
	local lru.LRU[string, []int32] // 本地 LRU 缓存，用于存储用户在线的 platformIDs。
	rdb   redis.UniversalClient    // Redis 客户端，用于与 Redis 通信。
}

// getUserOnlinePlatform 从缓存或 Redis 获取指定用户的在线平台 ID 列表。
func (o *OnlineCache) getUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.local.Get(userID, func() ([]int32, error) {
		// 如果缓存中没有，则从 Redis 获取在线状态。
		return o.GetUserOnlinePlatformWithRedis(ctx, userID)
	})
	if err != nil {
		// 如果获取失败，记录错误并返回。
		glog.Slog.ErrorKVs(ctx, "OnlineCache GetUserOnlinePlatform", err, "userID", userID)
		return nil, err
	}
	glog.Slog.DebugKVs(ctx, "OnlineCache GetUserOnlinePlatform", "userID", userID, "platformIDs", platformIDs)
	return platformIDs, nil
}

// GetUserOnlinePlatformWithRedis 从 Redis 获取指定用户的在线平台 ID 列表。
func (o *OnlineCache) GetUserOnlinePlatformWithRedis(ctx context.Context, userID string) ([]int32, error) {
	// 从 Redis 中获取指定用户的在线平台 ID，范围为当前时间到未来。
	members, err := o.rdb.ZRangeByScore(ctx, constant.GetOnlineKey(userID), &redis.ZRangeBy{
		Min: strconv.FormatInt(time.Now().Unix(), 10),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, gerr.Wrap(err)
	}

	// 将获取到的平台 ID 转换为 int32 数组。
	platformIDs := make([]int32, 0, len(members))
	for _, member := range members {
		val, err := strconv.Atoi(member)
		if err != nil {
			return nil, gerr.Wrap(err)
		}
		platformIDs = append(platformIDs, int32(val))
	}
	return platformIDs, nil
}

// GetUserOnlinePlatform 获取指定用户的在线平台 ID 列表，并返回其副本。
func (o *OnlineCache) GetUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.getUserOnlinePlatform(ctx, userID)
	if err != nil {
		return nil, err
	}
	// 创建平台 ID 的副本，以避免外部修改。
	tmp := make([]int32, len(platformIDs))
	copy(tmp, platformIDs)
	return tmp, nil
}

// GetUserOnline 检查指定用户是否在线。
func (o *OnlineCache) GetUserOnline(ctx context.Context, userID string) (bool, error) {
	platformIDs, err := o.getUserOnlinePlatform(ctx, userID)
	if err != nil {
		return false, err
	}
	// 如果平台 ID 列表不为空，则表示用户在线。
	return len(platformIDs) > 0, nil
}

// setUserOnline 更新本地缓存中指定用户的在线平台 ID 列表。
func (o *OnlineCache) setUserOnline(userID string, platformIDs []int32) bool {
	return o.local.SetHas(userID, platformIDs)
}

// ParseUserOnlineStatus 解析 Redis 消息中的用户在线状态信息。
func ParseUserOnlineStatus(payload string) (string, []int32, error) {
	// 分割 Redis 消息负载，格式为 "platformID1:platformID2:...:userID"。
	arr := strings.Split(payload, ":")
	if len(arr) == 0 {
		return "", nil, errors.New("invalid data")
	}

	// 获取用户 ID，为分割数组中的最后一个元素。
	userID := arr[len(arr)-1]
	if userID == "" {
		return "", nil, errors.New("userID is empty")
	}

	// 将前面的平台 ID 转换为 int32 类型。
	platformIDs := make([]int32, len(arr)-1)
	for i := range platformIDs {
		platformID, err := strconv.Atoi(arr[i])
		if err != nil {
			return "", nil, err
		}
		platformIDs[i] = int32(platformID)
	}
	return userID, platformIDs, nil
}
