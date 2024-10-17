package main

import (
	"framework-gin/common"
	"framework-gin/controllers"
	_ "framework-gin/docs"
	"framework-gin/ws"
	"github.com/qiafan666/gotato/commons"
	"github.com/qiafan666/gotato/v2"
)

// @title framework API Document
// @description framework API Document
// @version 1
// @schemes http
// @produce json
// @consumes json
func main() {
	server := v2.GetGotatoInstance()
	server.RegisterErrorCodeAndMsg(commons.MsgLanguageChinese, common.GetCodeMsg())
	server.StartServer(v2.GinService, v2.DatabaseService)
	controllers.RegisterRouter(server.App())
	ws.Register(server.App())
	server.WaitClose()
}
