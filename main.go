package main

import (
	_ "framework-gin/docs"
	"framework-gin/internal/http"
	"framework-gin/internal/ws"
	"framework-gin/pkg/common/errs"
	"github.com/qiafan666/gotato"
	"github.com/qiafan666/gotato/commons/gerr"
	_ "net/http/pprof"
)

// @title framework API Document
// @description framework API Document
// @version 1
// @schemes http
// @produce json
// @consumes json
func main() {
	server := gotato.GetGotato()
	server.RegisterErrorCodeAndMsg(gerr.MsgLanguageChinese, errs.ChineseCodeMsg())
	server.StartServer(gotato.GinService, gotato.DatabaseService, gotato.RedisService)
	http.Register(server.App())
	ws.Register(server.GetCtx(), server.App())
	server.WaitClose()
}
