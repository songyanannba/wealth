package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
)

func GetCardConfVersionByVersion(version int) (record []*table.MbCardConfVersion) {
	byVersion, err := table.CardConfVersionListByVersion(version)
	if err != nil {
		global.GVA_LOG.Error("CardConfVersionListByVersion", zap.Any("err", err))
	}
	return byVersion
}

func GetCardConfVersion() (record []*table.MbCardConfVersion) {
	cardConfVersionList, err := table.CardConfVersionList()
	if err != nil {
		global.GVA_LOG.Error("GetCardConfVersion", zap.Any("err", err))
	}
	return cardConfVersionList
}
