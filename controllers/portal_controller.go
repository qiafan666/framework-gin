package controllers

import (
	"framework-gin/common/function"
	"framework-gin/pojo/request"
	"framework-gin/services"
	"github.com/gin-gonic/gin"
	"github.com/qiafan666/gotato/commons/ggin"
)

type PortalController struct {
	portalService services.IPortalService
}

func NewPortalController() *PortalController {
	return &PortalController{
		portalService: services.NewPortalServiceInstance(),
	}
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
func (p *PortalController) UserCreate(c *gin.Context) {
	input := request.UserCreate{}
	if bindErr := function.BindAndValid(&input, c); bindErr != nil {
		ggin.GinError(c, bindErr)
	} else {
		if out, err := p.portalService.UserCreate(input); err != nil {
			ggin.GinError(c, err)
		} else {
			ggin.GinSuccess(c, out)
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
func (p *PortalController) UserDelete(c *gin.Context) {
	input := request.UserDelete{}
	if bindErr := function.BindAndValid(&input, c); bindErr != nil {
		ggin.GinError(c, bindErr)
	} else {
		if out, err := p.portalService.UserDelete(input); err != nil {
			ggin.GinError(c, err)
		} else {
			ggin.GinSuccess(c, out)
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
func (p *PortalController) UserUpdate(c *gin.Context) {
	input := request.UserUpdate{}
	if bindErr := function.BindAndValid(&input, c); bindErr != nil {
		ggin.GinError(c, bindErr)
	} else {
		if out, err := p.portalService.UserUpdate(input); err != nil {
			ggin.GinError(c, err)
		} else {
			ggin.GinSuccess(c, out)
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
func (p *PortalController) UserList(c *gin.Context) {
	input := request.UserList{}
	if bindErr := function.BindAndValid(&input, c); bindErr != nil {
		ggin.GinError(c, bindErr)
	} else {
		if out, err := p.portalService.UserList(input); err != nil {
			ggin.GinError(c, err)
		} else {
			ggin.GinSuccess(c, out)
		}
	}
}
