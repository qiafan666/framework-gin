package internal

import (
	"github.com/qiafan666/gotato/commons/gcommon"
	"sync"
	"time"
)

type UserMap interface {
	AllCons() []*Client
	GetAll(userID string) ([]*Client, bool)
	Get(userID string, platformID int) ([]*Client, bool, bool)
	Set(userID string, v *Client)
	DeleteClients(userID string, clients []*Client) (isDeleteUser bool)
	UserState() <-chan UserState
	GetAllUserStatus(deadline time.Time, nowtime time.Time) []UserState
	RecvSubChange(userID string, platformIDs []int32) bool
}

type UserState struct {
	UserID  string
	Online  []int32
	Offline []int32
}

type UserPlatform struct {
	Time    time.Time
	Clients []*Client
}

func (u *UserPlatform) PlatformIDs() []int32 {
	if len(u.Clients) == 0 {
		return nil
	}
	platformIDs := make([]int32, 0, len(u.Clients))
	for _, client := range u.Clients {
		platformIDs = append(platformIDs, int32(client.UserCtx.PlatformID))
	}
	return platformIDs
}

func (u *UserPlatform) PlatformIDSet() map[int32]struct{} {
	if len(u.Clients) == 0 {
		return nil
	}
	platformIDs := make(map[int32]struct{})
	for _, client := range u.Clients {
		platformIDs[int32(client.UserCtx.PlatformID)] = struct{}{}
	}
	return platformIDs
}

func newUserMap() UserMap {
	return &userMap{
		data: make(map[string]*UserPlatform),
		ch:   make(chan UserState, 10000),
	}
}

type userMap struct {
	lock sync.RWMutex
	data map[string]*UserPlatform
	ch   chan UserState
}

func (u *userMap) RecvSubChange(userID string, platformIDs []int32) bool {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result, ok := u.data[userID]
	if !ok {
		return false
	}
	localPlatformIDs := result.PlatformIDSet()
	for _, platformID := range platformIDs {
		delete(localPlatformIDs, platformID)
	}
	if len(localPlatformIDs) == 0 {
		return false
	}
	u.push(userID, result, nil)
	return true
}

func (u *userMap) push(userID string, userPlatform *UserPlatform, offline []int32) bool {
	select {
	case u.ch <- UserState{UserID: userID, Online: userPlatform.PlatformIDs(), Offline: offline}:
		userPlatform.Time = time.Now()
		return true
	default:
		return false
	}
}

func (u *userMap) AllCons() []*Client {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result := make([]*Client, 0, 1000)
	for _, userPlatform := range u.data {
		for _, client := range userPlatform.Clients {
			result = append(result, client)
		}
	}
	return result
}

func (u *userMap) GetAll(userID string) ([]*Client, bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result, ok := u.data[userID]
	if !ok {
		return nil, false
	}
	return result.Clients, true
}

func (u *userMap) Get(userID string, platformID int) ([]*Client, bool, bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result, ok := u.data[userID]
	if !ok {
		return nil, false, false
	}
	var clients []*Client
	for _, client := range result.Clients {
		if client.UserCtx.PlatformID == platformID {
			clients = append(clients, client)
		}
	}
	return clients, true, len(clients) > 0
}

func (u *userMap) Set(userID string, client *Client) {
	u.lock.Lock()
	defer u.lock.Unlock()
	result, ok := u.data[userID]
	if ok {
		result.Clients = append(result.Clients, client)
	} else {
		result = &UserPlatform{
			Clients: []*Client{client},
		}
		u.data[userID] = result
	}
	u.push(client.parseToken.UserId, result, nil)
}

func (u *userMap) DeleteClients(userID string, clients []*Client) (isDeleteUser bool) {
	if len(clients) == 0 {
		return false
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	result, ok := u.data[userID]
	if !ok {
		return false
	}
	offline := make([]int32, 0, len(clients))
	deleteAddr := gcommon.SliceSetAny(clients, func(client *Client) string {
		return client.UserCtx.GetRemoteAddr()
	})
	tmp := result.Clients
	result.Clients = result.Clients[:0]
	for _, client := range tmp {
		if _, delCli := deleteAddr[client.UserCtx.GetRemoteAddr()]; delCli {
			offline = append(offline, int32(client.UserCtx.PlatformID))
		} else {
			result.Clients = append(result.Clients, client)
		}
	}
	defer u.push(userID, result, offline)
	if len(result.Clients) > 0 {
		return false
	}
	delete(u.data, userID)
	return true
}

// GetAllUserStatus 返回死亡时间在用户状态时间之后的用户状态
func (u *userMap) GetAllUserStatus(deadline time.Time, nowtime time.Time) (result []UserState) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	result = make([]UserState, 0, len(u.data))
	for userID, userPlatform := range u.data {
		if deadline.Before(userPlatform.Time) {
			continue
		}
		userPlatform.Time = nowtime
		online := make([]int32, 0, len(userPlatform.Clients))
		for _, client := range userPlatform.Clients {
			online = append(online, int32(client.UserCtx.PlatformID))
		}
		result = append(result, UserState{UserID: userID, Online: online})
	}
	return result
}

func (u *userMap) UserState() <-chan UserState {
	return u.ch
}
