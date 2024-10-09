package errs

import "github.com/qiafan666/gotato/commons/gerr"

var (
	ErrDatabase         = gerr.NewCodeError(DatabaseError, "DatabaseError")
	ErrNetwork          = gerr.NewCodeError(NetworkError, "NetworkError")
	ErrCallback         = gerr.NewCodeError(CallbackError, "CallbackError")
	ErrCallbackContinue = gerr.NewCodeError(CallbackError, "ErrCallbackContinue")

	ErrInternalServer = gerr.NewCodeError(ServerInternalError, "ServerInternalError")
	ErrArgs           = gerr.NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission   = gerr.NewCodeError(NoPermissionError, "NoPermissionError")
	ErrDuplicateKey   = gerr.NewCodeError(DuplicateKeyError, "DuplicateKeyError")
	ErrRecordNotFound = gerr.NewCodeError(RecordNotFoundError, "RecordNotFoundError")

	ErrUserIDNotFound    = gerr.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrRegisteredAlready = gerr.NewCodeError(RegisteredAlreadyError, "RegisteredAlreadyError")

	ErrData             = gerr.NewCodeError(DataError, "DataError")
	ErrTokenExpired     = gerr.NewCodeError(TokenExpiredError, "TokenExpiredError")
	ErrTokenInvalid     = gerr.NewCodeError(TokenInvalidError, "TokenInvalidError")         //
	ErrTokenMalformed   = gerr.NewCodeError(TokenMalformedError, "TokenMalformedError")     //
	ErrTokenNotValidYet = gerr.NewCodeError(TokenNotValidYetError, "TokenNotValidYetError") //
	ErrTokenUnknown     = gerr.NewCodeError(TokenUnknownError, "TokenUnknownError")         //
	ErrTokenKicked      = gerr.NewCodeError(TokenKickedError, "TokenKickedError")
	ErrTokenNotExist    = gerr.NewCodeError(TokenNotExistError, "TokenNotExistError") //

	ErrConnOverMaxNumLimit = gerr.NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")

	ErrConnArgsErr          = gerr.NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
	ErrPushMsgErr           = gerr.NewCodeError(PushMsgErr, "push msg err")
	ErrIOSBackgroundPushErr = gerr.NewCodeError(IOSBackgroundPushErr, "ios background push err")
)
