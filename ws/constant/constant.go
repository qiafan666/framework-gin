package constant

// TokenStatus
const (
	NormalToken  = 0
	InValidToken = 1
	KickedToken  = 2
	ExpiredToken = 3
)

// 登录策略
const (
	// DefaultNotKick 默认不踢人
	DefaultNotKick = 0
	// AllLoginButSameTermKick 所有端登录，但同端踢人
	AllLoginButSameTermKick = 1
	// SingleTerminalLogin 只允许单端登录
	SingleTerminalLogin = 2
	// WebAndOther 网页端可以同时在线，其他端只能在一个端登录
	WebAndOther = 3
	// PcMobileAndWeb PC端互斥，移动端互斥，但网页端可以同时在线
	PcMobileAndWeb = 4
	// PCAndOther PC端可以同时在线，但其他端只允许一个端登录
	PCAndOther = 5
)

const (
	OnlineStatus  = "online"
	OfflineStatus = "offline"
	Registered    = "registered"
	UnRegistered  = "unregistered"

	Online  = 1
	Offline = 0
)

// UserContext 相关常量
const (
	RequestID       = "requestID"
	OpUserID        = "opUserID"
	ConnID          = "connID"
	OpUserPlatform  = "platform"
	Token           = "token"
	RpcCustomHeader = "customHeader" // rpc中间件自定义ctx参数
	CheckKey        = "CheckKey"
	TriggerID       = "triggerID"
	RemoteAddr      = "remoteAddr"
)
