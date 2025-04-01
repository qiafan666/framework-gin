package internal

import (
	"github.com/qiafan666/gotato/commons/gcache"
)

type IClientManager interface {
	Set(v *Client)
	GetOldClients(v *Client) ([]*Client, bool, bool)
	DeleteClients(clients []*Client) bool
	GetAllStatelessConnections() *gcache.ShardLockMap[string, []*Client]
	GetAllStatefulConnections() *gcache.ShardLockMap[string, []*Client]
}

func newClientManager(numShards int) IClientManager {
	return &clientManager{
		statelessConnections: gcache.NewShardLockMap[[]*Client](numShards),
		statefulConnections:  gcache.NewShardLockMap[[]*Client](numShards),
	}
}

type clientManager struct {
	statelessConnections *gcache.ShardLockMap[string, []*Client] // 无状态连接 key: uuid, value: []*Client
	statefulConnections  *gcache.ShardLockMap[string, []*Client] // 有状态连接 key: userID, value: []*Client
}

// Set 设置连接信息
func (cm *clientManager) Set(v *Client) {
	if v.GetClientState() {
		if oldClients, ok := cm.statefulConnections.SetIfAbsent(v.parseToken.UserId, []*Client{v}); ok {
			cm.statelessConnections.Set(v.parseToken.Uuid, append(oldClients, v))
		}
	} else {
		if oldClients, ok := cm.statelessConnections.Get(v.parseToken.Uuid); ok {
			cm.statelessConnections.Set(v.parseToken.Uuid, append(oldClients, v))
		}
	}
}

// GetOldClients 获取指定平台的老连接
func (cm *clientManager) GetOldClients(v *Client) ([]*Client, bool, bool) {
	if v.GetClientState() {
		if statefulClients, ok := cm.statefulConnections.Get(v.parseToken.UserId); ok {
			var samePlatformClients []*Client
			for _, client := range statefulClients {
				if client.UserCtx.PlatformID == v.UserCtx.PlatformID {
					samePlatformClients = append(samePlatformClients, client)
				}
			}
			return samePlatformClients, true, true
		}
		return nil, false, true
	}

	if statelessClients, ok := cm.statelessConnections.Get(v.parseToken.Uuid); ok {
		return statelessClients, true, false
	}
	return nil, false, false
}

// 删除客户端的通用函数
func (cm *clientManager) deleteClientFromConnections(connections *gcache.ShardLockMap[string, []*Client], key string, connID string) {
	if shardClients, ok := connections.Get(key); ok {
		var updatedClients []*Client
		for _, client := range shardClients {
			if client.UserCtx.ConnID != connID {
				updatedClients = append(updatedClients, client)
			}
		}

		if len(updatedClients) == 0 {
			connections.Remove(key)
		} else {
			connections.Set(key, updatedClients)
		}
	}
}

// DeleteClients 优化后的实现
func (cm *clientManager) DeleteClients(clients []*Client) bool {
	if len(clients) == 0 {
		return false
	}

	// 按类型将客户端分类
	for _, client := range clients {
		if client.GetClientState() {
			cm.deleteClientFromConnections(cm.statefulConnections, client.parseToken.UserId, client.UserCtx.ConnID)
		} else {
			cm.deleteClientFromConnections(cm.statelessConnections, client.parseToken.Uuid, client.UserCtx.ConnID)
		}
	}
	return true
}

func (cm *clientManager) GetAllStatelessConnections() *gcache.ShardLockMap[string, []*Client] {
	return cm.statelessConnections
}

func (cm *clientManager) GetAllStatefulConnections() *gcache.ShardLockMap[string, []*Client] {
	return cm.statefulConnections
}
