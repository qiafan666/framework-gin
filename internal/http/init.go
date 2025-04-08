package http

import (
	"framework-gin/internal/http/controllers"
	"framework-gin/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	//default setting
	r.Use(middleware.Cors).
		Use(middleware.Common).
		Use(middleware.CheckToken).
		GET("/health", middleware.Health)
	rootGroup := r.Group("v1")

	//注册controller
	portalController := controllers.NewPortalController()

	portalGroup := rootGroup.Group("portal")
	{
		//portalGroup.POST("/user/create", portalController.UserCreate)
		//portalGroup.POST("/user/delete", portalController.UserDelete)
		//portalGroup.POST("/user/update", portalController.UserUpdate)
		portalGroup.POST("/user/list", portalController.UserList)
	}
}
