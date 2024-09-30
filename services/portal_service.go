package services

import (
	"framework-gin/dao"
	"framework-gin/model"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/commons/gcommon"
	"gorm.io/gorm"
	"sync"
)

// PortalService service layer interface
type PortalService interface {
	UserCreate(info request.UserCreate) (out response.UserCreate, code commons.ResponseCode, err error)
	UserDelete(info request.UserDelete) (out response.UserDelete, code commons.ResponseCode, err error)
	UserUpdate(info request.UserUpdate) (out response.UserUpdate, code commons.ResponseCode, err error)
	UserList(info request.UserList) (out response.UserList, code commons.ResponseCode, err error)
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

func (g portalServiceImp) UserCreate(info request.UserCreate) (out response.UserCreate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g portalServiceImp) UserDelete(info request.UserDelete) (out response.UserDelete, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g portalServiceImp) UserUpdate(info request.UserUpdate) (out response.UserUpdate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g portalServiceImp) UserList(info request.UserList) (out response.UserList, code commons.ResponseCode, err error) {

	count, err := g.dao.WithContext(info.Ctx).Count(model.User{}, nil, nil)
	if err != nil {
		return response.UserList{}, 0, err
	}

	var users []model.User
	err = g.dao.WithContext(info.Ctx).Find([]string{}, nil, func(db *gorm.DB) *gorm.DB {
		return db.Scopes(gcommon.Paginate(info.CurrentPage, info.PageCount))
	}, &users)
	if err != nil {
		return response.UserList{}, 0, err
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
