package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
)

func GetCoinConsumeConfigByType(typ int) *table.CoinConsumeConfig {
	record, err := table.GetCoinConsumeConfigByType(typ)
	if err != nil {
		global.GVA_LOG.Error("GetCoinConsumeConfigByType failed", zap.Error(err))
	}
	return record
}
