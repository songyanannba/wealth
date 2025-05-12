package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models/table"
)

func CreateCoinConsumeLog(userId string, opType int, currNum, afterNum, beforeNum float64, desc string) {
	coinConsumeLog := &table.CoinConsumeLog{
		UserId:    userId,
		OpType:    opType,
		CurrNum:   currNum,
		AfterNum:  afterNum,
		BeforeNum: beforeNum,
		Desc:      desc,
		Date:      helper.YearMonthDayStr(),
		DateTime:  helper.LocalTime(),
	}
	err := table.CreateCoinConsumeLog(coinConsumeLog)
	if err != nil {
		global.GVA_LOG.Error("create coin consume log error", zap.Any("err", err))
	}
}
