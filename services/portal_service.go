package services

import (
	"framework-gin/dao"
	"framework-gin/model"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/glog"
	"gorm.io/gorm"
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
	count, err := g.dao.WithContext(info.Ctx).Count(model.User{}, nil, nil)
	if err != nil {
		glog.Slog.ErrorKVs(info.Ctx, "Count error", "err", err)
		return out, gerr.NewLang(gerr.UnKnowError, info.Language, info.RequestId)
	}

	var users []model.User
	err = g.dao.WithContext(info.Ctx).Find([]string{}, nil, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(gcommon.Paginate(info.CurrentPage, info.PageCount))
	}, &users)
	if err != nil {
		glog.Slog.ErrorKVs(info.Ctx, "Find error", "err", err)
		return out, gerr.NewLang(gerr.UnKnowError, info.Language, info.RequestId)
	}

	out.UserList = gcommon.SliceConvert(users, func(user model.User) response.User {
		return response.User{
			UUID:        user.UUID,
			Name:        user.Name,
			Age:         user.Age,
			CreatedTime: user.CreatedTime,
		}
	})
	out.CurrentPage = info.CurrentPage
	out.PageCount = info.PageCount
	out.Count = count
	return
}
