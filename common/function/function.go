package function

import (
	"framework-gin/common"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/commons/utils"
	"reflect"
)

// BindAndValid binds and validates data
func BindAndValid(entity interface{}, ctx *gin.Context) (commons.ResponseCode, error) {

	//set base request parameter
	object := reflect.ValueOf(entity)

	baseRequest := ctx.Keys[(common.BaseRequest)].(request.BaseRequest)
	elem := object.Elem()
	base := elem.FieldByName("BaseRequest")
	if base.Kind() != reflect.Invalid {
		base.Set(reflect.ValueOf(baseRequest))
	}

	baseTokenRequest, _ := ctx.Keys[(common.BaseTokenRequest)].(request.BaseTokenRequest)
	baseToken := elem.FieldByName("BaseTokenRequest")
	if baseToken.Kind() != reflect.Invalid {
		baseToken.Set(reflect.ValueOf(baseTokenRequest))
	}

	err := ctx.Bind(entity)
	if err != nil {
		return commons.ParameterError, err
	}

	if err = utils.Validate(entity); err != nil {
		return commons.ValidateError, err
	}

	return commons.OK, nil
}
