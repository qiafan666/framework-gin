package portal_controller

import (
	"framework-gin/services"
	"github.com/gin-gonic/gin"
	"sync"
)

var once sync.Once
var portalService services.PortalService

func PortalControllerInit(routeGroup gin.IRoutes) {
	once.Do(func() {
		portalService = services.NewPortalServiceInstance()
	})

	routeGroup.POST("/portal/test", Test)
}
