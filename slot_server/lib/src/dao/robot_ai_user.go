package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
)

func GetRobotAiUser() (record []*table.RobotAiUser, err error) {
	records, err := table.GetRobotAiUser()
	if err != nil {
		global.GVA_LOG.Error("get robot ai user error", zap.Error(err))
	}
	return records, nil
}
