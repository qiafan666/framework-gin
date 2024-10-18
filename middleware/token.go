package middleware

import (
	"errors"
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/v2"
	"net/http"
	"time"
)

var jwtConfig struct {
	JWT struct {
		Secret string `yaml:"secret"`
		Expire int    `yaml:"expire"` //天数
	} `yaml:"jwt"`
}

func init() {
	v2.GetGotatoInstance().LoadCustomizeConfig(&jwtConfig)
}

var witheList = map[string]string{
	"/ws":                  "",
	"/health":              "",
	"/v1/portal/user/list": "",
}

func CheckToken(ctx *gin.Context) {

	//check white list
	if _, ok := witheList[ctx.Request.RequestURI]; !ok {
		var language, requestId string
		baseRequest := ctx.Keys[(common.BaseRequest)].(request.BaseRequest)
		language = baseRequest.Language
		requestId = baseRequest.RequestId

		//check jwt
		resultMap, err := ParseToken(ctx.Request.Header.Get(common.HeaderAuthorization))
		if err != nil || len(resultMap) == 0 {
			ctx.JSON(http.StatusOK, commons.BuildFailed(gerr.TokenError, language, requestId))
			ctx.Abort()
			return
		}

		//todo 获取用户信息
	}

	ctx.Set(common.BaseTokenRequest, request.BaseTokenRequest{})

	ctx.Next()
}

// CreateToken 创建token
func CreateToken(uuid string, iss string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		common.TOKENUuid: uuid,
		common.TOKENIss:  iss,
		common.TOKENIat:  time.Now().Unix(),
		common.TOKENExp:  time.Now().Add(time.Hour * 24 * time.Duration(jwtConfig.JWT.Expire)).Unix(),
	}).SignedString([]byte(jwtConfig.JWT.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ParseToken 解析token
func ParseToken(token string) (result map[string]interface{}, err error) {

	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.JWT.Secret), nil
	})
	if err != nil {
		return result, err
	}
	if _, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		result[common.TOKENUuid] = parseToken.Claims
	} else {
		return result, errors.New("token parse error")
	}
	return result, nil
}
