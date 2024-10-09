package internal

import (
	"github.com/qiafan666/gotato/commons/ggin"
)

func httpError(ctx *UserConnContext, err error) {
	ggin.HttpError(ctx.RespWriter, err)
}
