package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `meme_room_config` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`bet` decimal(10,2) DEFAULT NULL COMMENT '押注',
//`admission_price` decimal(10,2) DEFAULT NULL COMMENT '入场费',
//`room_level` tinyint NOT NULL DEFAULT '0' COMMENT '房间等级 0=初级 1=中级 2=高级',
//`desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `room_level` (`room_level`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='房间配置';

type MemeRoomConfig struct {
	GVA_MODEL
	Bet            float64 `json:"bet" form:"bet" gorm:"column:bet;comment:押注;"`
	AdmissionPrice float64 `json:"admission_price" form:"admission_price" gorm:"column:admission_price;comment:押注;"`
	RoomLevel      int8    `json:"room_level" form:"room_level" gorm:"column:room_level;default:0;comment:房间等级 0 初级 1 中级 2 高级;"`
	Desc           string  `json:"desc" form:"desc" gorm:"column:desc;comment:年月日;"`
}

// TableName 房间配置
func (o *MemeRoomConfig) TableName() string {
	return "meme_room_config"
}

func GetRoomConfigs() (records []*MemeRoomConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(MemeRoomConfig{}).
		Find(&records).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("sql GetRoomConfigs error: %s", zap.Error(err))
		return records, err
	}
	return records, nil
}

func GetRoomConfigByRoomLevel(roomLevel int8) (record *MemeRoomConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(MemeRoomConfig{}).
		Where("room_level=?", roomLevel).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("sql GetRoomConfigByRoomLevel error: %s", zap.Error(err))
		return record, err
	}
	return record, nil
}
