package services

import (
	"framework-gin/dao"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons"
	"sync"
)

// GenService service layer interface
type GenService interface {
	UserCreate(info request.UserCreate) (out response.UserCreate, code commons.ResponseCode, err error)
	UserDelete(info request.UserDelete) (out response.UserDelete, code commons.ResponseCode, err error)
	UserUpdate(info request.UserUpdate) (out response.UserUpdate, code commons.ResponseCode, err error)
	UserList(info request.UserList) (out response.UserList, code commons.ResponseCode, err error)
	VersionCreate(info request.VersionCreate) (out response.VersionCreate, code commons.ResponseCode, err error)
	VersionDelete(info request.VersionDelete) (out response.VersionDelete, code commons.ResponseCode, err error)
	VersionUpdate(info request.VersionUpdate) (out response.VersionUpdate, code commons.ResponseCode, err error)
	VersionList(info request.VersionList) (out response.VersionList, code commons.ResponseCode, err error)
}

var genServiceIns *genServiceImp
var genServiceInitOnce sync.Once

func NewGenServiceInstance() GenService {

	genServiceInitOnce.Do(func() {
		genServiceIns = &genServiceImp{
			dao: dao.Instance(),
		}
	})

	return genServiceIns
}

type genServiceImp struct {
	dao dao.Dao
}

// ================================================================================
// -----------------------User service layer implementation------------------------
// ================================================================================

func (g genServiceImp) UserCreate(info request.UserCreate) (out response.UserCreate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserDelete(info request.UserDelete) (out response.UserDelete, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserUpdate(info request.UserUpdate) (out response.UserUpdate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserList(info request.UserList) (out response.UserList, code commons.ResponseCode, err error) {
	//todo
	return
}

// ================================================================================
// ----------------------Version service layer implementation----------------------
// ================================================================================

func (g genServiceImp) VersionCreate(info request.VersionCreate) (out response.VersionCreate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) VersionDelete(info request.VersionDelete) (out response.VersionDelete, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) VersionUpdate(info request.VersionUpdate) (out response.VersionUpdate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) VersionList(info request.VersionList) (out response.VersionList, code commons.ResponseCode, err error) {
	//todo
	return
}
