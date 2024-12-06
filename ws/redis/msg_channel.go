package redis

import (
	"context"
	"framework-gin/common/function"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/golang/protobuf/proto"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/glog"
	"github.com/redis/go-redis/v9"
)

type ChannelInterface interface {
	MsgPublish(ctx context.Context, message *pb.ReqPushMsgToOther) error
	MsgSubscribe(ctx context.Context, channel string) <-chan *pb.ReqPushMsgToOther
}

func NewMsgChannel(rdb redis.UniversalClient) ChannelInterface {

	return &channelImpl{
		rdb:        rdb,
		msgChannel: constant.MsgChannel,
	}
}

type channelImpl struct {
	rdb        redis.UniversalClient
	msgChannel string
}

// MsgPublish 发布消息
func (m *channelImpl) MsgPublish(ctx context.Context, message *pb.ReqPushMsgToOther) error {
	msg, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return m.rdb.Publish(ctx, m.msgChannel, msg).Err()
}

// MsgSubscribe 订阅消息
func (m *channelImpl) MsgSubscribe(ctx context.Context, channel string) <-chan *pb.ReqPushMsgToOther {
	pubsub := m.rdb.Subscribe(ctx, channel)
	ch := make(chan *pb.ReqPushMsgToOther)

	gcommon.Go(func() {
		defer pubsub.Close()
		for msg := range pubsub.Channel() {
			var pubSubMsg *pb.ReqPushMsgToOther
			if err := proto.Unmarshal([]byte(msg.Payload), pubSubMsg); err != nil {
				glog.Slog.ErrorKVs(function.WsCtx, "msgpack unmarshal error", "err", err)
				continue
			}
			ch <- pubSubMsg
		}
		close(ch)
	})

	return ch
}
