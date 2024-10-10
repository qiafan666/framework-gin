package internal

import (
	"context"
	"framework-gin/common"
	"framework-gin/ws/constant"
	"framework-gin/ws/errs"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gencrypt"
	"github.com/qiafan666/gotato/commons/glog"
	"github.com/qiafan666/gotato/commons/gtime"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type UserConnContext struct {
	RespWriter http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	RemoteAddr string
	ConnID     string
	Ctx        context.Context
}

func (c *UserConnContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *UserConnContext) Done() <-chan struct{} {
	return nil
}

func (c *UserConnContext) Err() error {
	return nil
}

func (c *UserConnContext) Value(key any) any {
	switch key {
	case constant.OpUserID:
		return c.GetUserID()
	case constant.RequestID:
		return c.GetRequestID()
	case constant.ConnID:
		return c.GetConnID()
	case constant.OpUserPlatform:
		return constant.PlatformIDToName(gcast.ToInt(c.GetPlatformID()))
	case constant.RemoteAddr:
		return c.RemoteAddr
	default:
		return ""
	}
}

func newContext(respWriter http.ResponseWriter, req *http.Request) *UserConnContext {

	ctx := glog.SetTraceId(constant.PlatformIDToName(gcast.ToInt(req.Header.Get(common.HeaderPlatformID))) + "-" + req.Header.Get(common.HeaderSendID))
	return &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		ConnID:     gencrypt.Md5(req.RemoteAddr + "_" + strconv.Itoa(int(gtime.GetCurrentTimestampByMill()))),
		Ctx:        ctx,
	}
}

func newTempContext() *UserConnContext {
	return &UserConnContext{
		Req: &http.Request{URL: &url.URL{}},
	}
}

func (c *UserConnContext) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *UserConnContext) Query(key string) (string, bool) {
	var value string
	if value = c.Req.Header.Get(key); value == "" {
		return value, false
	}
	return value, true
}

func (c *UserConnContext) GetHeader(key string) (string, bool) {
	var value string
	if value = c.Req.Header.Get(key); value == "" {
		return value, false
	}
	return value, true
}

func (c *UserConnContext) SetHeader(key, value string) {
	c.RespWriter.Header().Set(key, value)
}

func (c *UserConnContext) ErrReturn(error string, code int) {
	http.Error(c.RespWriter, error, code)
}

func (c *UserConnContext) GetConnID() string {
	return c.ConnID
}

func (c *UserConnContext) GetUserID() string {
	return c.Req.Header.Get(common.HeaderSendID)
}

func (c *UserConnContext) GetPlatformID() string {
	return c.Req.Header.Get(common.HeaderPlatformID)
}

func (c *UserConnContext) GetRequestID() string {
	return c.Req.Header.Get(common.RequestID)
}

func (c *UserConnContext) SetRequestID(requestID string) {
	c.Req.Header.Set(common.RequestID, requestID)
}

func (c *UserConnContext) GetToken() string {
	return c.Req.Header.Get(common.HeaderAuthorization)
}

func (c *UserConnContext) GetCompression() bool {
	compression, exists := c.Query(common.HeaderCompression)
	if exists && compression == common.CompressionGzip {
		return true
	} else {
		compression, exists = c.GetHeader(common.HeaderCompression)
		if exists && compression == common.CompressionGzip {
			return true
		}
	}
	return false
}

func (c *UserConnContext) ShouldSendResp() bool {
	errResp, exists := c.GetHeader(common.HeaderSendResponse)
	if exists {
		b, err := strconv.ParseBool(errResp)
		if err != nil {
			return false
		} else {
			return b
		}
	}
	return false
}

func (c *UserConnContext) SetToken(token string) {
	c.Req.Header.Set(common.HeaderAuthorization, token)
}

func (c *UserConnContext) ParseEssentialArgs() error {
	_, exists := c.GetHeader(common.HeaderAuthorization)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("token is empty")
	}
	_, exists = c.GetHeader(common.HeaderSendID)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("sendID is empty")
	}
	platformIDStr, exists := c.GetHeader(common.HeaderPlatformID)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("platformID is empty")
	}
	_, err := strconv.Atoi(platformIDStr)
	if err != nil {
		return errs.ErrConnArgsErr.WrapMsg("platformID is not int")

	}
	return nil
}
