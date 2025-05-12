package grpc_conterller

import (
	"context"
)

func InitGameInfo(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
	//var (
	//	request = pbs.MTInitReq{}
	//	comResp = component.NewNetMessage(int32(pbs.Meb_mtInitAck))
	//	ack     = pbs.MTInitAck{}
	//)
	//
	//if req.MsgId != int32(pbs.Meb_mtInitReq) {
	//	global.GVA_LOG.Error("InitGameInfo 协议号不正确", zap.Any("GetBetList", req))
	//	comResp.AckHead.Code = pbs.Code_ProtocNumberError
	//	return nil, errors.New("协议号不正确")
	//}
	//
	//err := proto.Unmarshal(req.Content, &request)
	//if err != nil {
	//	global.GVA_LOG.Error("InitGameInfo", zap.Error(err))
	//	comResp.AckHead.Code = pbs.Code_DataCompileError
	//	return comResp, err
	//}
	//global.GVA_LOG.Infof(" InitGameInfo：%v", &request)
	//
	////err = logic.SaveMtUser(request.UserId, request.Nickname)
	////if err != nil {
	////	global.GVA_LOG.Error("InitGameInfo", zap.Error(err))
	////}
	//
	//global.GVA_LOG.Infof(" InitGameInfo ack:%v", &ack)
	////返回数据
	//ackMarshal, _ := proto.Marshal(&ack)
	//comResp.Content = ackMarshal
	//return comResp, nil
	return nil, nil
}

//func GetBetList(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.BetReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnBetListAck))
//		ack     = pbs.BetListAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnBetReq) {
//		global.GVA_LOG.Error("GetBetList 协议号不正确", zap.Any("GetBetList", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("GetBetList", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" GetBetList：%v", &request)
//
//	//lists, err := logic.GetBetList()
//	//if err != nil {
//	//	global.GVA_LOG.Error("GetBetList", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//var betListData []*pbs.BetListData
//	//for k, _ := range lists {
//	//	list := lists[k]
//	//	item := pbs.BetListData{
//	//		Id:  int32(list.ID),
//	//		Bet: float32(list.Bet),
//	//	}
//	//	betListData = append(betListData, &item)
//	//}
//
//	//ack.BetListData = betListData
//
//	global.GVA_LOG.Infof(" GetBetList ack:%v", &ack)
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtCurrStatus(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.MTStatusReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnmStatusAck))
//		ack     = pbs.MTStatusAck{
//			PeriodId:  0,
//			State:     0,
//			Layer:     0,
//			StartTime: "",
//			PlayerMeta: &pbs.PlayerMeta{
//				UserId:         "",
//				PeriodId:       0,
//				IsBetAuto:      false,
//				ResidueCoinNum: 0,
//			},
//			LayerMeta: &pbs.LayerMeta{
//				BossDp:           0,
//				AutoSoldiersDp:   0,
//				NoAutoSoldiersDp: 0,
//				State:            0,
//			},
//			GameTurnStatus: 0,
//		}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnStatusReq) {
//		global.GVA_LOG.Error(" 协议号不正确", zap.Any("MtCurrStatus", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtCurrStatus", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtCurrStatus：%v", &request)
//
//	//record, err := logic.GetMtRoomPeriod()
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtCurrStatus", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//
//	//user, err := logic.GetMtUser(request.UserId)
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtCurrStatus GetMtUser", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//
//	//if user.IsAuto == 1 {
//	//	ack.PlayerMeta.IsBetAuto = true
//	//}
//	//
//	//userBetAutoInfo, err := table.GetUserBetAutoByUid(request.UserId)
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtCurrStatus GetUserBetAutoByUid", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//ack.PlayerMeta.ResidueCoinNum = float32(userBetAutoInfo.Bet)
//	//
//	////获取当前层的配置
//	//layerConfig, err := dao.GetLayerConfigByLayer(mt.MTRoomManager.Turn)
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtCurrStatus GetUserBetAutoByUid", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//
//	//state := 0
//	////当前时间和开始时间的差值
//	//subTime := time.Now().Sub(record.StartTime).Seconds()
//	////准备阶段
//	//if subTime > 0 && subTime < float64(layerConfig.AttackCountdownTime) {
//	//	state = 0
//	//} else {
//	//	state = 1 //todo
//	//}
//	//
//	////所有自动玩家小兵
//	//var autoSoldiersDp float64
//	//autoSoldiersDp = dao.GetAutoSoldiersDp()
//	//
//	////所有手动玩家小兵
//	//var noAutoSoldiersDp float64
//	//noAutoSoldiersDp = dao.NoAutoSoldiersDp(record.RoomNo, mt.MTRoomManager.Turn)
//	//
//	//layerMeta := &pbs.LayerMeta{
//	//	BossDp:           float32(layerConfig.BloodVolume),
//	//	AutoSoldiersDp:   float32(autoSoldiersDp),
//	//	NoAutoSoldiersDp: float32(noAutoSoldiersDp),
//	//	State:            int32(state),
//	//	PrepareTime:      mt.MTRoomManager.PrepareTime.Format("2006-01-02 15:04:05"),
//	//}
//	//
//	//ack.LayerMeta = layerMeta
//	//ack.PeriodId = int32(record.PeriodId)
//	//ack.State = int32(record.Status)
//	//ack.StartTime = record.StartTime.Format("2006-01-02 15:04:05")
//	//ack.PlayerMeta.UserId = request.UserId
//	//ack.Layer = int32(mt.MTRoomManager.Turn)
//	//
//	//ack.GameTurnStatus = int32(mt.MTRoomManager.GameStart)
//	//
//	//global.GVA_LOG.Infof(" MtCurrStatus ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtIsAutoUser(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.IsAutoReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnIsAutoAck))
//		ack     = pbs.IsAutoAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnIsAutoReq) {
//		global.GVA_LOG.Error(" 协议号不正确", zap.Any("MtIsAutoUser", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtIsAutoUser", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtIsAutoUser：%v", &request)
//
//	//{
//	//	//如果游戏没有开始 返回
//	//	if !mt.MTRoomManager.IsPeriodStart {
//	//		comResp.AckHead.Code = pbs.Code_GameNotStart
//	//		return comResp, err
//	//	}
//	//	//如果当前期 当前层设置过 返回
//	//	record, err := table.GetUserAutoLogByUId(request.UserId)
//	//	if err != nil {
//	//		global.GVA_LOG.Error("MtIsAutoUser", zap.Error(err))
//	//		comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//		return comResp, err
//	//	}
//	//	if record.ID > 0 {
//	//		if record.PeriodId == mt.MTRoomManager.RoomInfo.PeriodId && record.Layer == mt.MTRoomManager.Turn {
//	//			comResp.AckHead.Code = pbs.Code_AlreadySetIsAuto
//	//			return comResp, err
//	//		}
//	//	}
//	//
//	//}
//	//
//	//logic.UpdateMtUser(request.UserId, int(request.IsAuto))
//
//	global.GVA_LOG.Infof(" MtIsAutoUser ack:%v", &ack)
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtBetAutoNum(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.AutoNumReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnAutoNumAck))
//		ack     = pbs.AutoNumAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnAutoNumReq) {
//		global.GVA_LOG.Error(" 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtBetAutoNum", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof("MtBetAutoNum：%v", &request)
//
//	//if !mt.MTRoomManager.IsPeriodStart {
//	//	comResp.AckHead.Code = pbs.Code_GameNotStart
//	//	return comResp, err
//	//}
//	//
//	//if mt.MTRoomManager.GetStage() == 1 {
//	//	global.GVA_LOG.Infof("MtBetAutoNum 狂暴阶段 不让自动押注")
//	//	comResp.AckHead.Code = pbs.Code_CrazyLimitAutoBet
//	//	return comResp, err
//	//}
//	//
//	//user, err := table.GetMtUser(request.UserId)
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtBetAutoNum", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_DBErr
//	//	return comResp, err
//	//}
//	//if user.IsAuto != 1 {
//	//	comResp.AckHead.Code = pbs.Code_NotAutoUser
//	//	return comResp, err
//	//}
//	//
//	////更新自动池
//	//logic.UpdateUserAutoBet(request.UserId, float64(request.Bet))
//
//	global.GVA_LOG.Infof(" MtBetAutoNum ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtBetNum(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.BetNumReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnBetNumAck))
//		ack     = pbs.BetNumAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnAutoNumReq) {
//		global.GVA_LOG.Error(" 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtBetAutoNum", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtBetAutoNum：%v", &request)
//
//	//if !mt.MTRoomManager.IsPeriodStart {
//	//	comResp.AckHead.Code = pbs.Code_GameNotStart
//	//	return comResp, err
//	//}
//	//
//	////不是狂化阶段 每层之间能押注一次
//	//if mt.MTRoomManager.GetStage() == 0 {
//	//	userBetRecords, err := table.GetUserBetByRoomNoAndUserAndLayer(request.UserId, mt.MTRoomManager.RoomInfo.RoomNo, mt.MTRoomManager.Turn)
//	//	if err != nil {
//	//		global.GVA_LOG.Error("MtBetAutoNum", zap.Error(err))
//	//		comResp.AckHead.Code = pbs.Code_DBErr
//	//		return comResp, err
//	//	}
//	//	if userBetRecords.ID > 0 {
//	//		comResp.AckHead.Code = pbs.Code_HaveCallSoldiers
//	//		return comResp, err
//	//	}
//	//
//	//	//如果开启自动召唤 手动召唤按钮就不可以适应
//	//
//	//}
//	//
//	////添加手动池
//	//layerConfigRecord, err := dao.GetLayerConfigByLayer(mt.MTRoomManager.Turn)
//	//if err != nil {
//	//	global.GVA_LOG.Error("UpdateGameStatus GetLayerConfigByLayer err:%v", zap.Any("err", err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//	//
//	//if request.Bet > int32(layerConfigRecord.CallNum) {
//	//	comResp.AckHead.Code = pbs.Code_ParameterIllegal
//	//	return comResp, err
//	//}
//	//
//	//err = dao.CreateUserBet(request.UserId, uuid.New().String(), float64(request.Bet))
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtCurrStatus GetMtUser", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//
//	global.GVA_LOG.Infof(" MtBetAutoNum ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtOneTouchAddBetNum(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.BetOneTouchAddBetNumReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_mtOneTouchAddBetAck))
//		ack     = pbs.BetOneTouchAddBetNumAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_mtOneTouchAddBetReq) {
//		global.GVA_LOG.Error("MtOneTouchAddBetNum 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtOneTouchAddBetNum", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//
//	global.GVA_LOG.Infof(" MtOneTouchAddBetNum：%v", &request)
//	//
//	//if mt.MTRoomManager.GetStage() != 1 {
//	//	global.GVA_LOG.Infof("MtOneTouchAddBetNum 非狂暴阶段,不能一键召唤")
//	//	comResp.AckHead.Code = pbs.Code_NotCrazyStageBet
//	//	return comResp, err
//	//}
//	//
//	////更新自动池
//	//logic.UpdateUserAutoBet(request.UserId, float64(request.Bet))
//
//	global.GVA_LOG.Infof(" MtOneTouchAddBetNum ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtJoinGameSelectCamp(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.JoinGameSelectCampReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnJoinGameSelectCampAck))
//		ack     = pbs.JoinGameSelectCampAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_mtOneTouchBetReq) {
//		global.GVA_LOG.Error("MtJoinGameSelectCamp 协议号不正确", zap.Any("MtJoinGameSelectCamp", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtJoinGameSelectCamp", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//
//	global.GVA_LOG.Infof(" MtJoinGameSelectCamp：%v", &request)
//
//	global.GVA_LOG.Infof(" MtJoinGameSelectCamp ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtUserPeriodLayerList(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.UserPeriodLayerListReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnUserPeriodLayerListAck))
//		ack     = pbs.UserPeriodLayerListAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnUserPeriodLayerListReq) {
//		global.GVA_LOG.Error("MtUserPeriodLayerList 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("MtUserPeriodLayerList 协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtUserPeriodLayerList", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtUserPeriodLayerList：%v", &request)
//
//	//userPeriodLayerList := dao.GetLayerPassAwardByPeriodIdAndUserId(int(request.PeriodId), request.UserId)
//	//
//	//ack.UserPeriodLayerList = userPeriodLayerList
//
//	global.GVA_LOG.Infof(" MtUserPeriodLayerList ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//func MtPeriodUserRevenueRank(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.UserRevenueRankReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_pnUserRevenueRankAck))
//		ack     = pbs.UserRevenueRankAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_pnUserRevenueRankReq) {
//		global.GVA_LOG.Error("MtPeriodUserRevenueRank 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("MtPeriodUserRevenueRank 协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtPeriodUserRevenueRank", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtPeriodUserRevenueRank：%v", &request)
//
//	//userRevenueRank := dao.PeriodUserRevenueRank(int(request.PeriodId), request.UserId)
//	//
//	//ack.UserRevenueRank = userRevenueRank
//
//	global.GVA_LOG.Infof(" MtPeriodUserRevenueRank ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
//
//// MtOneTouchBet 一键召唤  押注金额
//func MtOneTouchBet(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error) {
//	var (
//		request = pbs.OneTouchBetReq{}
//		comResp = component.NewNetMessage(int32(pbs.Mmb_mtOneTouchBetAck))
//		ack     = pbs.OneTouchBetAck{}
//	)
//
//	if req.MsgId != int32(pbs.Mmb_mtOneTouchBetReq) {
//		global.GVA_LOG.Error("MtOneTouchBet 协议号不正确", zap.Any("MtBetAutoNum", req))
//		comResp.AckHead.Code = pbs.Code_ProtocNumberError
//		return nil, errors.New("协议号不正确")
//	}
//	err := proto.Unmarshal(req.Content, &request)
//	if err != nil {
//		global.GVA_LOG.Error("MtOneTouchBet", zap.Error(err))
//		comResp.AckHead.Code = pbs.Code_DataCompileError
//		return comResp, err
//	}
//	global.GVA_LOG.Infof(" MtOneTouchBet：%v", &request)
//
//	//if !mt.MTRoomManager.IsPeriodStart {
//	//	comResp.AckHead.Code = pbs.Code_GameNotStart
//	//	return comResp, err
//	//}
//	//
//	//if mt.MTRoomManager.GetStage() != 1 {
//	//	global.GVA_LOG.Infof("MtBetAutoNum 狂暴阶段 不让自动押注")
//	//	comResp.AckHead.Code = pbs.Code_CrazyLimitAutoBet
//	//	return comResp, err
//	//}
//	//
//	//userBetAutoInfo, err := table.GetUserBetAutoByUid(request.UserId)
//	//if err != nil {
//	//	global.GVA_LOG.Error(err.Error())
//	//	comResp.AckHead.Code = pbs.Code_DBErr
//	//	return comResp, err
//	//}
//	//
//	//if userBetAutoInfo.Bet < float64(request.Bet) {
//	//	comResp.AckHead.Code = pbs.Code_SCoin
//	//	return comResp, err
//	//}
//	//
//	////添加手动池
//	//err = dao.CreateUserBet(request.UserId, uuid.New().String(), float64(request.Bet))
//	//if err != nil {
//	//	global.GVA_LOG.Error("MtOneTouchBet GetMtUser", zap.Error(err))
//	//	comResp.AckHead.Code = pbs.Code_GetDataFromDbErr
//	//	return comResp, err
//	//}
//
//	global.GVA_LOG.Infof("MtOneTouchBet ack:%v", &ack)
//
//	//返回数据
//	ackMarshal, _ := proto.Marshal(&ack)
//	comResp.Content = ackMarshal
//	return comResp, nil
//}
