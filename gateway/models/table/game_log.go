package table

import (
	"gateway/global"
	"go.uber.org/zap"
)

//CREATE TABLE `lc_game_log` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0',
//`method` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
//`req` text COLLATE utf8mb4_unicode_ci COMMENT '请求',
//`resq` text COLLATE utf8mb4_unicode_ci COMMENT '返回',
//`memo` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
//`room_no` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '房间号',
//`updatetime` int DEFAULT '0' COMMENT '更新时间',
//`createtime` int DEFAULT NULL COMMENT '创建时间',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB  COMMENT='请求3方服务的日志';

type GameLog struct {
	GVA_MODEL
	UserId string `json:"user_id" form:"user_id" gorm:"column:user_id;comment:uid"`
	Method string `json:"method" form:"method" gorm:"column:method;comment:method"`
	Req    string `json:"req" form:"req" gorm:"column:req;comment:请求;"`
	Resq   string `json:"resq" form:"resq" gorm:"column:resq;comment:返回;"`
	Memo   string `json:"memo" form:"memo" gorm:"column:memo;comment:备注;"`
	RoomNo string `json:"room_no" form:"room_no" gorm:"column:room_no;comment:对局房间;"`
}

func (o *GameLog) TableName() string {
	return "lc_game_log"
}

func CreateGameLog(val *GameLog) error {
	global.GVA_LOG.Info("CreateGameLog")
	if err := global.GVA_USER_DB.
		Model(GameLog{}).
		Create(&val).Error; err != nil {
		global.GVA_LOG.Error("insert sql CreateGameLog error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateGameLog(id int, values map[string]interface{}) error {
	global.GVA_LOG.Info("UpdateGameLog")
	err := global.GVA_USER_DB.
		Model(GameLog{}).
		Where("id = ?", id).
		Updates(values).Error

	if err != nil {
		global.GVA_LOG.Error("UpdateGameLog error: %s", zap.Error(err))
		return err
	}
	return nil
}
