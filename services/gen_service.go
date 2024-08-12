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

	UserVersionCreate(info request.UserVersionCreate) (out response.UserVersionCreate, code commons.ResponseCode, err error)
	UserVersionDelete(info request.UserVersionDelete) (out response.UserVersionDelete, code commons.ResponseCode, err error)
	UserVersionUpdate(info request.UserVersionUpdate) (out response.UserVersionUpdate, code commons.ResponseCode, err error)
	UserVersionList(info request.UserVersionList) (out response.UserVersionList, code commons.ResponseCode, err error)

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
// --------------------UserVersion service layer implementation--------------------
// ================================================================================

func (g genServiceImp) UserVersionCreate(info request.UserVersionCreate) (out response.UserVersionCreate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserVersionDelete(info request.UserVersionDelete) (out response.UserVersionDelete, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserVersionUpdate(info request.UserVersionUpdate) (out response.UserVersionUpdate, code commons.ResponseCode, err error) {
	//todo
	return
}
func (g genServiceImp) UserVersionList(info request.UserVersionList) (out response.UserVersionList, code commons.ResponseCode, err error) {
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
