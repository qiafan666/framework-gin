package internal

import (
	"context"
	"framework-gin/common"
	"framework-gin/ws/constant"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/gid"
	"github.com/qiafan666/gotato/commons/glog"
	"net/http"
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
	Language   string
	PlatformID int
	IsCompress bool
	TraceCtx   context.Context
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

	x := &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: gcommon.RemoteIP(req),
		ConnID:     gcast.ToString(gid.RandID()),
		PlatformID: gcast.ToInt(req.Header.Get(common.HeaderPlatformID)),
	}
	x.IsCompress = x.GetCompression()
	x.Language = x.GetLanguage()
	return x
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

func (c *UserConnContext) GetLanguage() string {
	headerLanguage := c.Req.Header.Get(common.HeaderLanguage)

	if headerLanguage == "" {
		return gerr.DefaultLanguage
	}
	if headerLanguage != gerr.MsgLanguageChinese && headerLanguage != gerr.MsgLanguageEnglish {
		return gerr.DefaultLanguage
	}
	return headerLanguage
}

func (c *UserConnContext) GetPlatformID() int {
	return gcast.ToInt(c.Req.Header.Get(common.HeaderPlatformID))
}

func (c *UserConnContext) GetToken() string {
	return c.Req.Header.Get(common.HeaderAuthorization)
}

func (c *UserConnContext) GetCompression() bool {
	compression, exists := c.GetHeader(common.HeaderCompression)
	if exists && compression == common.CompressionGzip {
		return true
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
		return gerr.NewLang(gerr.ParameterError, c.Language).WithDetail("auth token is empty")
	}
	platformIDStr, exists := c.GetHeader(common.HeaderPlatformID)
	if !exists {
		return gerr.NewLang(gerr.ParameterError, c.Language).WithDetail("platformID is empty")
	}
	_, err := strconv.Atoi(platformIDStr)
	if err != nil {
		return gerr.NewLang(gerr.ParameterError, c.Language).WithDetail("platformID is not a number")
	}
	return nil
}

// ------------------------ ws logic context ------------------------

func GetOpUserPlatform(ctx context.Context) string {
	GetCtxInfos(ctx)
	platform, _, _, _, _ := GetCtxInfos(ctx)
	return platform
}
func GetConnID(ctx context.Context) string {
	_, connID, _, _, _ := GetCtxInfos(ctx)
	return connID
}
func GetUserID(ctx context.Context) string {
	_, _, userID, _, _ := GetCtxInfos(ctx)
	return userID
}
func GetRequestID(ctx context.Context) string {
	_, _, _, requestID, _ := GetCtxInfos(ctx)
	return requestID
}
func GetRemoteAddr(ctx context.Context) string {
	_, _, _, _, addr := GetCtxInfos(ctx)
	return addr
}

func GetCtxInfosE(ctx context.Context) (platform, connID, userID, requestID, remoteAddr string, err error) {
	platform, connID, userID, requestID, remoteAddr = GetCtxInfos(ctx)
	if platform == "" {
		err = gerr.New("platform is empty")
		return
	}
	if connID == "" {
		err = gerr.New("connID is empty")
		return
	}
	if userID == "" {
		err = gerr.New("userID is empty")
	}

	if requestID == "" {
		err = gerr.New("requestID is empty")
	}

	if remoteAddr == "" {
		err = gerr.New("remoteAddr is empty")
	}
	return platform, connID, userID, requestID, remoteAddr, nil
}

func GetCtxInfos(ctx context.Context) (platform, connID, userID, requestID, remoteAddr string) {
	if traceId, ok := ctx.Value("trace_id").(string); ok {
		slice := gcommon.Str2Slice(traceId, "-")
		return slice[0], slice[1], slice[2], slice[3], slice[4]
	} else {
		return "", "", "", "", ""
	}
}

// SetTraceCtx platform-connID-remoteAddr-userID
func SetTraceCtx(values []any) context.Context {
	return glog.SetTraceId(gcommon.Slice2Str(values, "-"))
}

// AppendTraceCtx platform-connID-remoteAddr-userID-requestID-grp-cmd
func AppendTraceCtx(ctx context.Context, values []any) context.Context {
	return glog.SetTraceId(glog.GetTraceId(ctx) + "-" + gcommon.Slice2Str(values, "-"))
}
