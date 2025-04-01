package internal

import (
	v2 "github.com/qiafan666/gotato/v2"
	"time"
)

var config struct {
	Ws struct {
		MulitLoginPolicy int `yaml:"mulit_login_policy"`
		Protocol         int `yaml:"protocol"`
		PrivateLiveTime  int `yaml:"private_live_time"`
		PublicLiveTime   int `yaml:"public_live_time"`
	} `yaml:"ws"`
}

func init() {
	v2.GetGotato().LoadCustomCfg(&config)
}

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

const (
	// MessageText 表示 UTF-8 编码的文本消息，例如 JSON。
	MessageText = iota + 1
	// MessageBinary 表示二进制消息，例如 protobufs。
	MessageBinary
	// CloseMessage 表示关闭控制消息。可选消息负载包含一个数字代码和文本。
	// 使用 FormatCloseMessage 函数格式化关闭消息负载。
	CloseMessage = 8

	// PingMessage 表示 ping 控制消息。可选消息负载为 UTF-8 编码文本。
	PingMessage = 9

	// PongMessage 表示 pong 控制消息。可选消息负载为 UTF-8 编码文本。
	PongMessage = 10
)
