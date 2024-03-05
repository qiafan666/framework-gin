package middleware

import (
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato"
	"github.com/qiafan666/gotato/commons"
	"net/http"
)

var jwtConfig struct {
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
}

func init() {
	gotato.GetGotatoInstance().LoadCustomizeConfig(&jwtConfig)
}

var witheList = map[string]string{
	"/health":         "",
	"/v1/portal/test": "",
}

func CheckPortalAuth(ctx *gin.Context) {

	var language, requestId, phone, companyName string
	var userId, roleId int

	baseRequest := ctx.Keys[(common.BaseRequest)].(request.BaseRequest)
	language = baseRequest.Language
	requestId = baseRequest.RequestId
	//check white list
	if _, ok := witheList[ctx.Request.RequestURI]; !ok {

		//check jwt
		parseToken, err := jwt.Parse(ctx.Request.Header.Get("Authorization"), func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.JWT.Secret), nil
		})
		if err != nil {
			ctx.JSON(http.StatusOK, commons.BuildFailed(commons.TokenError, language, requestId))
			ctx.Abort()
			return
		}

		if _, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {

		} else {
			ctx.JSON(http.StatusOK, commons.BuildFailed(commons.TokenError, language, requestId))
			ctx.Abort()
			return
		}

	}

	ctx.Set(common.BaseTokenRequest, request.BaseTokenRequest{
		BaseID:      int64(userId),
		Phone:       phone,
		Role:        roleId,
		CompanyName: companyName,
	})

	ctx.Next()
}
