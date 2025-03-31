package ws

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"framework-gin/ws/proto/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/qiafan666/gotato/commons/gcompress"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

type Req struct {
	RequestId string `json:"request_id"   validate:"required"`
	GrpId     uint8  `json:"grp_id" validate:"required"` // 消息组id
	CmdId     uint8  `json:"cmd_id" validate:"required"` // 消息的ID
	Data      []byte `json:"data"`
}

func TestEncode(t *testing.T) {

	marshal, err := proto.Marshal(&pb.ReqHealth{Msg: "hello"})
	if err != nil {
		return
	}

	req := &Req{
		RequestId: "abcdefg",
		GrpId:     2,
		CmdId:     1,
		Data:      marshal,
	}

	gobEncoder := gcompress.NewGobEncoder()

	encode, err := gobEncoder.Encode(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(encode)
	t.Log(base64.StdEncoding.EncodeToString(encode))

	decodeString, err := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(encode))
	if err != nil {
		t.Error(err)
	}
	var decode *Req
	err = gobEncoder.Decode(decodeString, &decode)
	if err != nil {
		t.Error(err)
	}
	t.Log(decode)
	t.Log(string(decode.Data))
}

func TestJson(t *testing.T) {
	marshal, err := json.Marshal(&pb.ReqHealth{Msg: "hello"})
	if err != nil {
		return
	}

	t.Log(string(marshal))

	req := &Req{
		RequestId: "abcdefg",
		GrpId:     2,
		CmdId:     1,
		Data:      []byte(string(marshal)),
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(jsonData))

	compressor := gcompress.NewGzipCompressor()
	compress, err := compressor.Compress(jsonData)
	if err != nil {
		return
	}

	toString := base64.StdEncoding.EncodeToString(compress)
	t.Log(toString)
}

func TestWsPing(t *testing.T) {
	// WebSocket 服务器地址
	url := "ws://localhost:8080/ws"
	header := http.Header{}
	header.Set("UUID", "sdfgadf")
	header.Set("PlatformID", "5")

	// 连接到 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// 设置一个定时任务，每隔 5 秒发送一次 ping 消息
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				// 发送 ping 消息
				if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					log.Println("Error sending ping:", err)
					return
				}
				fmt.Println("Sent ping")
			}
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			return
		}

		switch messageType {
		case websocket.PongMessage:
			fmt.Println("Received pong:", string(message))
		case websocket.TextMessage:
			fmt.Println("Received text:", string(message))
		default:
			fmt.Println("Received message:", messageType, string(message))
		}
	}
}

func TestSuperConn(t *testing.T) {
	var wg sync.WaitGroup
	numConnections := 10000          // 1万连接
	maxConcurrentConnections := 1000 // 最大并发连接数

	// 控制并发数，避免一次性启动过多连接
	semaphore := make(chan struct{}, maxConcurrentConnections)

	// 启动 1 万个 WebSocket 连接
	for i := 1; i <= numConnections; i++ {
		semaphore <- struct{}{} // 确保当前并发数不超过最大限制
		wg.Add(1)

		// 使用 goroutine 启动每个 WebSocket 连接
		go func(id int) {
			defer func() { <-semaphore }()
			connectWebSocket(id, &wg)
		}(i)
	}

	// 等待所有连接完成
	wg.Wait()
}
func connectWebSocket(id int, wg *sync.WaitGroup) {
	defer wg.Done()

	// WebSocket 服务器地址
	url := "ws://localhost:8080/stream"
	header := http.Header{}
	header.Set("UUID", fmt.Sprintf("UUID-%d", id)) // 每个连接有不同的 UUID
	header.Set("PlatformID", "5")

	// 连接到 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Printf("Connection %d failed: %s", id, err)
		return
	}
	defer conn.Close()

	// 设置一个定时任务，每隔 5 秒发送一次 ping 消息
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				// 发送 ping 消息
				if err := conn.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					log.Printf("Connection %d: Error sending ping: %s", id, err)
					return
				}
				fmt.Printf("Connection %d: Sent ping\n", id)
			}
		}
	}()

	conn.WriteMessage(1, []byte("{\"request_id\":\"abcdefg\",\"grp_id\":2,\"cmd_id\":1,\"data\":\"eyJtc2ciOiJoZWxsbyJ9\"}"))

	// 读取并处理消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Connection %d: Read error: %s", id, err)
			return
		}

		switch messageType {
		case websocket.PongMessage:
			fmt.Printf("Connection %d: Received pong: %s\n", id, string(message))
		case websocket.TextMessage:
			fmt.Printf("Connection %d: Received text: %s\n", id, string(message))
		default:
			fmt.Printf("Connection %d: Received message: %s\n", id, string(message))
		}
	}
}
