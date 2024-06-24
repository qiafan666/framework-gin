package middleware

import (
	"framework-gin/common"
	"framework-gin/common/function"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/v2/middleware"
	"net/http"
	"sync"
)

// 拉黑的url不会被记录到日志中
var blackList = []string{
	"/health",
}

var once sync.Once

func init() {
	once.Do(func() {
		middleware.RegisterIgnoreRequest(blackList...)
	})
}

func Common(ctx *gin.Context) {

	//get language
	language := ctx.Request.Header.Get("Language")
	if language == "" {
		language = commons.DefaultLanguage
	}

	ctx.Set(common.BaseRequest, request.BaseRequest{
		Ctx:       function.GetCtx(ctx),
		RequestId: function.GetTraceId(ctx),
		Language:  language,
	})

	ctx.Next()
}

func Cors(ctx *gin.Context) {
	method := ctx.Request.Method

	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
	ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
	ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	ctx.Header("Access-Control-Allow-Credentials", "true")

	if method == "OPTIONS" {
		ctx.AbortWithStatus(http.StatusNoContent)
	}
	ctx.Next()
}

func Health(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
