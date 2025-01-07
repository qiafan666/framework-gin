package services

import (
	"framework-gin/common/errs"
	"framework-gin/dao"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"sync"
)

// PortalService service layer interface
type PortalService interface {
	UserCreate(info request.UserCreate) (out response.UserCreate, err error)
	UserDelete(info request.UserDelete) (out response.UserDelete, err error)
	UserUpdate(info request.UserUpdate) (out response.UserUpdate, err error)
	UserList(info request.UserList) (out response.UserList, err error)
}

var portalServiceIns *portalServiceImp
var portalServiceInitOnce sync.Once

func NewPortalServiceInstance() PortalService {

	portalServiceInitOnce.Do(func() {
		portalServiceIns = &portalServiceImp{
			dao: dao.Instance(),
		}
	})

	return portalServiceIns
}

type portalServiceImp struct {
	dao dao.Dao
}

// ================================================================================
// -----------------------User service layer implementation------------------------
// ================================================================================

func (g *portalServiceImp) UserCreate(info request.UserCreate) (out response.UserCreate, err error) {
	//todo
	return
}
func (g *portalServiceImp) UserDelete(info request.UserDelete) (out response.UserDelete, err error) {
	//todo
	return
}
func (g *portalServiceImp) UserUpdate(info request.UserUpdate) (out response.UserUpdate, err error) {
	//todo
	return
}
func (g *portalServiceImp) UserList(info request.UserList) (out response.UserList, err error) {
	glog.Slog.ErrorF(info.Ctx, "UserList not implemented")
	return out, gerr.NewLang(errs.BusinessError, info.Language)
}
