package constant

import "time"

const (
	OnlineKey     = "ONLINE:"
	OnlineChannel = "online_change"
	MsgChannel    = "msg_channel"
	OnlineExpire  = time.Hour / 2
)

func GetOnlineKey(userID string) string {
	return OnlineKey + userID
}
