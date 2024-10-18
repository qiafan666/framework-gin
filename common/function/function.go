package function

import (
	"context"
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"reflect"
)

// BindAndValid binds and validates data
func BindAndValid(entity interface{}, ctx *gin.Context) error {

	//set base request parameter
	object := reflect.ValueOf(entity)

	baseRequest := ctx.Keys[(common.BaseRequest)].(request.BaseRequest)
	elem := object.Elem()
	base := elem.FieldByName(common.BaseRequest)
	if base.Kind() != reflect.Invalid {
		base.Set(reflect.ValueOf(baseRequest))
	}

	baseTokenRequest, _ := ctx.Keys[(common.BaseTokenRequest)].(request.BaseTokenRequest)
	baseToken := elem.FieldByName(common.BaseTokenRequest)
	if baseToken.Kind() != reflect.Invalid {
		baseToken.Set(reflect.ValueOf(baseTokenRequest))
	}

	err := ctx.MustBindWith(entity, binding.JSON)
	if err != nil {
		glog.Slog.ErrorF(ctx.Value("ctx").(context.Context), "BindAndValid error: %v", err)
		return gerr.NewLang(gerr.ValidateError, baseRequest.Language, baseRequest.RequestId)
	}

	if err = gcommon.Validate(entity); err != nil {
		glog.Slog.ErrorF(ctx.Value("ctx").(context.Context), "Validate error: %v", err)
		return gerr.NewLang(gerr.ValidateError, baseRequest.Language, baseRequest.RequestId)
	}

	return nil
}

func GetTraceId(ctx *gin.Context) string {
	if traceId, ok := ctx.Value("trace_id").(string); ok {
		return traceId
	} else {
		return ""
	}
}
func GetCtx(ctx *gin.Context) context.Context {
	if v, ok := ctx.Value("ctx").(context.Context); ok {
		return v
	} else {
		return context.Background()
	}
}

var WsCtx = context.WithValue(context.Background(), "trace_id", "ws_server")
