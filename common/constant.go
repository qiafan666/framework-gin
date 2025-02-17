package common

import (
	"github.com/qiafan666/gotato/gconfig"
)

var DevEnv bool

func init() {
	if gconfig.SC.SConfigure.Profile == "dev" {
		DevEnv = true
	}
}

// ctx value enum
const (
	BaseRequest      = "BaseRequest"
	BaseTokenRequest = "BaseTokenRequest"
)

// header value enum
const (
	HeaderAuthorization                 = "Authorization"
	HeaderLanguage                      = "Language"
	HeaderPlatformID                    = "PlatformID"
	HeaderCompression                   = "Compression"
	HeaderSendResponse                  = "SendResponse"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
)

// token value enum
const (
	TOKENUuid = "uuid"
	TOKENIss  = "iss"
	TOKENIat  = "iat"
	TOKENExp  = "exp"
)

// 常量定义
const (
	CompressionGzip = "gzip"
)
