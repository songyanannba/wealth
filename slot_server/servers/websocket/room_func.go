package websocket

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/common"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/lib/src/dao"
	"slot_server/lib/src/logic"
	"slot_server/protoc/pbs"
	"strconv"
	"time"
)

func NatsSendAllUserMsg(trs *RoomSpace, msg *pbs.NetMessage) {
	trs.ComRoomSpace.NatsSendAllUserMsg(msg)
}

func NatsSendAimUserMsg(trs *RoomSpace, msg *pbs.NetMessage, userId string) {
	trs.ComRoomSpace.NatsSendAimUserMsg(msg, userId)
}

func (trs *RoomSpace) InItFunc() {
	//初始房间卡牌
	if len(trs.RoomBaseCard) == 0 {
		trs.RoomBaseCard = SlotRoomManager.RoomBaseCard
	}

	if trs.RoomVersionCard != nil {
		trs.RoomVersionCard = SlotRoomManager.RoomVersionCard
	}

	//初始房间问题
	if len(trs.RoomIssueConfig) == 0 {
		trs.RoomIssueConfig = SlotRoomManager.RoomIssueConfig
	}

	//加入房间
	trs.Register(strconv.Itoa(int(pbs.Meb_joinRoom)), JoinRoom)

	//重新加入房间
	trs.Register(strconv.Itoa(int(pbs.Meb_reJoinRoom)), ReJoinRoom)

	//离开房间
	trs.Register(strconv.Itoa(int(pbs.Meb_leaveRoom)), LeaveRoom)

	//就绪
	trs.Register(strconv.Itoa(int(pbs.Meb_readyMsg)), Ready)

	//取消就绪
	trs.Register(strconv.Itoa(int(pbs.Meb_cancelReady)), CancelReady)

	//被踢
	trs.Register(strconv.Itoa(int(pbs.Meb_kickRoom)), KickRoom)

	//开始游戏
	trs.Register(strconv.Itoa(int(pbs.Meb_startPlay)), StartPlay)

	//确认加载完成
	trs.Register(strconv.Itoa(int(pbs.Meb_loadCompleted)), LoadCompleted)

	//邀请
	trs.Register(strconv.Itoa(int(pbs.Meb_inviteFriend)), InviteFriend)

	//操作牌 Meb_operateCards
	trs.Register(strconv.Itoa(int(pbs.Meb_operateCards)), OperateCards)

	//点赞
	trs.Register(strconv.Itoa(int(pbs.Meb_likeCards)), LikeCards)

}

// LoadCompleted 确认加载完成
func LoadCompleted(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.LoadCompletedReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("LoadCompleted: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("LoadCompleted %v", request)

	userId := request.UserId
	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_loadCompleted), config.SlotServer)

	if !trs.ComRoomSpace.IsStartGame {
		global.GVA_LOG.Error("LoadCompleted 还没有开始游戏")
		netMessageResp.AckHead.Code = pbs.Code(common.GameNotStart)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//获取房间用户信息
	userInfo, err := trs.ComRoomSpace.GetUserInfo(userId)
	if err != nil {
		global.GVA_LOG.Error("LoadCompleted GetUserInfo  ", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	if trs.RoomInfo.UserNumLimit != len(trs.ComRoomSpace.UserInfos) {
		global.GVA_LOG.Error("LoadCompleted 房间人数不对")
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//确认者数据
	trs.AddLoadComps(userId, userInfo)

	//是否存在问题
	currIssue := &models.Issue{}
	currIssue, err = trs.ComRoomSpace.GetSelectIssue()
	if err != nil {
		global.GVA_LOG.Infof("LoadCompleted 房间{%v},本轮{%v} Issue :%v", trs.RoomInfo.RoomNo, trs.ComRoomSpace.GetTurn(), currIssue)
	} else {
		currIssue = currIssue
	}

	//已经加载完成 断线重连
	if trs.IsAllLoadComps == true {
		trs.ReLoadCompleted(userId, currIssue)
	} else {
		trs.LoadCompletedFirst(userId, currIssue)
	}

	//返回
	return nil, nil
}

// LikeCards 给牌点赞
func LikeCards(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.LikeCardReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("LikeCards: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("LikeCards %v", request)

	var (
		userId     = request.UserId     //点赞的用户
		likeUserId = request.LikeUserId //被点赞的用户ID
		reqCards   = request.Card       //请求中
	)

	//获取房间用户信息
	_, err = trs.ComRoomSpace.GetUserInfo(userId)
	if err != nil {
		netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_likeCards), config.SlotServer)
		global.GVA_LOG.Error("LikeCards GetUserInfo  ", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//那个用户没点赞
	likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(request.UserId)
	if len(likeUserInfo) > 0 {
		//该用户已经给别人点过赞
		global.GVA_LOG.Infof("LikeCards 该用户已经给别人点过赞 userID %v", request.UserId)
		netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_likeCards), config.SlotServer)
		netMessageResp.AckHead.Code = pbs.Code(common.HaveLikeCard)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, err
	}

	//获取用户已经出过的牌
	outCards := trs.ComRoomSpace.GetUserOutEdCards(request.LikeUserId)
	if len(outCards) == 0 {
		global.GVA_LOG.Error("LikeCards GetCurrCard", zap.Error(err))
		//netMessageResp := helper.NewNetMessage(int32(helper.GetIntUserId(userId)), 0, int32(pbs.Meb_likeCards), config.SlotServer)
		//netMessageResp.AckHead.Code = pbs.Code(common.NotData)
		//NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, err
	}

	likeCard := models.LikeCard{}

	//每轮每次只出一张牌
	isOutCard := false
	for _, outCard := range outCards {
		for _, rc := range reqCards {
			if outCard.CardId == int(rc.CardId) {
				likeCard = models.LikeCard{
					CardId:     outCard.CardId,
					LikeUserId: likeUserId,
					Level:      outCard.Level,
					AddRate:    outCard.AddRate,
				}
				isOutCard = true
			}
		}
	}
	if !isOutCard {
		global.GVA_LOG.Error("LikeCards GetCurrCard 点赞的牌不存在")
		netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_likeCards), config.SlotServer)
		netMessageResp.AckHead.Code = pbs.Code(common.NotData)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, err
	}

	trs.DoLikeCard(request.UserId, request.LikeUserId, likeCard, outCards)

	return nil, nil
}

// OperateCards 操作牌
func OperateCards(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.OperateCardReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("OperateCards: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("OperateCards %v", request)

	var (
		userId   = request.UserId
		opeType  = request.OpeType
		reqCards = request.Card //请求中 要出的牌
		cards    []*models.Card //用户手中当前牌
	)

	//获取房间用户信息
	userInfo, err := trs.ComRoomSpace.GetUserInfo(userId)
	if err != nil {
		netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_operateCards), config.SlotServer)
		global.GVA_LOG.Error("OperateCards GetUserInfo  ", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//获取用户当前的牌
	currCards, err := trs.ComRoomSpace.GetCurrCard(userId)
	if err != nil {
		global.GVA_LOG.Error("OperateCards GetCurrCard", zap.Error(err))
		netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_operateCards), config.SlotServer)
		netMessageResp.AckHead.Code = pbs.Code(common.NotData)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, err
	}
	cards = currCards

	//0:看牌 1:出牌 2:表情 3:重随
	switch opeType {
	case int32(common.LookCards):
		//netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_LookCards), config.SlotServer)
		//protoNum := strconv.Itoa(int(pbs.Meb_LookCards))
		//msgData := models.OperateCardsMsg{
		//	ProtoNum: protoNum,
		//	UserId:   userId,
		//	CardNum:  len(cards),
		//	//Card:     cards,
		//}
		//
		////给客户消息
		//global.GVA_LOG.Infof("OperateCards 看牌的广播: %v", msgData)
		//
		//responseHeadByte, _ := json.Marshal(msgData)
		//netMessageResp.Content = responseHeadByte
		//NatsSendAllUserMsg(trs, netMessageResp) //OperateCards
		return nil, nil
	case int32(common.OpeEmoji):
		netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_opeEmoji), config.SlotServer)
		protoNum := strconv.Itoa(int(pbs.Meb_opeEmoji))
		msgData := models.OperateCardsMsg{
			ProtoNum: protoNum,
			UserId:   userId,
			CardNum:  len(cards),
			EmojiId:  request.EmojiId,
		}

		//给客户消息
		global.GVA_LOG.Infof("OperateCards 表情: %v", msgData)

		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp.Content = responseHeadByte
		NatsSendAllUserMsg(trs, netMessageResp) //OperateCards
		return nil, nil
	case int32(common.ReMakeCards): //随牌
		//netMessageResp := helper.NewNetMessage(int32(helper.GetIntUserId(userId)), 0, int32(pbs.Meb_reMakeCards), config.SlotServer)
		////判断是否在重随 时间段
		//gameStatus, turnTime := trs.CurrGameTurnStateAndDownTime()
		//global.GVA_LOG.Infof("CurrGameTurnStateAndDownTime 随牌 OperateCards gameStatus %v, turnTime %v", gameStatus, turnTime)
		//if gameStatus != int(CliRemakeCard) {
		//	netMessageResp.AckHead.Code = pbs.Code(common.OutReMakeCardTime)
		//	NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		//	return nil, nil
		//}
		//
		////第几次从随 金币是否够
		//extractCard := trs.ComRoomSpace.GetExtractCard(userInfo.UserID)
		//if len(extractCard) >= 4 {
		//	//重随过 需要扣积分
		//	extractCardNum := helper.CeilDiv(len(extractCard), 4)
		//	if extractCardNum > 0 {
		//		coinConsumeConfig := dao.GetCoinConsumeConfigByType(int(table.GetCoinConsume(extractCardNum)))
		//		//判断金币是否够
		//		experienceInfo := dao.GetUserCoinExperience(request.UserId)
		//		if experienceInfo.CoinNum <= coinConsumeConfig.CoinNum {
		//			netMessageResp.AckHead.Code = pbs.Code_KingCoinNotEnough
		//			NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		//			return nil, nil
		//		}
		//		//扣金币
		//		err := dao.UpdateUserCoinNumOrExperience(request.UserId, -coinConsumeConfig.CoinNum, 0, 2)
		//		if err != nil {
		//			global.GVA_LOG.Error("UnpackCardController ,UpdateUserCoinNumOrExperience err", zap.Any("err", err))
		//		}
		//	}
		//}

		//判断是否有可被选择的牌
		//extractCards, err := trs.ComRoomSpace.GetNotExtractCard(userInfo.UserID)
		//if err != nil || len(extractCards) == 0 {
		//	global.GVA_LOG.Error("OperateCards GetNotExtractCard", zap.Error(err))
		//	netMessageResp.AckHead.Code = pbs.Code(common.NotExtractCard)
		//	NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		//	return nil, nil
		//}

		//重置手中的牌

		//所有牌随机打乱
		//helper.SliceShuffle(extractCards)
		//var (
		//	newCards = make([]*table.MbCardConfig, 0) //未被选的牌
		//	outCards = make([]*table.MbCardConfig, 0) //给用户要发的牌
		//	perCards = make([]*models.Card, 0)        //每个人获取4张牌
		//)
		//
		//for _, val := range extractCards {
		//	if len(perCards) < 4 {
		//		v := &models.Card{
		//			CardId:  val.ID,
		//			Name:    val.Name,
		//			Suffix:  val.SuffixName,
		//			Level:   val.Level,
		//			AddRate: val.AddRate,
		//			UserID:  userId}
		//		perCards = append(perCards, v)
		//		outCards = append(outCards, val)
		//	} else {
		//		newCards = append(newCards, val)
		//	}
		//}

		//把手里的牌放到本轮随过的数据结构里面
		//用户当前的牌
		//用户重随 这个要跟着重置
		//err = trs.ComRoomSpace.AddCurrCard(userInfo.UserID, perCards)
		//if err != nil {
		//	global.GVA_LOG.Error("dealCards  AddCurrCard", zap.Error(err))
		//}
		//
		////重置 未抽过的牌
		//trs.ComRoomSpace.ReMakeExtractCard(userInfo.UserID, newCards)
		//
		////把手里的牌放到本轮随过的数据结构里面
		////保留抽过的牌
		//trs.ComRoomSpace.SaveExtractCard(userInfo.UserID, outCards)
		//
		////发牌:给每一个用户发对应的牌
		//msgData := models.DealCardsMsg{
		//	ProtoNum:  strconv.Itoa(int(pbs.Meb_reMakeCards)),
		//	Timestamp: time.Now().Unix(),
		//	UserId:    userInfo.UserID,
		//	RoomNo:    trs.RoomInfo.RoomNo,
		//	Turn:      trs.ComRoomSpace.GetTurn(),
		//	Cards:     perCards,
		//}
		//
		////给用户发送消息
		//global.GVA_LOG.Infof("发牌的广播: %v", msgData)
		//userStateRespMarshal, _ := json.Marshal(msgData)
		//netMessageResp.Content = userStateRespMarshal
		//NatsSendAimUserMsg(trs, netMessageResp, userInfo.UserID)

		//随完牌纪录一下具体数据 （mysql 暂时先不纪录）
		//perCardsMarshal, _ := json.Marshal(perCards)
		//dao.AddTurnDetails(trs.RoomInfo.RoomNo, userInfo.LikeUserId, userInfo.Nickname, trs.ComRoomSpace.GetTurn(), string(perCardsMarshal), "{}")
	case int32(common.OperateCards):
		netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_operateCards), config.SlotServer)
		global.GVA_LOG.Infof("操作牌 用户信息 OperateCards:%v", userInfo)

		if len(cards) != 4 {
			netMessageResp.AckHead.Code = pbs.Code(common.IsOutCards)
			NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
			return nil, nil
		}

		// 1:出牌
		if len(cards) < len(reqCards) || len(reqCards) != 1 {
			netMessageResp.AckHead.Code = pbs.Code(common.ErrorOperateCardsNum)
			NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
			return nil, nil
		}

		paramCards := []*models.Card{}
		for _, card := range request.Card {
			v := models.Card{
				CardId: int(card.CardId),
			}
			paramCards = append(paramCards, &v)
		}

		//出牌逻辑
		_, resCode := trs.OutCart(paramCards, cards, userId) //端侧请求出牌
		if resCode != common.OK {
			netMessageResp.AckHead.Code = pbs.Code(resCode)
			NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		}
		return nil, nil
	default:
		global.GVA_LOG.Error("OperateCards:", zap.Error(err), zap.Any("request", request))
	}

	return nil, nil
}

// KickRoom 被踢 离开房间
func KickRoom(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.KickRoomReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("KickRoom: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("KickRoom %v", request)

	netMessageResp := helper.NewNetMessage(request.OwnerId, "", int32(pbs.Meb_kickRoom), config.SlotServer)

	if trs.ComRoomSpace.UserOwner.UserID != request.OwnerId {
		global.GVA_LOG.Error("KickRoom 不是房主 不能踢人 ", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.NotRoomOwnKick)
		NatsSendAimUserMsg(trs, netMessageResp, request.OwnerId)
		return
	}

	//获取房间被踢用户信息
	_, err = trs.ComRoomSpace.GetUserInfo(request.UserId)
	if err != nil {
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.OwnerId)
		return
	}

	//发送广播 被踢 离开房间
	msgData := models.KickRoomMsg{
		ProtoNum: strconv.Itoa(int(pbs.Meb_kickRoom)),
		UserId:   request.UserId,
		RoomNo:   request.RoomNo,
	}
	//给客户消息
	global.GVA_LOG.Infof("KickRoom 被踢的广播: %v", msgData)

	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAllUserMsg(trs, netMessageResp) //reJoinRoom

	//从房间里面把用户删除
	trs.ComRoomSpace.DelUserInfoAndUserClient(request.UserId)

	//从数据库里面删除
	err = table.DelRoomUsers(request.RoomNo, request.UserId)
	if err != nil {
		global.GVA_LOG.Error("KickRoom ", zap.Error(err))
	}

	//payPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", trs.MemeRoomConfig.Bet+trs.MemeRoomConfig.AdmissionPrice), 64)
	//logic.ReturnUserCoin(request.LikeUserId, trs.RoomInfo.RoomNo, payPrice, "用户被踢离开房间")

	//如果游戏没开始 房间满员状态需要改变
	if trs.ComRoomSpace.IsStartGame != true && trs.RoomInfo.IsOpen == table.RoomStatusFill {
		trs.RoomInfo.IsOpen = table.RoomStatusOpen
		//更新数据库 房间状态
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("KickRoom", zap.Error(err))
		}
	}

	//更新用户维度信息
	updateMap := dao.MakeUpdateData("room_no", "")
	updateMap["is_leave"] = 1
	updateMap["is_owner"] = 0
	updateMap["seat"] = 0
	dao.UpdateUsersRoomRoomNo(request.UserId, updateMap)
	return nil, nil
}

func InviteFriend(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.InviteFriendReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("InviteFriend: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("InviteFriend %v", request)

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_inviteFriend), config.SlotServer)

	if trs.ComRoomSpace.UserOwner.UserID != request.OwnerId {
		global.GVA_LOG.Error("InviteFriend 不是房主 不能邀请", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.NotRoomInviteFriend)
		NatsSendAimUserMsg(trs, netMessageResp, request.OwnerId)
		return
	}

	//todo 检查是否是好友

	//发送广播 给被邀请人
	msgData := models.InviteFriendMsg{
		ProtoNum: strconv.Itoa(int(pbs.Meb_inviteFriend)),
		UserId:   request.InviteUserId,
		RoomNo:   request.RoomNo,
		OwnerInfo: models.MemeRoomUser{
			UserID:   request.OwnerId,
			Nickname: "待定", //todo
			IsOwner:  true,
			Seat:     1,
		},
	}
	//给客户消息
	global.GVA_LOG.Infof("InviteFriend 给被邀请人: %v", msgData)

	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte

	//返回的用户id
	netMessageResp.AckHead.Uid = request.InviteUserId

	global.GVA_LOG.Infof("InviteFriend LikeUserId:{%v} 给客户端发消息:{%v}", request.InviteUserId, msgData)

	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
	NastManager.Producer(netMessageRespMarshal)

	return nil, nil
}

// StartPlay 开始游戏
func StartPlay(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.StartPlayReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("StartPlay: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("StartPlay %v", request)

	userId := request.UserId

	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_startPlay), config.SlotServer)

	userInfo, err := trs.ComRoomSpace.GetUserInfo(userId)
	if err != nil {
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return
	}

	if len(trs.ComRoomSpace.UserInfos) != trs.RoomInfo.UserNumLimit {
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return
	}

	if !userInfo.UserProperty.IsOwner {
		netMessageResp.AckHead.Code = pbs.Code(common.NotRoomOwner)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, err
	}

	//开始标识
	trs.ComRoomSpace.IsStartGame = true
	trs.RoomInfo.IsOpen = table.RoomStatusIng

	//游戏每 小轮状态 游戏开始
	trs.ComRoomSpace.ChangeGameState(EnGameStartExec)

	//var (
	//	userInfos         = trs.ComRoomSpace.UserInfos
	//	memeRoomUserLists []models.MemeRoomUser
	//	index             int
	//)

	//err = table.SaveMemeRoom(trs.RoomInfo)
	//if err != nil {
	//	global.GVA_LOG.Error("StartPlay SaveTavernRoom ", zap.Error(err))
	//}
	//
	////先获取房间全部的用户
	//for k, _ := range userInfos {
	//	userItem := userInfos[k]
	//	userItem.UserProperty.Turn = trs.ComRoomSpace.GetTurn()
	//	tavernRoomUser := models.MemeRoomUser{
	//		LikeUserId:       userItem.LikeUserId,
	//		Nickname:     userItem.Nickname,
	//		Turn:         userItem.UserProperty.Turn,
	//		IsLeave:      userItem.UserProperty.IsLeave,
	//		IsOwner:      userItem.UserProperty.IsOwner,
	//		IsReady:      userItem.UserProperty.IsReady,
	//		Seat:         userItem.UserProperty.Seat,
	//		UserLimitNum: userItem.UserProperty.UserLimitNum,
	//		WinPrice:     userItem.UserProperty.WinPrice,
	//		Bet:          userItem.UserProperty.Bet,
	//	}
	//	memeRoomUserLists = append(memeRoomUserLists, tavernRoomUser)
	//	index++
	//}
	//
	////发送广播
	//msgData := models.StartPlayMsg{
	//	ProtoNum:         strconv.Itoa(int(pbs.Meb_startPlay)),
	//	RoomNo:           request.RoomNo,
	//	Timestamp:        time.Now().Unix(),
	//	MemeRoomUserList: memeRoomUserLists,
	//}
	//
	////给用户消息
	//global.GVA_LOG.Infof("StartPlay 开始游戏的广播: %v", msgData)
	//responseHeadByte, _ := json.Marshal(msgData)
	//netMessageResp.Content = responseHeadByte
	//NatsSendAllUserMsg(trs, netMessageResp)

	return nil, nil
}

// JoinRoom 加入房间
func JoinRoom(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	var (
		roomUserLists = make([]models.MemeRoomUser, 0)
		request       = &pbs.JoinRoomReq{}
		userId        = ""
		roomNo        = ""
	)

	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("JoinRoom: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("JoinRoom %v", request)

	userId = request.UserId
	roomNo = request.RoomNo
	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_joinRoom), config.SlotServer)

	//房间满员
	if int(trs.RoomInfo.IsOpen) == table.RoomStatusFill || trs.RoomInfo.IsOpen != table.RoomStatusOpen {
		global.GVA_LOG.Info("MebJoinRoom JoinRoom")
		netMessageResp.AckHead.Code = pbs.Code(common.JoinRoomFull)
		NatsSendAimUserMsg(trs, netMessageResp, userId)
		return
	}

	if trs.ComRoomSpace.IsStartGame == true && trs.RoomInfo.IsOpen == table.RoomStatusIng {
		//已经开始的对局
		global.GVA_LOG.Info("MebJoinRoom JoinRoom 已经开始的对局")
		netMessageResp.AckHead.Code = pbs.Code(common.JoinRoomErr)
		NatsSendAimUserMsg(trs, netMessageResp, userId)
		return
	}

	//匹配类型的房间 只能邀请一个好友
	if trs.RoomInfo.RoomType == 2 && len(trs.ComRoomSpace.UserInfos) == 2 {
		global.GVA_LOG.Info("匹配类型的房间 只能邀请一个好友")
		netMessageResp.AckHead.Code = pbs.Code(common.JoinRoomFull)
		NatsSendAimUserMsg(trs, netMessageResp, userId)
		return
	}

	userInfo, _ := trs.ComRoomSpace.GetUserInfo(userId)
	if userInfo == nil {
		//加入房间
		//todo 用户昵称
		//保存用户信息
		user := models.NewUserInfo(userId, "", models.NewUserProperty(0, 0, false, trs.MemeRoomConfig.Bet), models.UserExt{
			RoomNo: request.RoomNo,
		})
		userInfo = &user

		global.GVA_LOG.Infof("JoinRoom 加入房间:%v", userInfo)

		//添加用户到 房间用户里面
		trs.ComRoomSpace.AddUserInfos(userId, userInfo) //JoinRoom
	}

	//继续游戏的房间 需要设置房主
	if trs.RoomInfo.IsGoOn == 1 {
		err := GoOnJoinRoom(trs, userId, roomNo, userInfo)
		if err != nil {
			global.GVA_LOG.Error("JoinRoom GoOnJoinRoom ", zap.Error(err))
			return nil, err
		}
	}

	if len(trs.ComRoomSpace.UserInfos) == trs.RoomInfo.UserNumLimit {
		trs.RoomInfo.IsOpen = table.RoomStatusFill
		//更新数据库 房间状态
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("JoinRoom ", zap.Error(err))
			netMessageResp.AckHead.Code = pbs.Code(common.ModelAddError)
			NatsSendAimUserMsg(trs, netMessageResp, userId)
			return nil, err
		}
	}

	//查看房间里面是否有用户信息
	roomUsersByRoomNoAndUid, err := table.RoomUsersByRoomNoAndUid(trs.RoomInfo.RoomNo, userId)
	if err != nil {
		global.GVA_LOG.Error("JoinRoom RoomUsersByRoomNoAndUid", zap.Error(err))
		netMessageResp.AckHead.Code = pbs.Code(common.ModelAddError)
		NatsSendAimUserMsg(trs, netMessageResp, userId)
		return nil, err
	}

	//已经加入过房间 重回游戏
	if roomUsersByRoomNoAndUid.ID > 0 && roomUsersByRoomNoAndUid.UserId == userId {
		//查看加入房间的信息
		roomUserLists, _ = dao.GetRoomUser(userInfo.UserExt.RoomNo, trs.ComRoomSpace.GetTurn())
		msgData := models.JoinRoomMsg{
			ProtoNum:  strconv.Itoa(int(pbs.Meb_joinRoom)),
			Timestamp: time.Now().Unix(),
			RoomCom: models.NewRoomCom(trs.RoomInfo.RoomNo,
				userId,
				trs.RoomInfo.Name,
				trs.RoomInfo.ID, 0,
				trs.RoomInfo.UserNumLimit,
				int(trs.RoomInfo.RoomType),
				int(trs.RoomInfo.RoomLevel),
				trs.RoomInfo.IsOpen),
			RoomUserList: roomUserLists,
		}

		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp.Content = responseHeadByte

		global.GVA_LOG.Infof("JoinRoom  加入房间的广播: %v", msgData)
		NatsSendAllUserMsg(trs, netMessageResp) //JoinRoom
	} else {
		seat := len(trs.ComRoomSpace.UserInfos)
		//获取座位
		if seat != 1 {
			seat, err = logic.GetSeat(trs.RoomInfo.RoomNo, seat)
			if err != nil {
				global.GVA_LOG.Error("JoinRoom GetSeat", zap.Error(err))
			}
			if seat != 0 {
				seat = seat
			}
		}

		//添加房间用户
		isOwner := 0
		if userInfo.UserProperty.IsOwner {
			isOwner = 1
		}

		tavernRoomUsers := table.NewRoomUsers(userId, userInfo.UserExt.RoomNo, userInfo.Nickname, seat, 1, isOwner, trs.MemeRoomConfig.Bet, 0)
		err := table.CreateRoomUsers(tavernRoomUsers)
		if err != nil {
			global.GVA_LOG.Error("JoinRoom CreateRoomUsers", zap.Error(err))
			netMessageResp.AckHead.Code = pbs.Code(common.ModelAddError)
			NatsSendAimUserMsg(trs, netMessageResp, userId)
			return nil, err
		}

		//修改用户维度的用户信息
		tavernUserRoomData := table.NewUserRoom(tavernRoomUsers.UserId, tavernRoomUsers.RoomNo, tavernRoomUsers.Nickname, tavernRoomUsers.IsLeave, tavernRoomUsers.IsKilled,
			tavernRoomUsers.IsOwner, tavernRoomUsers.Turn, tavernRoomUsers.Seat, tavernRoomUsers.IsRobot, tavernRoomUsers.IsReady)
		err = dao.CreateOrUpdateUsersRoom(tavernUserRoomData)
		if err != nil {
			global.GVA_LOG.Error("JoinRoom", zap.Any("CreateOrUpdateUsersRoom", err))
		}

		//设置座位
		if err := trs.ComRoomSpace.SetUserSeat(userId, seat); err != nil {
			global.GVA_LOG.Error("JoinRoom SetUserSeat", zap.Error(err))
			netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
			NatsSendAimUserMsg(trs, netMessageResp, userId)
			return nil, err
		}

		//查看加入房间的信息
		roomUserLists, _ = dao.GetRoomUser(userInfo.UserExt.RoomNo, 0)

		//发送广播 谁加入房间
		msgData := models.JoinRoomMsg{
			ProtoNum:  strconv.Itoa(int(pbs.Meb_joinRoom)),
			Timestamp: time.Now().Unix(),
			RoomCom: models.NewRoomCom(trs.RoomInfo.RoomNo,
				userId, trs.RoomInfo.Name,
				trs.RoomInfo.ID,
				0,
				trs.RoomInfo.UserNumLimit,
				int(trs.RoomInfo.RoomType),
				int(trs.RoomInfo.RoomLevel),
				trs.RoomInfo.IsOpen,
			),
			RoomUserList: roomUserLists,
		}

		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp.Content = responseHeadByte
		global.GVA_LOG.Infof("JoinRoom  加入房间的广播: %v", msgData)
		NatsSendAllUserMsg(trs, netMessageResp) //JoinRoom
	}

	return nil, nil
}

func GoOnJoinRoom(trs *RoomSpace, userId, roomNo string, userInfo *models.UserInfo) error {
	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_joinRoom), config.SlotServer)
	//如果原来的房主 回到房间 还是房主
	if userId == trs.RoomInfo.UserId {
		trs.RoomInfo.Owner = userId
		//更新数据库 房主
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("JoinRoom SaveMemeRoom", zap.Error(err))
			netMessageResp.AckHead.Code = pbs.Code(common.ModelAddError)
			NatsSendAimUserMsg(trs, netMessageResp, userId)
			return err
		}

		userInfo.SetUserIsOwner(true)
		trs.ComRoomSpace.UserOwner = userInfo

		//把其他用户的房主信息修改过来
		for _, val := range trs.ComRoomSpace.UserInfos {
			if val.UserID != userInfo.UserID && val.UserProperty.IsOwner == true {
				val.SetUserIsOwner(false)
				record, _ := table.GetRoomUser(val.UserID, roomNo)
				//更新新房主标识
				if record.ID > 0 {
					err = dao.UpdateRoomOwner(val.UserID, roomNo, table.NotBeOwner)
					if err != nil {
						global.GVA_LOG.Error("JoinRoom LeaveRoom UpdateTavernRoomUsers", zap.Error(err))
					}
				}
				//更新用户维度状态 更新新房主用户状态
				updateMap := dao.MakeUpdateData("room_no", roomNo)
				updateMap["is_owner"] = 0
				dao.UpdateUsersRoomRoomNo(val.UserID, updateMap)
			}
		}
	}

	//房主是否已经提前加入
	isBeForJoin := false
	for _, val := range trs.ComRoomSpace.UserInfos {
		if val.UserID == trs.RoomInfo.UserId {
			isBeForJoin = true
		}
	}

	//第一个回到房间的是房主 ；后面如果房主回到房间，需要把房主还给以前的房主
	if !isBeForJoin && len(trs.ComRoomSpace.UserInfos) == 1 {
		if trs.RoomInfo.UserId != userId {
			trs.RoomInfo.Owner = userId
			//更新数据库 房间状态
			err := table.SaveMemeRoom(trs.RoomInfo)
			if err != nil {
				global.GVA_LOG.Error("JoinRoom SaveMemeRoom", zap.Error(err))
				netMessageResp.AckHead.Code = pbs.Code(common.ModelAddError)
				NatsSendAimUserMsg(trs, netMessageResp, userId)
				return err
			}

			userInfo.SetUserIsOwner(true)
			trs.ComRoomSpace.UserOwner = userInfo
		}
	}

	return nil
}

// ReJoinRoom 重新回到已经开始的游戏对局
func ReJoinRoom(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.JoinRoomReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("ReJoinRoom: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("ReJoinRoom %v", request)
	userId := request.UserId

	//保存用户信息
	userInfo := &models.UserInfo{}
	info, _ := trs.ComRoomSpace.GetUserInfo(userId)
	if info != nil {
		userInfo = info
	} else {
		global.GVA_LOG.Infof("ReJoinRoom: %v", userInfo)
		return nil, nil
	}

	//添加用户到 房间用户里面
	trs.ComRoomSpace.UpdateUserInfoAndUserClient(userId, userInfo) //reJoinRoom
	if len(trs.ComRoomSpace.UserInfos) == trs.RoomInfo.UserNumLimit && trs.RoomInfo.IsOpen == table.RoomStatusOpen {
		//更新数据库 房间状态
		trs.RoomInfo.IsOpen = table.RoomStatusFill
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("ReJoinRoom ", zap.Error(err))
			return nil, nil
		}
	}

	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_reJoinRoom), config.SlotServer)

	msgData := models.JoinRoomMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_reJoinRoom)),
		Timestamp: time.Now().Unix(),
		RoomCom: models.NewRoomCom(trs.RoomInfo.RoomNo,
			userId, trs.RoomInfo.Name,
			trs.RoomInfo.ID,
			0,
			trs.RoomInfo.UserNumLimit,
			int(trs.RoomInfo.RoomType),
			int(trs.RoomInfo.RoomLevel),
			trs.RoomInfo.IsOpen),
	}

	//已经加入过房间 重回游戏
	roomUserLists, _ := trs.ComRoomSpace.UserInfoToRoomUser()
	msgData.RoomUserList = roomUserLists

	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte

	global.GVA_LOG.Infof("ReJoinRoom 重回游戏 加入房间的广播: %v", msgData)
	NatsSendAllUserMsg(trs, netMessageResp) //reJoinRoom

	return nil, nil
}

// Ready  就绪标识
func Ready(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.ReadyReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("Ready: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("Ready %v", request)
	userId := request.UserId

	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_readyMsg), config.SlotServer)

	//获取房间用户信息
	userInfo, err := trs.ComRoomSpace.GetUserInfo(request.UserId)
	if err != nil {
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return
	}

	//查看是否就绪
	if userInfo.GetUserIsReady() == int(models.Ready) {
		//已经就绪 直接返回
		return nil, nil
	}

	//就绪标识
	userInfo.SetUserIsReady(int(models.Ready))

	//获取房间人数
	//发送广播 谁 就绪
	msgData := models.ReadyMsg{
		ProtoNum: strconv.Itoa(int(pbs.Meb_readyMsg)),
		UserId:   userId,
		RoomNo:   request.RoomNo,
	}
	global.GVA_LOG.Infof("Ready 就绪的广播: %v", msgData)
	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAllUserMsg(trs, netMessageResp) //JoinRoom

	//修改数据库了
	err = dao.UpdateRoomUserReady(userId, request.RoomNo, int8(models.Ready))
	if err != nil {
		global.GVA_LOG.Error("Ready UpdateTavernRoomUserReady", zap.Error(err))
	}
	return nil, nil
}

// CancelReady 取消就绪
func CancelReady(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.CancelReadyReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("CancelReady: %v %v", zap.Error(err))
		return nil, err
	}
	global.GVA_LOG.Infof("CancelReady %v", request)
	userId := request.UserId

	netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_cancelReadyMsg), config.SlotServer)

	//获取房间用户信息
	userInfo, err := trs.ComRoomSpace.GetUserInfo(request.UserId)
	if err != nil {
		netMessageResp.AckHead.Code = pbs.Code(common.UserNotInRoom)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		global.GVA_LOG.Error("CancelReady roomSpaceInfo.ComRoomSpace.GetUserInfo  ", zap.Error(err))
		return nil, nil
	}

	//查看是否就绪
	if userInfo.GetUserIsReady() != int(models.Ready) {
		//本来就没有就绪
		netMessageResp.AckHead.Code = pbs.Code(common.NotReadyNotCancel)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//如果已经开始游戏 不让取消
	if trs.ComRoomSpace.IsStartGame {
		netMessageResp.AckHead.Code = pbs.Code(common.RoomGameStatusIng)
		NatsSendAimUserMsg(trs, netMessageResp, request.UserId)
		return nil, nil
	}

	//取消 就绪标识
	userInfo.SetUserIsReady(int(models.NotReady))

	//获取房间人数
	//发送广播 谁取消就绪
	msgData := models.ReadyMsg{
		ProtoNum: strconv.Itoa(int(pbs.Meb_cancelReadyMsg)),
		UserId:   userId,
		RoomNo:   request.RoomNo,
	}

	global.GVA_LOG.Infof("CancelReady 取消 就绪的广播: %v", msgData)
	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAllUserMsg(trs, netMessageResp) //CancelReady

	//修改数据库了
	err = dao.UpdateRoomUserReady(userId, request.RoomNo, int8(models.NotReady))
	if err != nil {
		global.GVA_LOG.Error("CancelReady UpdateTavernRoomUserReady", zap.Error(err))
	}

	return nil, nil
}

func LeaveRoom(message []byte, trs *RoomSpace) (resMessage []byte, err error) {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	request := &pbs.LeaveRoomReq{}
	if err := proto.Unmarshal(message, request); err != nil {
		global.GVA_LOG.Error("LeaveRoomController:", zap.Error(err))
		return nil, err
	}

	global.GVA_LOG.Infof("LeaveRoomController %v", request)
	userId := request.UserId

	//获取房间用户信息
	userInfo, err := trs.ComRoomSpace.GetUserInfo(userId)
	if err != nil {
		global.GVA_LOG.Error("LeaveRoomController  ComRoomSpace GetUserInfo  ", zap.Error(err))
		//return common.UserNotInRoom
	}
	global.GVA_LOG.Infof("LeaveRoom 离开房间:%v", userInfo)

	if trs.ComRoomSpace.IsStartGame != true {
		leaveRoomNotStartGameCode := trs.LeaveRoomNotStartGame(userInfo)

		if leaveRoomNotStartGameCode != common.OK {
			//发送广播 离开房间
			msgData := models.LeaveRoomMsg{
				ProtoNum: strconv.Itoa(int(pbs.Meb_leaveRoom)),
				UserId:   userId,
				RoomNo:   userInfo.UserExt.RoomNo,
			}
			responseHeadByte, _ := json.Marshal(msgData)

			netMessageResp := helper.NewNetMessage(userId, "", int32(pbs.Meb_leaveRoom), config.SlotServer)
			netMessageResp.AckHead.Code = pbs.Code(leaveRoomNotStartGameCode)
			netMessageResp.Content = responseHeadByte

			global.GVA_LOG.Infof("LeaveRoom 离开房间的广播: %v", string(responseHeadByte))
			NatsSendAimUserMsg(trs, netMessageResp, userId)
		}
	} else {
		//游戏开始后 只能死亡 才能离开房间
		//从房间里面把用户删除
		//if userInfo.GetUserIsKilled() != 1 {
		//	SendAimUserMsg(trs, userId, models.LeaveRoomMsg{
		//		ProtoNum: models.TavernStoryLeaveRoom,
		//		UserId:   userId,
		//		RoomNo:   userInfo.UserExt.RoomNo,
		//	}, models.TavernStoryLeaveRoom, common.NotLeaveRoom)
		//	return nil, nil
		//}
		//
		////修改数据库
		//err := dao.UpdateTavernRoomUserLeave(userId, userInfo.UserExt.RoomNo, int8(models.Leave))
		//if err != nil {
		//	global.GVA_LOG.Error("LeaveRoom UpdateTavernRoomUserReady", zap.Error(err))
		//	msgData := models.LeaveRoomMsg{
		//		ProtoNum: models.TavernStoryLeaveRoom,
		//		UserId:   userId,
		//		RoomNo:   userInfo.UserExt.RoomNo,
		//	}
		//	SendAimUserMsg(trs, userId, msgData, models.TavernStoryLeaveRoom, common.ModelDeleteError)
		//}
		//
		////更新用户维度状态 离开后不能在回到房间
		//updateMap := dao.MakeUpdateData("room_no", "")
		//updateMap["is_leave"] = 1
		//updateMap["is_owner"] = 0
		//updateMap["seat"] = 0
		//dao.UpdateTavernUsersRoomRoomNo(userId, updateMap)
		//
		////设置属性 离开房间标识
		//userInfo.SetUserIsLeave(models.Leave)
		//
		////发送广播 离开房间
		//msgData := models.LeaveRoomMsg{
		//	ProtoNum:  models.TavernStoryLeaveRoom,
		//	UserId:    userId,
		//	RoomNo:    userInfo.UserExt.RoomNo,
		//	Timestamp: time.Now().Unix(),
		//}
		//
		//SendAllUserMsg(trs, msgData, models.TavernStoryLeaveRoom, common.OK)
	}

	return nil, nil
}
