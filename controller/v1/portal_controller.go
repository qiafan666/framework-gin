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

var portalControllerIns *PortalControllerImp
var portalControllerInitOnce sync.Once

func NewPortalControllerInstance() *PortalControllerImp {
	portalControllerInitOnce.Do(func() {
		portalControllerIns = &PortalControllerImp{
			portalService: services.NewPortalServiceInstance(),
		}
	})
	return portalControllerIns
}

type PortalControllerImp struct {
	portalService services.PortalService
}

// Test godoc
// @Summary Test
// @Description Test
// @Tags test
// @Accept  json
// @Produce  json
// @Router /v1/portal/test [post]
// @param data body request.Test true "request.Test"
// @Success 200 {object} response.Test
func (p *PortalControllerImp) Test(c *gin.Context) {
	input := request.Test{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := p.portalService.Test(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}
