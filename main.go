package main

import (
	"framework-gin/common/errs"
	"framework-gin/controllers"
	_ "framework-gin/docs"
	"framework-gin/ws"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/v2"
	_ "net/http/pprof"
)

// @title framework API Document
// @description framework API Document
// @version 1
// @schemes http
// @produce json
// @consumes json
func main() {
	server := v2.GetGotato()
	server.ReadConfig()
	server.RegisterErrorCodeAndMsg(gerr.MsgLanguageChinese, errs.ChineseCodeMsg())
	server.StartServer(v2.GinService, v2.DatabaseService, v2.RedisService)
	controllers.RegisterRouter(server.App())
	ws.Register(server.GetCtx(), server.App())
	server.WaitClose()
}
