package rpccache

import (
	"context"
	"github.com/qiafan666/gotato/commons/gcache"
)

func NewLocalOnlineCache() *OnlineCache {
	return &OnlineCache{
		localCache: gcache.NewCache(gcache.DefaultExpiration, 0),
	}
}

type OnlineCache struct {
	localCache *gcache.Cache
}

func (o *OnlineCache) getUserOnlinePlatform(ctx context.Context, userID string) ([]int32, bool) {
	get, b := o.localCache.Get(userID)
	if b {
		return get.([]int32), true
	}
	return nil, false
}

func (o *OnlineCache) GetUserOnline(ctx context.Context, userID string) bool {
	platformIDs, flag := o.getUserOnlinePlatform(ctx, userID)
	if !flag {
		return false
	}
	return len(platformIDs) > 0
}

func (o *OnlineCache) setUserOnline(userID string, platformIDs []int32) {
	o.localCache.SetDefault(userID, platformIDs)
}
