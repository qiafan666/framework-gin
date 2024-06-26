package main

import (
	"framework-gin/common"
	"framework-gin/controllers"
	_ "framework-gin/docs"
	"github.com/qiafan666/gotato/commons"
	v2 "github.com/qiafan666/gotato/v2"
)

// @title framework API Document
// @description framework API Document
// @version 1
// @schemes http
// @produce json
// @consumes json
func main() {
	server := v2.GetGotatoInstance()
	server.RegisterErrorCodeAndMsg(commons.MsgLanguageEnglish, common.EnglishCodeMsg)
	server.StartServer(v2.GinService, v2.DatabaseService)
	controllers.RegisterRouter(server.App())
	server.WaitClose()
}
