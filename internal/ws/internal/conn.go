package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin"
	"net/http"
	"time"
)

type ConnInterface interface {
	// Close 关闭此连接
	Close() error
	// WriteMessage 写入消息到连接，messageType 表示数据类型，可以设置为二进制(2)或文本(1)
	WriteMessage(messageType int, message []byte) error
	// ReadMessage 从连接中读取消息
	ReadMessage() (int, []byte, error)
	// SetReadDeadline 设置读取超时时间，当读取超时后，将返回错误
	SetReadDeadline(timeout time.Duration) error
	// SetWriteDeadline 设置发送消息时的写入超时时间，当写入超时后，将返回错误
	SetWriteDeadline(timeout time.Duration) error
	// Dial 尝试拨号连接，url 必须设置认证参数，header 可以控制数据压缩
	Dial(urlStr string, requestHeader http.Header) (*http.Response, error)
	// IsNil 判断当前长连接的连接是否为空
	IsNil() bool
	// SetConnNil 将当前长连接的连接设置为空
	SetConnNil()
	// SetReadLimit 设置从对等端读取消息的最大大小（字节数）
	SetReadLimit(limit int64)
	// SetPongHandler 设置 Pong 消息的处理器
	SetPongHandler(handler PingPongHandler)
	// SetPingHandler 设置 Ping 消息的处理器
	SetPingHandler(handler PingPongHandler)
	// GenerateLongConn 检查当前连接和发送时的连接是否相同，生成长连接
	GenerateLongConn(w http.ResponseWriter, r *http.Request) error
}

type WebSocket struct {
	conn             *websocket.Conn
	handshakeTimeout time.Duration
	writeBufferSize  int
}

func newGWebSocket(handshakeTimeout time.Duration, wbs int) *WebSocket {
	return &WebSocket{handshakeTimeout: handshakeTimeout, writeBufferSize: wbs}
}

func (w *WebSocket) Close() error {
	return w.conn.Close()
}

func (w *WebSocket) GenerateLongConn(res http.ResponseWriter, req *http.Request) error {
	upgrade := &websocket.Upgrader{
		HandshakeTimeout: w.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	if w.writeBufferSize > 0 { // default is 4kb.
		upgrade.WriteBufferSize = w.writeBufferSize
	}

	conn, err := upgrade.Upgrade(res, req, nil)
	if err != nil {
		// The upgrader.Upgrade method usually returns enough error messages to diagnose problems that may occur during the upgrade
		return gerr.WrapMsg(err, "GenerateLongConn: WebSocket upgrade failed")
	}
	w.conn = conn
	return nil
}

func (w *WebSocket) WriteMessage(messageType int, message []byte) error {
	return w.conn.WriteMessage(messageType, message)
}

func (w *WebSocket) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WebSocket) SetReadDeadline(timeout time.Duration) error {
	return w.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (w *WebSocket) SetWriteDeadline(timeout time.Duration) error {
	if timeout <= 0 {
		return gerr.New("timeout must be greater than 0")
	}

	if err := w.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return gerr.WrapMsg(err, "GWebSocket.SetWriteDeadline failed")
	}
	return nil
}

func (w *WebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err != nil {
		return httpResp, gerr.WrapMsg(err, "GWebSocket.Dial failed", "url", urlStr)
	}
	w.conn = conn
	return httpResp, nil
}

func (w *WebSocket) IsNil() bool {
	if w.conn != nil {
		return false
	}
	return true
}

func (w *WebSocket) SetConnNil() {
	w.conn = nil
}

func (w *WebSocket) SetReadLimit(limit int64) {
	w.conn.SetReadLimit(limit)
}

func (w *WebSocket) SetPongHandler(handler PingPongHandler) {
	w.conn.SetPongHandler(handler)
}

func (w *WebSocket) SetPingHandler(handler PingPongHandler) {
	w.conn.SetPingHandler(handler)
}

func (w *WebSocket) RespondWithError(err error, res http.ResponseWriter, req *http.Request) error {
	if err = w.GenerateLongConn(res, req); err != nil {
		return err
	}
	data, err := json.Marshal(ggin.ParseError(err))
	if err != nil {
		_ = w.Close()
		return gerr.WrapMsg(err, "json marshal failed")
	}

	if err = w.WriteMessage(MessageText, data); err != nil {
		_ = w.Close()
		return gerr.WrapMsg(err, "WriteMessage failed")
	}
	_ = w.Close()
	return nil
}

func (w *WebSocket) RespondWithSuccess() error {
	data, err := json.Marshal(ggin.ApiSuccess(nil, "", "init suc"))
	if err != nil {
		_ = w.Close()
		return gerr.WrapMsg(err, "json marshal failed")
	}

	if err = w.WriteMessage(MessageText, data); err != nil {
		_ = w.Close()
		return gerr.WrapMsg(err, "WriteMessage failed")
	}
	return nil
}
