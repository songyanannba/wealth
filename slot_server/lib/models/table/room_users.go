package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"time"
)

//CREATE TABLE `room_users` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户ID',
//`win_price` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '赢钱',
//`bet` decimal(10,2) NOT NULL DEFAULT '0.00' COMMENT '押注',
//`is_leave` tinyint(1) NOT NULL DEFAULT '0' COMMENT '1:离开 0 未离开',
//`is_killed` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被杀 1 被杀 0 没有被杀',
//`is_owner` int NOT NULL DEFAULT '0' COMMENT '1：房主 0：不是房主',
//`turn` int NOT NULL DEFAULT '0' COMMENT '第几轮',
//`seat` int NOT NULL DEFAULT '0' COMMENT '房间座位',
//`room_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间号',
//`nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
//`is_ready` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否就绪 0:未就绪 1:就绪',
//`is_robot` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否机器人:0=否,1=是',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `room_no` (`room_no`) USING BTREE,
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='游戏玩家表';

const NotBeOwner = 0
const BeOwner = 1

type RoomUsers struct {
	GVA_MODEL
	UserId   string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	RoomNo   string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	Nickname string     `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	IsRobot  int8       `json:"is_robot" form:"is_robot" gorm:"column:is_robot;default:0;comment:是否机器人:0=否,1=是;"`
	IsReady  int8       `json:"is_ready" form:"is_ready" gorm:"column:is_ready;default:0;comment:是否就绪 0:未就绪 1:就绪;"`
	Seat     int        `json:"seat" form:"seat" gorm:"column:seat;default:0;comment:座位次序"`
	Date     string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
	Turn     int        `json:"turn" form:"turn" gorm:"column:turn;default:0;comment:第几轮"`
	IsLeave  int        `json:"is_leave" form:"is_leave" gorm:"column:is_leave;default:0;comment:1:离开 0:未离开"`
	IsKilled int        `json:"is_killed" form:"is_killed" gorm:"column:is_killed;default:0;comment:是否被杀 1:被杀 0:没有被杀"`
	IsOwner  int        `json:"is_owner" form:"is_owner" gorm:"column:is_owner;default:0;comment:1:房主 0:不是房主"`
	WinPrice float64    `json:"win_price" form:"win_price" gorm:"column:win_price;comment:赢钱;"`
	Bet      float64    `json:"bet" form:"bet" gorm:"column:bet;comment:押注;"`
}

// TableName 游戏玩家表
func (o *RoomUsers) TableName() string {
	return "room_users"
}

func NewRoomUsers(userId, roomNo, nickname string, seat, turn, isOwner int, bet float64, isReady int8) *RoomUsers {
	return &RoomUsers{
		UserId:   userId,
		RoomNo:   roomNo,
		Nickname: nickname,
		Seat:     seat,
		Turn:     turn,
		IsOwner:  isOwner,
		Bet:      bet,
		IsRobot:  0,
		IsReady:  isReady,
		WinPrice: 0,
		IsLeave:  0,
		IsKilled: 0,
		Date:     helper.YearMonthDayStr(),
		DateTime: helper.LocalTime(),
	}
}

func NewestRoomUsersByUid(userId string) (record *RoomUsers, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("user_id = ?", userId).
		Order("id desc").
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CreateRoomUsers(record *RoomUsers) error {
	err := global.GVA_SLOT_SERVER_DB.Model(RoomUsers{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateRoomUsers error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetRoomUsers(roomNo string) (records []*RoomUsers, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("room_no = ? and is_leave = ?", roomNo, 0).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" sql GetTavernRoomUsers error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}

func RoomUsersByRoomNoAndUid(roomNo, uid string) (record *RoomUsers, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(RoomUsers{}).
		Where("room_no = ? and user_id = ? ", roomNo, uid).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" sql RoomUsersByRoomNoAndUid error: %s", zap.Error(err))
		return record, err
	}
	return record, nil
}

func DelRoomUsers(roomNo, userId string) error {
	global.GVA_LOG.Infof("DelRoomUsers %v %v", roomNo, userId)
	err := global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("room_no = ? and user_id = ? ", roomNo, userId).
		Delete(RoomUsers{}).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql DelRoomUsers error: %s", zap.Error(err))
		return err
	}
	return nil
}

func DelRoomUsersByRoomNo(roomNo string) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("room_no = ?", roomNo).
		Delete(RoomUsers{}).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql DelRoomUsersByRoomNo error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateRoomUsersReady(uid, roomNo string, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("user_id = ? and  room_no = ?", uid, roomNo).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateRoomUsersReady error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateRoomUsersLeave(uid, roomNo string, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("user_id = ? and  room_no = ?", uid, roomNo).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateSwingRodRecordBySwingRodNo error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateRoomUser(uid, roomNo string, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("user_id = ? and  room_no = ?", uid, roomNo).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateTavernRoomUser error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetRoomUser(uid, roomNo string) (record *RoomUsers, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(RoomUsers{}).
		Where("user_id = ? and room_no = ?", uid, roomNo).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

//func UpdateTavernRoomWinPrice(uid, roomNo string, values map[string]interface{}) error {
//	err := global.GVA_SLOT_SERVER_DB.
//		Model(RoomUsers{}).
//		Where("user_id = ? and  room_no = ?", uid, roomNo).
//		Updates(values).
//		Error
//	if err != nil {
//		global.GVA_LOG.Error("UpdateTavernRoomWinPrice error: %s", zap.Error(err))
//		return err
//	}
//	return nil
//}
