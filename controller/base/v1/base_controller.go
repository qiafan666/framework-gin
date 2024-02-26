package v1

import (
	"framework-gin/common/function"
	"framework-gin/pojo/request"
	"framework-gin/services"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"net/http"
	"sync"
)

var once sync.Once
var baseService services.BaseService

func init() {
	once.Do(func() {
		baseService = services.NewBaseServiceInstance()
	})
}

// PostTest godoc
// @Summary Test
// @Description Test
// @Tags test
// @Accept  json
// @Produce  json
// @Router /api/v1/test [post]
// @param data body request.Test true "request.Test"
// @Success 200 {object} response.Test
func PostTest(c *gin.Context) {
	input := request.Test{}
	if code, err := function.BindAndValid(&input, c); err != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(code, err.Error(), input.RequestId))
	}

	if out, code, err := baseService.Test(input); err != nil {
		c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
	} else {
		c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
	}

}
