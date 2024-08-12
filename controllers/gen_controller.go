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

var genControllerIns *GenControllerImp
var genControllerInitOnce sync.Once

func NewGenControllerInstance() *GenControllerImp {

	genControllerInitOnce.Do(func() {
		genControllerIns = &GenControllerImp{
			genService: services.NewGenServiceInstance(),
		}
	})

	return genControllerIns
}

type GenControllerImp struct {
	genService services.GenService
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
func (g *GenControllerImp) UserCreate(c *gin.Context) {
	input := request.UserCreate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserCreate(input); err != nil {
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
func (g *GenControllerImp) UserDelete(c *gin.Context) {
	input := request.UserDelete{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserDelete(input); err != nil {
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
func (g *GenControllerImp) UserUpdate(c *gin.Context) {
	input := request.UserUpdate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserUpdate(input); err != nil {
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
func (g *GenControllerImp) UserList(c *gin.Context) {
	input := request.UserList{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserList(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// ================================================================================
// -----------------------------UserVersion controller-----------------------------
// ================================================================================

// UserVersionCreate
// @Summary 创建接口
// @Description UserVersion创建
// @Tags UserVersion
// @Accept json
// @Produce json
// @Router /v1/userVersion/create [post]
// @param data body request.UserVersionCreate true "UserVersion创建请求参数"
// @Success 200 {object} response.UserVersionCreate "UserVersion创建返回结果"
func (g *GenControllerImp) UserVersionCreate(c *gin.Context) {
	input := request.UserVersionCreate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserVersionCreate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserVersionDelete
// @Summary 删除接口
// @Description UserVersion删除
// @Tags UserVersion
// @Accept json
// @Produce json
// @Router /v1/userVersion/delete [post]
// @param data body request.UserVersionDelete true "UserVersion删除请求参数"
// @Success 200 {object} response.UserVersionDelete "UserVersion删除返回结果"
func (g *GenControllerImp) UserVersionDelete(c *gin.Context) {
	input := request.UserVersionDelete{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserVersionDelete(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserVersionUpdate
// @Summary 更新接口
// @Description UserVersion更新
// @Tags UserVersion
// @Accept json
// @Produce json
// @Router /v1/userVersion/update [post]
// @param data body request.UserVersionUpdate true "UserVersion更新请求参数"
// @Success 200 {object} response.UserVersionUpdate "UserVersion更新返回结果"
func (g *GenControllerImp) UserVersionUpdate(c *gin.Context) {
	input := request.UserVersionUpdate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserVersionUpdate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// UserVersionList
// @Summary 列表接口
// @Description UserVersion列表
// @Tags UserVersion
// @Accept json
// @Produce json
// @Router /v1/userVersion/list [post]
// @param data body request.UserVersionList true "UserVersion列表请求参数"
// @Success 200 {object} response.UserVersionList "UserVersion列表返回结果"
func (g *GenControllerImp) UserVersionList(c *gin.Context) {
	input := request.UserVersionList{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.UserVersionList(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// ================================================================================
// -------------------------------Version controller-------------------------------
// ================================================================================

// VersionCreate
// @Summary 创建接口
// @Description Version创建
// @Tags Version
// @Accept json
// @Produce json
// @Router /v1/version/create [post]
// @param data body request.VersionCreate true "Version创建请求参数"
// @Success 200 {object} response.VersionCreate "Version创建返回结果"
func (g *GenControllerImp) VersionCreate(c *gin.Context) {
	input := request.VersionCreate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.VersionCreate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// VersionDelete
// @Summary 删除接口
// @Description Version删除
// @Tags Version
// @Accept json
// @Produce json
// @Router /v1/version/delete [post]
// @param data body request.VersionDelete true "Version删除请求参数"
// @Success 200 {object} response.VersionDelete "Version删除返回结果"
func (g *GenControllerImp) VersionDelete(c *gin.Context) {
	input := request.VersionDelete{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.VersionDelete(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// VersionUpdate
// @Summary 更新接口
// @Description Version更新
// @Tags Version
// @Accept json
// @Produce json
// @Router /v1/version/update [post]
// @param data body request.VersionUpdate true "Version更新请求参数"
// @Success 200 {object} response.VersionUpdate "Version更新返回结果"
func (g *GenControllerImp) VersionUpdate(c *gin.Context) {
	input := request.VersionUpdate{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.VersionUpdate(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}

// VersionList
// @Summary 列表接口
// @Description Version列表
// @Tags Version
// @Accept json
// @Produce json
// @Router /v1/version/list [post]
// @param data body request.VersionList true "Version列表请求参数"
// @Success 200 {object} response.VersionList "Version列表返回结果"
func (g *GenControllerImp) VersionList(c *gin.Context) {
	input := request.VersionList{}
	if bindCode, bindErr := function.BindAndValid(&input, c); bindErr != nil {
		c.JSON(http.StatusOK, commons.BuildFailedWithMsg(bindCode, bindErr.Error(), input.RequestId))
	} else {
		if out, code, err := g.genService.VersionList(input); err != nil {
			c.JSON(http.StatusOK, commons.BuildFailed(code, input.Language, input.RequestId))
		} else {
			c.JSON(http.StatusOK, commons.BuildSuccess(out, input.Language, input.RequestId))
		}
	}
}
