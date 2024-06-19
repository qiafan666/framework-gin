package services

import (
	"framework-gin/dao"
	"framework-gin/pojo/request"
	"framework-gin/pojo/response"
	"github.com/qiafan666/gotato/commons"
	"sync"
)

// PortalService service layer interface
type PortalService interface {
	Test(info request.Test) (out response.Test, code commons.ResponseCode, err error)
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

func (p portalServiceImp) Test(info request.Test) (out response.Test, code commons.ResponseCode, err error) {
	return
}
