package errs

const (
	DatabaseError = 50002 // 数据库错误 (Redis/MySQL 等)
	NetworkError  = 50004 // 网络错误
	DataError     = 50007 // 数据错误
	CallbackError = 50000 // 回调错误

	// 通用错误码。
	ServerInternalError = 500  // 服务器内部错误
	ArgsError           = 5001 // 输入参数错误
	NoPermissionError   = 5002 // 权限不足
	DuplicateKeyError   = 5003 // 主键重复
	RecordNotFoundError = 5004 // 记录不存在

	// 账户错误码。
	UserIDNotFoundError    = 5101 // 用户ID不存在或未注册
	RegisteredAlreadyError = 5102 // 用户已注册

	// Token 错误码。
	TokenExpiredError     = 5501 // Token 过期
	TokenInvalidError     = 5502 // Token 无效
	TokenMalformedError   = 5503 // Token 格式错误
	TokenNotValidYetError = 5504 // Token 尚未生效
	TokenUnknownError     = 5505 // Token 未知错误
	TokenKickedError      = 5506 // Token 已被踢出
	TokenNotExistError    = 5507 // Token 不存在

	// 长连接网关错误码。
	ConnOverMaxNumLimit  = 5601 // 超过最大连接数限制
	ConnArgsErr          = 5602 // 连接参数错误
	PushMsgErr           = 5603 // 推送消息错误
	IOSBackgroundPushErr = 5604 // iOS 后台推送错误

)
