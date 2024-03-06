package controller

import (
	"framework-gin/controller/v1/portal_controller"
	"framework-gin/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(r *gin.Engine) {
	//default setting
	r.Use(func(context *gin.Context) {
		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
	}).Use(middleware.Common).GET("/health", func(context *gin.Context) {
		context.Status(200)
	})

	v1 := r.Group("/v1").Use(middleware.CheckPortalAuth)
	//注册controller
	portal_controller.PortalControllerInit(v1)

}
