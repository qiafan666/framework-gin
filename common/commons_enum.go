package common

import (
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/config"
)

var DebugFlag bool

func init() {
	if config.SC.SConfigure.Profile == "dev" {
		DebugFlag = true
	}
}

// Error define the error code
const (
	Error = 1000
)

// ChineseCodeMsg local code and msg
var ChineseCodeMsg = map[commons.ResponseCode]string{
	1000: "业务错误",
}

// ctx value enum
const (
	BaseRequest      = "base_request"
	BaseTokenRequest = "base_token_request"
)
