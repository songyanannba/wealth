package dao

import (
	"errors"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
)

// UpdateRoomUserReady 用户就绪
func UpdateRoomUserReady(uid, roomNo string, isReady int8) error {
	global.GVA_LOG.Infof("UpdateRoomUserReady uid:%v ,roomNo:%v isReady:%v", uid, roomNo, isReady)

	updateMap := MakeUpdateData("is_ready", isReady)

	err := table.UpdateRoomUsersReady(uid, roomNo, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomUserReady ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateRoomUserLeave 用户离开
func UpdateRoomUserLeave(uid, roomNo string, isLeave int8) error {
	global.GVA_LOG.Infof("UpdateRoomUserLeave uid:%v ,roomNo:%v isReady:%v", uid, roomNo, isLeave)

	updateMap := MakeUpdateData("is_leave", isLeave)

	err := table.UpdateRoomUsersLeave(uid, roomNo, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomUserLeave ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateRoomOwner 更新房主信息
func UpdateRoomOwner(uid, roomNo string, isOwner int) error {
	global.GVA_LOG.Infof("UpdateRoomOwner uid:%v ,roomNo:%v isReady:%v", uid, roomNo)

	updateMap := MakeUpdateData("is_owner", isOwner)

	err := table.UpdateRoomUser(uid, roomNo, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomOwner ", zap.Error(err))
		return err
	}
	return nil
}

func UpdateRoomIsKill(uid, roomNo string) error {
	global.GVA_LOG.Infof("UpdateRoomIsKill uid:%v ,roomNo:%v isReady:%v", uid, roomNo)

	updateMap := MakeUpdateData("is_killed", 1)

	err := table.UpdateRoomUser(uid, roomNo, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomIsKill ", zap.Error(err))
		return err
	}
	return nil
}

// UpdateRoomWinPrice 最终赢钱用户更新赢钱金额
func UpdateRoomWinPrice(uid, roomNo string, winPrice float64) error {
	global.GVA_LOG.Infof("UpdateRoomWinPrice uid:%v ,roomNo:%v isReady:%v", uid, roomNo)

	updateMap := MakeUpdateData("win_price", winPrice)

	err := table.UpdateRoomUser(uid, roomNo, updateMap)
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomWinPrice ", zap.Error(err))
		return err
	}
	return nil
}

func NewestRoomUser(userId string) (tavernRoomUser *table.RoomUsers, err error) {
	record, err := table.NewestRoomUsersByUid(userId)
	if err != nil {
		global.GVA_LOG.Error("NewestTavernRoomUser", zap.Any("err", err))
		return tavernRoomUser, err
	}
	//已经离开
	if record.IsLeave == models.Leave {
		return tavernRoomUser, errors.New("用户已经离开房间")
	}
	return record, nil
}

func GetRoomUser(roomNo string, turn int) ([]models.MemeRoomUser, error) {
	roomUserLists := make([]models.MemeRoomUser, 0)
	roomUsers, err := table.GetRoomUsers(roomNo)
	if err != nil {
		global.GVA_LOG.Error("JoinRoom GetTavernRoomUsers: %v %v", zap.Error(err))
		return roomUserLists, err
	}
	for _, roomUser := range roomUsers {
		isOwner := false
		if roomUser.IsOwner == 1 {
			isOwner = true
		}
		//character, _ := table.GetUserCharacter(roomUser.UserId)
		userItem := models.MemeRoomUser{
			UserID:   roomUser.UserId,
			Seat:     roomUser.Seat,
			Turn:     turn,
			Nickname: roomUser.Nickname,
			IsOwner:  isOwner,
			//CharacterId: character.CharacterId,
			Bet:     roomUser.Bet,
			IsReady: int(roomUser.IsReady),
		}
		roomUserLists = append(roomUserLists, userItem)
	}
	return roomUserLists, nil
}
