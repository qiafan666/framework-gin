package portal_controller

import (
	"github.com/gin-gonic/gin"
)

func PortalControllerInit(routeGroup gin.IRoutes) {
	routeGroup.POST("/portal/test", Test)
}
