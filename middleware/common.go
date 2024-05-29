package middleware

import (
	"framework-gin/common"
	"framework-gin/common/function"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/v2/middleware"
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
