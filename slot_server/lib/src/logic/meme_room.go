package logic

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"slot_server/lib/common"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/lib/src/dao"
)

func SaveRoom(userID string, roomType, userNumLimit, roomLevel, roomTurnNum int, bet float64) (*table.MemeRoom, error) {
	room := &table.MemeRoom{}

	userInfo, err := table.GetGameUserByUid(userID)
	if err != nil {
		global.GVA_LOG.Error("SaveRoom GetGameUserByUid", zap.Error(err))
		return room, err
	}

	//创建
	rNo := uuid.New().String()
	rName := userInfo.Nickname + "'s" + " lobby"
	room = table.NewMemeRoom(userID, userID, rNo, rName, "", table.TavernRoomOpen, int8(roomType), int8(roomLevel), table.RoomClassInvite, roomTurnNum, userNumLimit)
	err = table.CreateMemeRoom(room)
	if err != nil {
		global.GVA_LOG.Error("SaveRoom CreateMemeRoom", zap.Error(err))
		return room, err
	}

	//添加房间用户
	tavernRoomUsers := table.NewRoomUsers(userID, room.RoomNo, userInfo.Nickname, 1, 1, 1, bet, 0)
	err = table.CreateRoomUsers(tavernRoomUsers)
	if err != nil {
		global.GVA_LOG.Error("SaveRoom CreateTavernRoomUsers", zap.Error(err))
		return room, err
	}

	//用户维度状态 ｜ 创建房间
	tavernUserRoomData := table.NewUserRoom(tavernRoomUsers.UserId, tavernRoomUsers.RoomNo, tavernRoomUsers.Nickname, tavernRoomUsers.IsLeave, tavernRoomUsers.IsKilled,
		tavernRoomUsers.IsOwner, tavernRoomUsers.Turn, tavernRoomUsers.Seat, tavernRoomUsers.IsRobot, tavernRoomUsers.IsReady)
	err = dao.CreateOrUpdateUsersRoom(tavernUserRoomData)
	if err != nil {
		global.GVA_LOG.Error("SaveRoom", zap.Any("CreateOrUpdateUsersRoom", err))
	}

	return room, nil
}

func IsCanCreateOrJoinRoom(userID string) int {
	code := common.OK
	//房间用户维度
	roomDetail := UserIsJoinRoom(userID)
	//先离开原来的房间
	if roomDetail != nil && len(roomDetail.RoomNo) != 0 && roomDetail.Status > 0 {
		//通知房间谁必须先离开房间
		code = common.LeavePreRoom
		global.GVA_LOG.Error("IsCanCreateOrJoinRoom 先离开原来的房间 ")
		return code
	}

	//房间维度在检查一次
	isCan, err := IsCanCreateTavernRoom(userID)
	if err != nil {
		code = common.LeavePreRoom
		global.GVA_LOG.Error("IsCanCreateOrJoinRoom IsCanCreateTavernRoom ", zap.Error(err))
		return code
	}
	if !isCan {
		code = common.LeavePreRoom
		global.GVA_LOG.Error("IsCanCreateOrJoinRoom 先离开原来的房间 ", zap.Error(err))
		return code
	}
	return code
}

// UserIsJoinRoom 用户如果已经加入房间 返回用户所在房间
func UserIsJoinRoom(userId string) *models.RoomItem {
	roomDetail := &models.RoomItem{}
	tavernRoom, err := dao.UserIsJoinRoom(userId)
	if err != nil {
		global.GVA_LOG.Error("UserIsJoinRoom", zap.Error(err))
	}
	if tavernRoom == nil || tavernRoom.ID <= 0 || tavernRoom.RoomNo == "" {
		return roomDetail
	}

	roomDetail = &models.RoomItem{
		RoomCom: models.RoomCom{
			RoomId:       tavernRoom.ID,
			RoomNo:       tavernRoom.RoomNo,
			RoomName:     tavernRoom.Name,
			Status:       tavernRoom.IsOpen,
			UserNumLimit: tavernRoom.UserNumLimit,
			RoomType:     int(tavernRoom.RoomType),
			RoomLevel:    int(tavernRoom.RoomLevel),
		},
		RoomUserList: nil,
	}

	roomUserLists, _ := dao.GetRoomUser(tavernRoom.RoomNo, 0)

	roomDetail.RoomUserList = roomUserLists
	return roomDetail
}

func IsCanCreateTavernRoom(userID string) (bool, error) {
	room, err := table.GetMemeRoomByUid(userID)
	if err != nil {
		global.GVA_LOG.Error("IsCanCreateTavernRoom GetMemeRoomByUid", zap.Error(err))
		return false, err
	}
	if room.ID > 0 {
		//查看是否在房间里面
		record, err := table.RoomUsersByRoomNoAndUid(room.RoomNo, userID)
		if err != nil {
			global.GVA_LOG.Error("SaveTavernRoomUsersByRoomNoAndUid", zap.Error(err))
			return false, err
		}
		if record.ID > 0 && record.IsLeave == 0 {
			return false, errors.New("请先离开加入的房间")
		}
	}

	room, err = table.GetMemeRoomByUid(userID)
	if err != nil {
		global.GVA_LOG.Error("SaveRoom GetTavernRoomByOwner", zap.Error(err))
		return false, err
	}
	if room.ID > 0 {
		record, err := table.RoomUsersByRoomNoAndUid(room.RoomNo, userID)
		if err != nil {
			global.GVA_LOG.Error("SaveTavernRoomUsersByRoomNoAndUid", zap.Error(err))
			return false, err
		}
		if record.ID > 0 && record.IsLeave == 0 {
			return false, errors.New("请先离开加入的房间")
		}
	}
	return true, nil
}

func GetSeat(roomNo string, preSeat int) (int, error) {
	//获取座位
	roomUsers, err := table.GetRoomUsers(roomNo)
	if err != nil {
		global.GVA_LOG.Error("JoinRoom SaveTavernRoom TavernRoomUsers", zap.Error(err))
		return preSeat, err
	}
	leaveSeatArr := []int{}
	for k, _ := range roomUsers {
		//说明作为已经被占
		if roomUsers[k].Seat != k+1 {
			leaveSeatArr = append(leaveSeatArr, k+1)
		}
	}
	if len(leaveSeatArr) > 0 {
		preSeat = leaveSeatArr[0]
	}
	return preSeat, nil
}
