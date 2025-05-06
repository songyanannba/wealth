package routers

import (
	"gateway/protoc/pbs"
	"gateway/servers/websocket"
)

// MemeBattle
type MemeBattleRouters struct {
}

var memeBattleRouters = MemeBattleRouters{}

func (rm *MemeBattleRouters) WayRouterInit() {

	//websocket.Register("398", websocket.MtTestReq)

	//magicTower
	//websocket.Register("400", websocket.MTHeart400)
	//websocket.RegisterNatsProtoResp(4, websocket.MtTestResp)

	//测试心跳返回
	websocket.RegisterNatsProtoResp(int32(pbs.Meb_mtHeart), websocket.MtHeartResp)

	//创建房间返回
	websocket.RegisterNatsProtoResp(int32(pbs.Meb_createRoom), websocket.MemeCreateRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_joinRoom), websocket.MemeJoinRoomRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_readyMsg), websocket.MemeReadyRoomRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_cancelReadyMsg), websocket.MemeCancelReadyRoomRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_leaveRoom), websocket.MemeLeaveRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_userState), websocket.MemeUserStateAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_reJoinRoom), websocket.MemeReJoinRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_kickRoom), websocket.MemeKickRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_inviteFriend), websocket.MemeInviteFriendAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_startPlay), websocket.MemeStartPlayAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_issueMsg), websocket.MemeIssueMsgAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_dealCardsMsg), websocket.MemeDealCardsAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_loadCompleted), websocket.MemeLoadCompletedAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_operateCards), websocket.MemeOperateCardsAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_lookCards), websocket.MemeLookCardsAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_opeEmoji), websocket.MemeOpeEmojiAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_outCards), websocket.MemeOutCardsAck)

	//重新随牌（洗牌）
	websocket.RegisterNatsProtoResp(int32(pbs.Meb_reMakeCards), websocket.MemeReMakeCardsAck)

	//进入点赞页面的关广播
	websocket.RegisterNatsProtoResp(int32(pbs.Meb_entryLikePage), websocket.MemeEntryLikePageAck)

	//点赞
	websocket.RegisterNatsProtoResp(int32(pbs.Meb_likeCards), websocket.MemeLikeCardsAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_calculateRank), websocket.MemeCalculateRankAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_memeMatchRoom), websocket.MemeMatchRoomAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_matchStart), websocket.MemeMatchStartAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_cancelMatchRoom), websocket.MemeCancelMatchAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_memeBattleOver), websocket.MemeBattleOverAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_roomAlive), websocket.MemeRoomAliveAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_handbookList), websocket.MemeHandbookListAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_unpackCard), websocket.MemeUnpackCardAck)

	websocket.RegisterNatsProtoResp(int32(pbs.Meb_cardVersionList), websocket.MemeCardVersionListAck)

}
