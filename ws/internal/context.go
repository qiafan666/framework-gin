package internal

import (
	"framework-gin/ws/constant"
	"framework-gin/ws/errs"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gencrypt"
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
	case constant.OperationID:
		return c.GetOperationID()
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
	return &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		ConnID:     gencrypt.Md5(req.RemoteAddr + "_" + strconv.Itoa(int(gtime.GetCurrentTimestampByMill()))),
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
	if value = c.Req.URL.Query().Get(key); value == "" {
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
	return c.Req.URL.Query().Get(WsUserID)
}

func (c *UserConnContext) GetPlatformID() string {
	return c.Req.URL.Query().Get(PlatformID)
}

func (c *UserConnContext) GetOperationID() string {
	return c.Req.URL.Query().Get(OperationID)
}

func (c *UserConnContext) SetOperationID(operationID string) {
	values := c.Req.URL.Query()
	values.Set(OperationID, operationID)
	c.Req.URL.RawQuery = values.Encode()
}

func (c *UserConnContext) GetToken() string {
	return c.Req.URL.Query().Get(Token)
}

func (c *UserConnContext) GetCompression() bool {
	compression, exists := c.Query(Compression)
	if exists && compression == GzipCompressionProtocol {
		return true
	} else {
		compression, exists = c.GetHeader(Compression)
		if exists && compression == GzipCompressionProtocol {
			return true
		}
	}
	return false
}

func (c *UserConnContext) ShouldSendResp() bool {
	errResp, exists := c.Query(SendResponse)
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
	c.Req.URL.RawQuery = Token + "=" + token
}

func (c *UserConnContext) ParseEssentialArgs() error {
	_, exists := c.Query(Token)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("token is empty")
	}
	_, exists = c.Query(WsUserID)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("sendID is empty")
	}
	platformIDStr, exists := c.Query(PlatformID)
	if !exists {
		return errs.ErrConnArgsErr.WrapMsg("platformID is empty")
	}
	_, err := strconv.Atoi(platformIDStr)
	if err != nil {
		return errs.ErrConnArgsErr.WrapMsg("platformID is not int")

	}
	return nil
}
