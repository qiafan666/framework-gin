package main

import (
	"encoding/base64"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gcommon"
	"testing"
)

type Req struct {
	RequestID string `json:"request_id"   validate:"required"`
	GrpID     uint8  `json:"grp_id" validate:"required"` // 消息组id
	CmdID     uint8  `json:"cmd_id" validate:"required"` // 消息的ID
	Data      []byte `json:"data"`
}

func TestEncode(t *testing.T) {

	marshal, err := proto.Marshal(&pb.ReqHealth{Msg: "hello"})
	if err != nil {
		return
	}

	req := &Req{
		RequestID: "abcdefg",
		GrpID:     2,
		CmdID:     1,
		Data:      marshal,
	}

	gobEncoder := gcommon.NewGobEncoder()

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
