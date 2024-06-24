package controllers

import (
	"framework-gin/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	//default setting
	r.Use(middleware.Cors).
		Use(middleware.Common).
		Use(middleware.CheckToken).
		GET("/health", middleware.Health)
	//注册controller
	portalController := NewPortalControllerInstance()

	portalGroup := r.Group("v1/portal")
	{
		portalGroup.POST("/test", portalController.Test)
	}
}
