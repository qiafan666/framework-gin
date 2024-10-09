package errs

import (
	"github.com/qiafan666/gotato/commons/gerr"
)

var (
	ErrArgs                = gerr.NewCodeError(ArgsError, "ArgsError")
	ErrTokenInvalid        = gerr.NewCodeError(TokenInvalidError, "TokenInvalidError") //
	ErrConnOverMaxNumLimit = gerr.NewCodeError(ConnOverMaxNumLimit, "ConnOverMaxNumLimit")
	ErrConnArgsErr         = gerr.NewCodeError(ConnArgsErr, "args err, need token, sendID, platformID")
	ErrInvalidRequest      = gerr.NewCodeError(InvalidRequestError, "InvalidRequest")
)
