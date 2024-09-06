package common

import "github.com/qiafan666/gotato/gconfig"

var DevEnv bool

func init() {
	if gconfig.SC.SConfigure.Profile == "dev" {
		DevEnv = true
	}
}

// ctx value enum
const (
	BaseRequest      = "base_request"
	BaseTokenRequest = "base_token_request"
)

// header value enum
const (
	HeaderAuthorization                 = "Authorization"
	HeaderLanguage                      = "Language"
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
