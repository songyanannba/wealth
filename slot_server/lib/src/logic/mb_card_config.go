package logic

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
)

func GetUserOwnCards(userId string) []*table.MbCardConfig {
	//给自己的牌赋值
	//todo 加缓存
	userHandbooks, err := table.GetUserHandbooks(userId)
	if err != nil {
		global.GVA_LOG.Error("LoadCompletedFirst...", zap.Error(err))
	}
	cardIds := []int{}
	for _, userHandbook := range userHandbooks {
		cardIds = append(cardIds, userHandbook.CardId)
	}
	cardConfigByIds, err := table.GetMbCardConfigByIds(cardIds)
	if err != nil {
		global.GVA_LOG.Error("LoadCompletedFirst...", zap.Error(err))
	}
	return cardConfigByIds
}
