package errs

import "github.com/qiafan666/gotato/commons/gcommon"

func ChineseCodeMsg() map[int]string {
	return gcommon.MapMerge(ChineseHttpCodeMsg, ChineseWsCodeMsg)
}

// http error
const (
	BusinessError = 10000 // 业务错误
)

// ChineseHttpCodeMsg http code	and msg
var ChineseHttpCodeMsg = map[int]string{
	10000: "业务错误",
}

// ws error
const (
	ConnOverMaxNumLimit = 20001 // 超过最大连接数限制
	ConnArgsErr         = 20002 // 连接参数错误
)

// ChineseWsCodeMsg ws code	and msg
var ChineseWsCodeMsg = map[int]string{
	20001: "超过最大连接数限制",
	20002: "连接参数错误",
}
