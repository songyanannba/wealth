package logic

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
)

func BatchUpdateUserHandbook(userId string, handListCard []*models.HandListCard) {
	for _, card := range handListCard {
		UpdateUserHandbook(userId, card.CardId)
	}
}

func UpdateUserHandbook(userId string, cardId int) {
	record, err := table.GetUserHandbookByCardId(userId, cardId)
	if err != nil {
		global.GVA_LOG.Error("UpdateUserHandbook", zap.Error(err))
		return
	}
	if record.ID > 0 {
		//修改
		record.Num += 1
		err := table.SaveMbUserHandbook(record)
		if err != nil {
			global.GVA_LOG.Error("UpdateUserHandbook SaveMbUserHandbook", zap.Error(err))
		}
	} else {
		err := table.CreateMbUserHandbook(&table.MbUserHandbook{
			CardId:   cardId,
			Num:      1,
			UserId:   userId,
			DateTime: helper.LocalTime(),
		})
		if err != nil {
			global.GVA_LOG.Error("UpdateUserHandbook CreateMbUserHandbook", zap.Error(err))
		}
	}
}
