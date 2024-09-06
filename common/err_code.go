package common

import (
	"github.com/qiafan666/gotato/commons"
)

// Error define the error code
const (
	Error = 1000
)

// ChineseCodeMsg local code and msg
var ChineseCodeMsg = map[commons.ResponseCode]string{
	1000: "业务错误",
}
