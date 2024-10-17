package internal

import (
	"context"
	"crypto/md5"
	"encoding/binary"
	"framework-gin/common/function"
	"framework-gin/ws/constant"
	"framework-gin/ws/proto/pb"
	"github.com/qiafan666/gotato/commons/gcommon"
	"github.com/qiafan666/gotato/commons/glog"
	"math/rand"
	"time"
)

func (ws *WsServer) ChangeOnlineStatus(concurrent int) {
	if concurrent < 1 {
		concurrent = 1
	}
	const renewalTime = constant.OnlineExpire / 3
	renewalTicker := time.NewTicker(renewalTime)

	// 初始化多个通道，用于并发处理 SetUserOnlineStatus 请求
	requestChs := make([]chan *pb.SetUserOnlineStatusReq, concurrent)
	changeStatus := make([][]UserState, concurrent)

	for i := 0; i < concurrent; i++ {
		requestChs[i] = make(chan *pb.SetUserOnlineStatusReq, 64) // 每个通道的缓冲大小为64
		changeStatus[i] = make([]UserState, 0, 100)               // 初始化每个状态数组，容量为100
	}

	mergeTicker := time.NewTicker(time.Second) // 每秒钟合并一次用户状态，确保状态不会积压

	// 将 UserState 转换为 pb.UserOnlineStatus 格式，用于网络传输
	local2pb := func(u UserState) *pb.UserOnlineStatus {
		return &pb.UserOnlineStatus{
			UserID:  u.UserID,
			Online:  u.Online,
			Offline: u.Offline,
		}
	}

	// 随机数用于将用户状态均匀分布到多个通道中，避免单个通道压力过大
	rNum := rand.Uint64()

	// 将用户状态放入对应的处理通道中，确保每个通道的状态数组不会超出容量限制
	pushUserState := func(us ...UserState) {
		for _, u := range us {
			// 使用用户ID的哈希值来决定将状态分配到哪个通道
			sum := md5.Sum([]byte(u.UserID))
			i := (binary.BigEndian.Uint64(sum[:]) + rNum) % uint64(concurrent)
			changeStatus[i] = append(changeStatus[i], u)

			status := changeStatus[i]
			// 当状态数组达到容量时，将其发送到通道进行处理
			if len(status) == cap(status) {
				req := &pb.SetUserOnlineStatusReq{
					Status: gcommon.SliceConvert(status, local2pb),
				}
				changeStatus[i] = status[:0] // 清空已处理的状态数组
				select {
				case requestChs[i] <- req: // 尝试将请求发送到通道
				default:
					// 当通道满时，记录处理过慢的日志
					glog.Slog.DebugKVs(function.WsCtx, "ChangeOnlineStatus user online processing is too slow")
				}
			}
		}
	}

	// 将所有未处理的用户状态推送到请求通道中，确保没有遗漏
	pushAllUserState := func() {
		for i, status := range changeStatus {
			if len(status) == 0 {
				continue
			}
			req := &pb.SetUserOnlineStatusReq{
				Status: gcommon.SliceConvert(status, local2pb),
			}
			changeStatus[i] = status[:0]
			select {
			case requestChs[i] <- req: // 尝试将请求发送到通道
			default:
				// 当通道满时，记录处理过慢的错误日志
				glog.Slog.ErrorKVs(function.WsCtx, "ChangeOnlineStatus user online processing is too slow")
			}
		}
	}

	// 执行 SetUserOnlineStatus 请求
	doRequest := func(req *pb.SetUserOnlineStatusReq) {
		// 为每个请求设置唯一的 requestID 并创建一个带超时的上下文
		ctx, cancel := context.WithTimeout(function.WsCtx, time.Second*5)
		defer cancel()

		for _, status := range req.Status {
			err := ws.rdbOnline.SetUserOnline(ctx, status.UserID, status.Online, status.Offline)
			if err != nil {
				glog.Slog.ErrorKVs(ctx, "ChangeOnlineStatus", "set user online status err", err, "userID", status.UserID, "online", status.Online, "offline", status.Offline)
			}
		}
		glog.Slog.DebugKVs(ctx, "ChangeOnlineStatus", "req", req)
	}

	// 启动多个 goroutine 处理用户在线状态的批量更新请求
	for i := 0; i < concurrent; i++ {
		go func(ch <-chan *pb.SetUserOnlineStatusReq) {
			for req := range ch {
				doRequest(req) // 处理通道中的每个请求
			}
		}(requestChs[i])
	}

	// 主循环，处理用户状态变化和定时任务
	for {
		select {
		case <-mergeTicker.C: // 定时合并用户状态并推送
			pushAllUserState()
		case now := <-renewalTicker.C: // 每次 renewalTicker 触发时，检查需要更新的用户状态
			deadline := now.Add(-constant.OnlineExpire / 3)
			users := ws.clients.GetAllUserStatus(deadline, now) // 获取当前时间段内的用户状态
			glog.Slog.DebugKVs(function.WsCtx, "CheckOnlineStatus renewal ticker", "deadline", deadline, "nowtime", now, "num", len(users), "users", users)
			pushUserState(users...) // 推送用户状态
		case state := <-ws.clients.UserState(): // 当有用户状态变化时
			glog.Slog.DebugKVs(function.WsCtx, "CheckOnlineStatus OnlineCache user online change", "userID", state.UserID, "online", state.Online, "offline", state.Offline)
			pushUserState(state) // 推送单个用户的状态变化
		}
	}
}
