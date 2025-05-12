package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models/table"
)

// AddTurnDetails 详情
func AddTurnDetails(roomNo, userId, nickname string, turn int, cards, outCards string) {
	roomDetails, err := table.GetTurnDetailsByNoAndUidAndTurn(roomNo, userId, turn)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}

	if roomDetails.ID <= 0 {
		roomDetails = &table.TurnDetails{
			Turn:     turn,
			State:    0,
			UserId:   userId,
			RoomNo:   roomNo,
			Nickname: nickname,
			Date:     helper.YearMonthDayStr(),
			DateTime: helper.LocalTime(),
		}
		err := table.CreateTurnDetails(roomDetails)
		if err != nil {
			global.GVA_LOG.Error(" AddTurnDetails", zap.Error(err))
		}
		AddTurnDetailsExt(roomDetails.ID, roomNo, turn, cards, outCards)
	} else {
		//随牌
		AddTurnDetailsExt(roomDetails.ID, roomNo, turn, cards, outCards)
	}

}

// AddTurnDetailsExt 详情
func AddTurnDetailsExt(detailsId int, roomNo string, turn int, cards, outCards string) {
	roomDetailsExt, err := table.GetTurnDetailsExtByDetailsId(detailsId)
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}
	if roomDetailsExt.ID > 0 {
		return
	}
	ext := table.TurnDetailsExt{
		RoomNo:    roomNo,
		DetailsId: detailsId,
		Turn:      turn,
		Cards:     cards,
		OutCards:  outCards,
		Date:      helper.YearMonthDayStr(),
		DateTime:  helper.LocalTime(),
	}
	err = table.CreateTurnDetailsExt(&ext)
	if err != nil {
		global.GVA_LOG.Error(" AddTurnDetailsExt", zap.Error(err))
	}

}
