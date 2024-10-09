package internal

//func (ws *WsServer) ChangeOnlineStatus(concurrent int) {
//	if concurrent < 1 {
//		concurrent = 1
//	}
//	const renewalTime = constant.OnlineExpire / 3
//	//const renewalTime = time.Second * 10
//	renewalTicker := time.NewTicker(renewalTime)
//
//	requestChs := make([]chan *pbuser.SetUserOnlineStatusReq, concurrent)
//	changeStatus := make([][]UserState, concurrent)
//
//	for i := 0; i < concurrent; i++ {
//		requestChs[i] = make(chan *pbuser.SetUserOnlineStatusReq, 64)
//		changeStatus[i] = make([]UserState, 0, 100)
//	}
//
//	mergeTicker := time.NewTicker(time.Second)
//
//	local2pb := func(u UserState) *pbuser.UserOnlineStatus {
//		return &pbuser.UserOnlineStatus{
//			UserID:  u.UserID,
//			Online:  u.Online,
//			Offline: u.Offline,
//		}
//	}
//
//	rNum := rand.Uint64()
//	pushUserState := func(us ...UserState) {
//		for _, u := range us {
//			sum := md5.Sum([]byte(u.UserID))
//			i := (binary.BigEndian.Uint64(sum[:]) + rNum) % uint64(concurrent)
//			changeStatus[i] = append(changeStatus[i], u)
//			status := changeStatus[i]
//			if len(status) == cap(status) {
//				req := &pbuser.SetUserOnlineStatusReq{
//					Status: gcommon.SliceConvert(status, local2pb),
//				}
//				changeStatus[i] = status[:0]
//				select {
//				case requestChs[i] <- req:
//				default:
//					glog.Slog.DebugKVs(context.Background(), "user online processing is too slow", nil)
//				}
//			}
//		}
//	}
//
//	pushAllUserState := func() {
//		for i, status := range changeStatus {
//			if len(status) == 0 {
//				continue
//			}
//			req := &pbuser.SetUserOnlineStatusReq{
//				Status: gcommon.SliceConvert(status, local2pb),
//			}
//			changeStatus[i] = status[:0]
//			select {
//			case requestChs[i] <- req:
//			default:
//				glog.Slog.ErrorKVs(context.Background(), "user online processing is too slow")
//			}
//		}
//	}
//
//	var count atomic.Int64
//	operationIDPrefix := fmt.Sprintf("p_%d_", os.Getpid())
//	doRequest := func(req *pbuser.SetUserOnlineStatusReq) {
//		opIdCtx := mcontext.SetOperationID(context.Background(), operationIDPrefix+strconv.FormatInt(count.Add(1), 10))
//		ctx, cancel := context.WithTimeout(opIdCtx, time.Second*5)
//		defer cancel()
//		if _, err := ws.userClient.Client.SetUserOnlineStatus(ctx, req); err != nil {
//			glog.Slog.ErrorKVs(ctx, "update user online status", "err", err)
//		}
//	}
//
//	for i := 0; i < concurrent; i++ {
//		go func(ch <-chan *pbuser.SetUserOnlineStatusReq) {
//			for req := range ch {
//				doRequest(req)
//			}
//		}(requestChs[i])
//	}
//
//	for {
//		select {
//		case <-mergeTicker.C:
//			pushAllUserState()
//		case now := <-renewalTicker.C:
//			deadline := now.Add(-constant.OnlineExpire / 3)
//			users := ws.clients.GetAllUserStatus(deadline, now)
//			glog.Slog.DebugKVs(context.Background(), "renewal ticker", "deadline", deadline, "nowtime", now, "num", len(users), "users", users)
//			pushUserState(users...)
//		case state := <-ws.clients.UserState():
//			glog.Slog.DebugKVs(context.Background(), "OnlineCache user online change", "userID", state.UserID, "online", state.Online, "offline", state.Offline)
//			pushUserState(state)
//		}
//	}
//}
