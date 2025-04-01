package services

import (
	"framework-gin/common/errs"
	"framework-gin/dao"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"github.com/qiafan666/gotato/commons/gredis"
	v2 "github.com/qiafan666/gotato/v2"
)

// IPortalService service layer interface
type IPortalService interface {
	UserCreate(info request.UserCreate) (out response.UserCreate, err error)
	UserDelete(info request.UserDelete) (out response.UserDelete, err error)
	UserUpdate(info request.UserUpdate) (out response.UserUpdate, err error)
	UserList(info request.UserList) (out response.UserList, err error)
}

type portalService struct {
	dao   dao.IDao
	redis *gredis.Client
}

func NewPortalServiceInstance() IPortalService {
	return &portalService{
		dao:   dao.New(),
		redis: gredis.SetRedis(v2.GetGotato().Redis("test")),
	}
}

// ================================================================================
// -----------------------User service layer implementation------------------------
// ================================================================================

func (p *portalService) UserCreate(info request.UserCreate) (out response.UserCreate, err error) {
	//todo
	return
}
func (p *portalService) UserDelete(info request.UserDelete) (out response.UserDelete, err error) {
	//todo
	return
}
func (p *portalService) UserUpdate(info request.UserUpdate) (out response.UserUpdate, err error) {
	//todo
	return
}
func (p *portalService) UserList(info request.UserList) (out response.UserList, err error) {
	glog.Slog.ErrorF(info.Ctx, "UserList not implemented")
	return out, gerr.NewLang(errs.BusinessError, info.Language)
}
