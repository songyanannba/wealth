package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models/table"
)

func InitTavernUsersRoom(userID string) {
	tavernUserStatus, err := table.GetUsersRoomByUid(userID)
	if err != nil {
		global.GVA_LOG.Error("InitTavernUsersRoom GetTavernUsersRoomByUid", zap.Error(err))
	}

	if tavernUserStatus.ID <= 0 {
		room := &table.UserRoom{
			UserId:   userID,
			IsLeave:  0,
			IsKilled: 0,
			IsOwner:  0,
			Turn:     0,
			Seat:     0,
			RoomNo:   "",
			Nickname: "",
			IsRobot:  0,
			IsReady:  0,
			Date:     helper.YearMonthDayStr(),
			DateTime: helper.LocalTime(),
		}
		err := table.CreateUsersRoom(room)
		if err != nil {
			global.GVA_LOG.Error("InitTavernUsersRoom CreateOrUpdateUsersRoom", zap.Error(err))
		}
	}
}

// CreateOrUpdateUsersRoom 更新用户维度的数据
func CreateOrUpdateUsersRoom(record *table.UserRoom) error {
	//创建房间的时候 ｜ 快速开始的时候
	//加入房间的时候调用
	//离开房间的时候
	tavernUserStatus, err := table.GetUsersRoomByUid(record.UserId)
	if err != nil {
		global.GVA_LOG.Error("CreateOrUpdateUsersRoom GetTavernUsersRoomByUid", zap.Error(err))
	}
	if tavernUserStatus.ID > 0 {
		//更

		val := &table.UserRoom{
			GVA_MODEL: table.GVA_MODEL{
				ID: tavernUserStatus.ID,
			},
			UserId:   record.UserId,
			RoomNo:   record.RoomNo,
			Nickname: record.Nickname,
			Seat:     record.Seat,
			Turn:     record.Turn,
			IsLeave:  record.IsLeave,
			IsKilled: record.IsKilled,
			IsOwner:  record.IsOwner,
			Date:     helper.YearMonthDayStr(),
			DateTime: helper.LocalTime(),
		}
		err := table.SaveUsersRoom(val)
		if err != nil {
			global.GVA_LOG.Error("CreateOrUpdateUsersRoom tavernUserStatus", zap.Error(err))
		}
	} else {
		return table.CreateUsersRoom(record)
	}
	return nil
}

func UpdateUsersRoomRoomNo(userId string, updateMap map[string]interface{}) {
	global.GVA_LOG.Info("UpdateUsersRoomRoomNo", zap.Any("updateMap", updateMap))

	usersRoom, err := table.GetUsersRoomByUid(userId)
	if err != nil {
		global.GVA_LOG.Error("UpdateUsersRoomRoomNo GetUsersRoomByUid", zap.Error(err))
	}
	if usersRoom.ID > 0 {
		global.GVA_LOG.Infof("UpdateUsersRoomRoomNo uid:%v ", userId)

		err := table.UpdateUsersRoom(userId, usersRoom.ID, updateMap)
		if err != nil {
			global.GVA_LOG.Error("UpdateUsersRoomRoomNo UpdateUsersRoom", zap.Error(err))
		}
	}
}
