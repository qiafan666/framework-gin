package internal

import (
	"context"
	"framework-gin/pkg/common"
	"github.com/qiafan666/gotato/commons/gcast"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/gid"
	"net/http"
	"strconv"
)

type BaseReq struct {
	RequestId string
	GrpId     uint8
	CmdId     uint8
}

type UserConnContext struct {
	BaseReq
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
	Uuid       string
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
		Uuid:       gcast.ToString(req.Header.Get(common.HeaderUUid)),
	}
	x.IsCompress = x.GetCompression()
	x.Language = x.GetLanguage()
	return x
}

func (c *UserConnContext) GetRemoteAddr() string {
	return c.RemoteAddr
}

func (c *UserConnContext) GetHeader(key string) (string, bool) {
	var value string
	if value = c.Req.Header.Get(key); value == "" {
		return value, false
	}
	return value, true
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

func (c *UserConnContext) ParseEssentialArgs() error {
	authToken, _ := c.GetHeader(common.HeaderAuthorization)
	uuid, _ := c.GetHeader(common.HeaderUUid)
	if authToken == "" && uuid == "" {
		return gerr.NewLang(gerr.ParameterError, c.Language).WithDetail("authToken and uuid are empty")
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

// Trace append当前请求的trace信息到context中
func (c *UserConnContext) Trace() context.Context {
	return AppendTraceCtx(c.TraceCtx, c.BaseReq)
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
	if traceId, ok := ctx.Value("request_id").(string); ok {
		slice := gcommon.Str2Slice(traceId, "-")
		return slice[0], slice[1], slice[2], slice[3], slice[4]
	} else {
		return "", "", "", "", ""
	}
}

// SetTraceCtx platform-connID-remoteAddr-userID
func SetTraceCtx(values []any) context.Context {
	return gcommon.SetRequestId(gcommon.Slice2Str(values, "-"))
}

// AppendTraceCtx platform-connID-remoteAddr-userID-requestID-grp-cmd
func AppendTraceCtx(ctx context.Context, baseReq BaseReq) context.Context {
	return gcommon.SetRequestId(gcommon.GetRequestId(ctx) + "-" +
		gcommon.Slice2Str([]any{baseReq.RequestId, baseReq.GrpId, baseReq.CmdId}, "-"))
}
