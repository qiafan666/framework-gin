package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/qiafan666/gotato/commons/gerr"
	"github.com/qiafan666/gotato/commons/ggin"
	"net/http"
	"time"
)

type LongConn interface {
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

type GWebSocket struct {
	conn             *websocket.Conn
	handshakeTimeout time.Duration
	writeBufferSize  int
}

func newGWebSocket(handshakeTimeout time.Duration, wbs int) *GWebSocket {
	return &GWebSocket{handshakeTimeout: handshakeTimeout, writeBufferSize: wbs}
}

func (d *GWebSocket) Close() error {
	return d.conn.Close()
}

func (d *GWebSocket) GenerateLongConn(w http.ResponseWriter, r *http.Request) error {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: d.handshakeTimeout,
		CheckOrigin:      func(r *http.Request) bool { return true },
	}
	if d.writeBufferSize > 0 { // default is 4kb.
		upgrader.WriteBufferSize = d.writeBufferSize
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// The upgrader.Upgrade method usually returns enough error messages to diagnose problems that may occur during the upgrade
		return gerr.WrapMsg(err, "GenerateLongConn: WebSocket upgrade failed")
	}
	d.conn = conn
	return nil
}

func (d *GWebSocket) WriteMessage(messageType int, message []byte) error {
	return d.conn.WriteMessage(messageType, message)
}

func (d *GWebSocket) ReadMessage() (int, []byte, error) {
	return d.conn.ReadMessage()
}

func (d *GWebSocket) SetReadDeadline(timeout time.Duration) error {
	return d.conn.SetReadDeadline(time.Now().Add(timeout))
}

func (d *GWebSocket) SetWriteDeadline(timeout time.Duration) error {
	if timeout <= 0 {
		return gerr.New("timeout must be greater than 0")
	}

	if err := d.conn.SetWriteDeadline(time.Now().Add(timeout)); err != nil {
		return gerr.WrapMsg(err, "GWebSocket.SetWriteDeadline failed")
	}
	return nil
}

func (d *GWebSocket) Dial(urlStr string, requestHeader http.Header) (*http.Response, error) {
	conn, httpResp, err := websocket.DefaultDialer.Dial(urlStr, requestHeader)
	if err != nil {
		return httpResp, gerr.WrapMsg(err, "GWebSocket.Dial failed", "url", urlStr)
	}
	d.conn = conn
	return httpResp, nil
}

func (d *GWebSocket) IsNil() bool {
	return d.conn == nil
	//
	// if d.conn != nil {
	// 	return false
	// }
	// return true
}

func (d *GWebSocket) SetConnNil() {
	d.conn = nil
}

func (d *GWebSocket) SetReadLimit(limit int64) {
	d.conn.SetReadLimit(limit)
}

func (d *GWebSocket) SetPongHandler(handler PingPongHandler) {
	d.conn.SetPongHandler(handler)
}

func (d *GWebSocket) SetPingHandler(handler PingPongHandler) {
	d.conn.SetPingHandler(handler)
}

func (d *GWebSocket) RespondWithError(err error, w http.ResponseWriter, r *http.Request) error {
	if err = d.GenerateLongConn(w, r); err != nil {
		return err
	}
	data, err := json.Marshal(ggin.ParseError(err))
	if err != nil {
		_ = d.Close()
		return gerr.WrapMsg(err, "json marshal failed")
	}

	if err = d.WriteMessage(MessageText, data); err != nil {
		_ = d.Close()
		return gerr.WrapMsg(err, "WriteMessage failed")
	}
	_ = d.Close()
	return nil
}

func (d *GWebSocket) RespondWithSuccess() error {
	data, err := json.Marshal(ggin.ApiSuccess(nil, "", "init suc"))
	if err != nil {
		_ = d.Close()
		return gerr.WrapMsg(err, "json marshal failed")
	}

	if err = d.WriteMessage(MessageText, data); err != nil {
		_ = d.Close()
		return gerr.WrapMsg(err, "WriteMessage failed")
	}
	return nil
}
