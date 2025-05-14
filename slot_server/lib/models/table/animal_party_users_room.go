package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"time"
)

//CREATE TABLE `users_room` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户ID',
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
//UNIQUE KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB COMMENT='用户维度信息';

type AnimalPartyUserRoom struct {
	GVA_MODEL
	UserId   string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	IsLeave  int        `json:"is_leave" form:"is_leave" gorm:"column:is_leave;default:0;comment:1:离开 0:未离开"`
	IsKilled int        `json:"is_killed" form:"is_killed" gorm:"column:is_killed;default:0;comment:是否被杀 1:被杀 0:没有被杀"`
	IsOwner  int        `json:"is_owner" form:"is_owner" gorm:"column:is_owner;default:0;comment:1:房主 0:不是房主"`
	Turn     int        `json:"turn" form:"turn" gorm:"column:turn;default:0;comment:第几轮"`
	Period   int        `json:"period" form:"period" gorm:"column:period;default:0;comment:第几期"`
	Seat     int        `json:"seat" form:"seat" gorm:"column:seat;default:0;comment:房间座位"`
	RoomNo   string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	Nickname string     `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	IsRobot  int8       `json:"is_robot" form:"is_robot" gorm:"column:is_robot;default:0;comment:是否机器人:0=否,1=是;"`
	IsReady  int8       `json:"is_ready" form:"is_ready" gorm:"column:is_ready;default:0;comment:是否就绪 0:未就绪 1:就绪;"`
	Date     string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 游戏玩家表
func (o *AnimalPartyUserRoom) TableName() string {
	return "animal_party_users_room"
}

func NewUserRoom(userId, roomNo, nickname string, isLeave, isKilled, isOwner, turn, period int, isRobot, isReady int8) *AnimalPartyUserRoom {
	return &AnimalPartyUserRoom{
		UserId:   userId,
		IsLeave:  isLeave,
		IsKilled: isKilled,
		IsOwner:  isOwner,
		Turn:     turn,
		Period:   period,
		RoomNo:   roomNo,
		Nickname: nickname,
		IsRobot:  isRobot,
		IsReady:  isReady,
		Date:     helper.YearMonthDayStr(),
		DateTime: helper.LocalTime(),
	}
}

func CreateUsersRoom(record *AnimalPartyUserRoom) error {
	err := global.GVA_SLOT_SERVER_DB.Model(AnimalPartyUserRoom{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql AnimalPartyUserRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetUsersRoom(roomNo string) (records []*AnimalPartyUserRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyUserRoom{}).
		Where("room_no = ?", roomNo).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" sql AnimalPartyUserRoom error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}

func GetUsersRoomByUid(userId string) (record *AnimalPartyUserRoom, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyUserRoom{}).
		Where("user_id = ?", userId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func SaveUsersRoom(record *AnimalPartyUserRoom) error {
	err := global.GVA_SLOT_SERVER_DB.Model(AnimalPartyUserRoom{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql AnimalPartyUserRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUsersRoom(uid string, id int, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(AnimalPartyUserRoom{}).
		Where("user_id = ? and id = ?", uid, id).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateUsersRoom error: %s", zap.Error(err))
		return err
	}
	return nil
}
