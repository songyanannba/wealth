package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `turn_details` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`turn` int NOT NULL DEFAULT '0' COMMENT '游戏局数',
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户ID',
//`nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
//`room_no` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '房间号',
//`state` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0',
//`date` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `room_no` (`room_no`) USING BTREE,
//KEY `user_id` (`user_id`) USING BTREE,
//KEY `date` (`date`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='游戏局数结算详情';

type TurnDetails struct {
	GVA_MODEL
	Turn     int        `json:"turn" form:"turn" gorm:"column:turn;default:0;comment:游戏局数;"`
	UserId   string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	Nickname string     `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	State    uint8      `json:"state" form:"state" gorm:"column:state;default:0;comment:0:生存;"`
	RoomNo   string     `json:"room_no" form:"room_no" gorm:"column:room_no;comment:房间号;"`
	Date     string     `json:"date" form:"date" gorm:"column:date;comment:年月日;"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 游戏局数结算详情
func (o *TurnDetails) TableName() string {
	return "turn_details"
}

func CreateTurnDetails(record *TurnDetails) error {
	err := global.GVA_SLOT_SERVER_DB.Model(TurnDetails{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateTurnDetails error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetTurnDetailsByNoAndUidAndTurn(roomNo, userId string, turn int) (record *TurnDetails, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(TurnDetails{}).
		Where("room_no = ? and user_id = ? and turn = ? ", roomNo, userId, turn).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(" sql GetTurnDetailsByNoAndUidAndTurn error: %s", zap.Error(err))
		return record, err
	}
	return record, nil
}

func GetTurnDetails(roomNo string) (records []*TurnDetails, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(TurnDetails{}).
		Where("room_no = ?", roomNo).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("insert sql GetTurnDetails error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}

func SaveTurnDetails(record *TurnDetails) error {
	err := global.GVA_SLOT_SERVER_DB.Model(TurnDetails{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql SaveTurnDetails error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateTurnDetails(id int, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(TurnDetails{}).
		Where("id = ?", id).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UpdateTurnDetails error: %s", zap.Error(err))
		return err
	}
	return nil
}
