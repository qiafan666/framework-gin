package internal

import (
	"time"
)

const (
	// Websocket Protocol.
	WSGetNewestSeq        = 1001
	WSPullMsgBySeqList    = 1002
	WSSendMsg             = 1003
	WSSendSignalMsg       = 1004
	WSPushMsg             = 2001
	WSKickOnlineMsg       = 2002
	WsLogoutMsg           = 2003
	WsSetBackgroundStatus = 2004
	WsSubUserOnlineStatus = 2005
	WSDataError           = 3001
)

const (
	// 写入消息到对端的允许时间。
	writeWait = 10 * time.Second

	// 读取下一个 pong 消息的允许时间。
	pongWait = 10 * time.Second

	// 以这个频率向对端发送 ping 消息。必须小于 pongWait。
	pingPeriod = (pongWait * 9) / 10

	// 从对端允许的最大消息大小。
	maxMessageSize = 51200
)
