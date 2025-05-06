package websocket

import (
	"encoding/json"
	"gateway/common"
	"gateway/global"
	"gateway/models"
	"go.uber.org/zap"
)

func MtHeartResp(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	return uint32(msgId), common.WebOK, nil
}

func MemeCreateRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeCreateRoomAck %v", string(message))
	msgData := &models.CreateRoomMsg{}
	//responseHead := models.NewResponseHead("", models.TavernStoryCreateRoom, common.OK, "", msgData)
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeCreateRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeJoinRoomRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeJoinRoomRoomAck %v", string(message))
	msgData := &models.JoinRoomMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeJoinRoomRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeReadyRoomRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeReadyRoomRoomAck %v", string(message))
	msgData := &models.ReadyMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeReadyRoomRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeCancelReadyRoomRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeCancelReadyRoomRoomAck %v", string(message))
	msgData := &models.ReadyMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeCancelReadyRoomRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeLeaveRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLeaveRoomAck %v", string(message))
	msgData := &models.LeaveRoomMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLeaveRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeUserStateAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeUserStateAck %v", string(message))
	msgData := &models.UserStateMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeUserStateAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeReJoinRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeReJoinRoomAck %v", string(message))
	msgData := &models.JoinRoomMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeReJoinRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeKickRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeKickRoomAck %v", string(message))
	msgData := &models.KickRoomMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeKickRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeInviteFriendAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeInviteFriendAck %v", string(message))
	msgData := &models.InviteFriendMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeInviteFriendAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeStartPlayAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeStartPlayAck %v", string(message))
	msgData := &models.StartPlayMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeStartPlayAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeLoadCompletedAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLoadCompletedAck %v", string(message))
	msgData := &models.LoadMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLoadCompletedAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeOperateCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeOperateCardsAck %v", string(message))
	msgData := &models.OperateCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeOperateCardsAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeLookCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLookCardsAck %v", string(message))
	msgData := &models.OperateCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLookCardsAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeOpeEmojiAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeOpeEmojiAck %v", string(message))
	msgData := &models.OperateCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeOpeEmojiAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeOutCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeOutCardsAck %v", string(message))
	msgData := &models.OperateCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeOutCardsAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeReMakeCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeReMakeCardsAck %v", string(message))
	msgData := &models.DealCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeReMakeCardsAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeEntryLikePageAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeEntryLikePageAck %v", string(message))
	msgData := &models.EntryLikePageMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeEntryLikePageAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeLikeCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLikeCardsAck %v", string(message))
	msgData := &models.LikeCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLikeCardsAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeCalculateRankAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeCalculateRankAck %v", string(message))
	msgData := &models.CalculateRankMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeCalculateRankAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeMatchRoomAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeMatchRoomAck %v", string(message))
	msgData := &models.MatchSuccResp{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeMatchRoomAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeMatchStartAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeMatchStartAck %v", string(message))
	msgData := &models.MatchSuccResp{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeMatchStartAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

// MemeCancelMatchAck 取消匹配
func MemeCancelMatchAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeCancelMatchAck %v", string(message))
	msgData := &models.MatchSuccResp{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeCancelMatchAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeBattleOverAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeBattleOverAck %v", string(message))
	msgData := &models.GameOverMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeBattleOverAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeIssueMsgAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLoadCompletedAck %v", string(message))
	msgData := &models.IssueMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLoadCompletedAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeDealCardsAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeLoadCompletedAck %v", string(message))
	msgData := &models.DealCardsMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeLoadCompletedAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeRoomAliveAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeRoomAliveAck %v", string(message))
	msgData := &models.UserStateMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeRoomAliveAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeHandbookListAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeHandbookListAck %v", string(message))
	msgData := &models.HandbookListMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeHandbookListAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeUnpackCardAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeUnpackCardAck %v", string(message))
	msgData := &models.UnpackCardMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeUnpackCardAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}

func MemeCardVersionListAck(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{}) {
	global.GVA_LOG.Infof("MemeCardVersionListAck %v", string(message))
	msgData := &models.CardVersionListMsg{}
	err := json.Unmarshal(message, msgData)
	if err != nil {
		global.GVA_LOG.Error("MemeCardVersionListAck Unmarshal fail", zap.Error(err))
	}
	return uint32(msgId), common.OK, msgData
}
