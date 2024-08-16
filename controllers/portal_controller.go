package controllers

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

// ================================================================================
// --------------------------------User controller---------------------------------
// ================================================================================

// UserCreate
// @Summary 创建接口
// @Description User创建
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/user/create [post]
// @param data body request.UserCreate true "User创建请求参数"
// @Success 200 {object} response.UserCreate "User创建返回结果"
func (g *PortalControllerImp) UserCreate(c *gin.Context) {
	input := request.UserCreate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.portalService.UserCreate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserDelete
// @Summary 删除接口
// @Description User删除
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/user/delete [post]
// @param data body request.UserDelete true "User删除请求参数"
// @Success 200 {object} response.UserDelete "User删除返回结果"
func (g *PortalControllerImp) UserDelete(c *gin.Context) {
	input := request.UserDelete{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.portalService.UserDelete(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserUpdate
// @Summary 更新接口
// @Description User更新
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/user/update [post]
// @param data body request.UserUpdate true "User更新请求参数"
// @Success 200 {object} response.UserUpdate "User更新返回结果"
func (g *PortalControllerImp) UserUpdate(c *gin.Context) {
	input := request.UserUpdate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.portalService.UserUpdate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserList
// @Summary 列表接口
// @Description User列表
// @Tags User
// @Accept json
// @Produce json
// @Router /v1/user/list [post]
// @param data body request.UserList true "User列表请求参数"
// @Success 200 {object} response.UserList "User列表返回结果"
func (g *PortalControllerImp) UserList(c *gin.Context) {
	input := request.UserList{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.portalService.UserList(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}
