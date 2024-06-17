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
	Error = 9000
)

// EnglishCodeMsg local code and msg

var EnglishCodeMsg = map[commons.ResponseCode]string{
	9000: "未知错误",
}

// ctx value enum
const (
	BaseRequest      = "base_request"
	BaseTokenRequest = "base_token_request"
)
