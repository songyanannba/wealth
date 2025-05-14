package dao

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models/table"
)

// UserIsJoinRoom 用户如果已经加入房间 返回用户所在房间
func UserIsJoinRoom(userId string) (tavernRoom *table.AnimalPartyRoom, err error) {
	//查看用户最近一次加入的房间
	user, err := NewestRoomUser(userId)
	if err != nil {
		global.GVA_LOG.Error("UserIsJoinRoom NewestTavernRoomUser", zap.Any("err", err))
		return tavernRoom, err
	}
	if len(user.RoomNo) == 0 {
		return tavernRoom, nil
	}

	//获取最近一次加入过的房间信息
	room, err := NewestTavernRoomByRoomNo(user.RoomNo)
	if err != nil {
		global.GVA_LOG.Error("UserIsJoinRoom NewestTavernRoom", zap.Any("err", err))
		return tavernRoom, err
	}
	return room, nil
}

func NewestTavernRoomByRoomNo(roomNo string) (tavernRoom *table.AnimalPartyRoom, err error) {
	record, err := table.NewestNormalMemeRoomByRoomNo(roomNo)
	if err != nil {
		global.GVA_LOG.Error("GetTavernRoom", zap.Any("err", err))
		return tavernRoom, err
	}
	return record, nil
}
