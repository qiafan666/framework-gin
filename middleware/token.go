package middleware

import (
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
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
		parseToken, err := jwt.Parse(ctx.Request.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.JWT.Secret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusOK, commons.BuildFailed(commons.TokenError, language, requestId))
			ctx.Abort()
			return
		}

		if _, ok = parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
			//TODO 查询用户信息
		} else {
			ctx.JSON(http.StatusOK, commons.BuildFailed(commons.TokenError, language, requestId))
			ctx.Abort()
			return
		}

	}

	ctx.Set(common.BaseTokenRequest, request.BaseTokenRequest{})

	ctx.Next()
}

func CreateToken(uuid string, iss string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": uuid,
		"iss":  iss,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24 * time.Duration(jwtConfig.JWT.Expire)).Unix(),
	}).SignedString([]byte(jwtConfig.JWT.Secret))
	if err != nil {
		return "", err
	}
	return token, nil
}
