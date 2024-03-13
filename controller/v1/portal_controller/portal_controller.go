package portal_controller

import (
	"framework-gin/common/function"
	"framework-gin/pojo/request"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons"
	"net/http"
)

// Test godoc
// @Summary Test
// @Description Test
// @Tags test
// @Accept  json
// @Produce  json
// @Router /api/v1/test [post]
// @param data body request.Test true "request.Test"
// @Success 200 {object} response.Test
func Test(c *gin.Context) {
	input := request.Test{}
	if code, err := function.BindAndValid(&input, c); err != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(code, err.Error(), input.RequestId))
	}

	if out, code, err := portalService.Test(input); err != nil {
		c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
	} else {
		c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
	}
}
