package middleware

import (
	"context"
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/v2/middleware"
	"sync"
)

// body and api print
var blackList = []string{
	"/health",
	"/swagger/*",
	"/debug/pprof/*",
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
	c := ctx.Value("ctx").(context.Context)
	requestId := c.Value("trace_id").(string)
	ctx.Set(common.BaseRequest, request.BaseRequest{
		Ctx:       c,
		RequestId: requestId,
		Language:  language,
	})

	ctx.Next()
}
