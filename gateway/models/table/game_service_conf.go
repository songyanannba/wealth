package table

import (
	"gateway/global"
	"gorm.io/gorm"
)

//CREATE TABLE `game_service_conf` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`maintenance` int unsigned NOT NULL DEFAULT '0' COMMENT '1 维护阶段 禁止加入游戏',
//`updatetime` int DEFAULT '0' COMMENT '更新时间',
//`createtime` int DEFAULT NULL COMMENT '创建时间',
//`deletetime` int DEFAULT NULL COMMENT '删除时间',
//`g_type` int DEFAULT '0' COMMENT '游戏类型',
//`desc` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '' COMMENT '描述',
//`g_id` int DEFAULT '0' COMMENT '游戏ID',
//PRIMARY KEY (`id`)
//) ENGINE=InnoDB  COMMENT='更新服务相关配置';

type GameServiceConf struct {
	GVA_MODEL
	GType       int    `json:"g_type" form:"g_type" gorm:"column:g_type;comment:1=游戏类型:size:64;"`
	GId         int    `json:"g_id" form:"g_id" gorm:"column:g_id;comment:游戏ID:size:64;"`
	Desc        string `json:"desc" form:"desc" gorm:"column:desc;comment:1 全部长链接服务;"`
	Maintenance int    `json:"maintenance" form:"maintenance" gorm:"column:maintenance;comment:奖池总额;"`
}

func (o *GameServiceConf) TableName() string {
	return "game_service_conf"
}

func GetGameServiceConf(id int) (val *GameServiceConf, err error) {
	err = global.GVA_USER_DB.
		Model(GameServiceConf{}).
		Where("id = ? ", id).
		First(&val).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return nil, err
	}

	global.GVA_LOG.Infof("GetGameServiceConf %v", val)
	return val, nil
}
