package table

import (
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `mb_card_conf_version` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`version` int NOT NULL DEFAULT '0' COMMENT '卡牌包版本 0=基础牌 1=V1;2=V2',
//`desc` varchar(300) DEFAULT NULL COMMENT '介绍',
//`is_show` int NOT NULL DEFAULT '0' COMMENT '状态 1=不展示',
//PRIMARY KEY (`id`),
//KEY `version` (`version`) USING BTREE,
//KEY `state` (`is_show`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='卡包版本配置';

type MbCardConfVersion struct {
	GVA_MODEL
	Version int    `json:"version" form:"version" gorm:"column:version;default:0;comment:卡牌包版本 0=基础牌 1=V1;2=V2"`
	IsShow  int    `json:"is_show" form:"is_show" gorm:"column:is_show;default:0;comment:状态 状态 1=不展示"`
	Desc    string `json:"desc" form:"date" gorm:"column:desc;comment:介绍;"`
}

func (o *MbCardConfVersion) TableName() string {
	return "mb_card_conf_version"
}

// CardConfVersionListByVersion 版本列表
func CardConfVersionListByVersion(version int) (record []*MbCardConfVersion, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbCardConfVersion{}).
		Where("version = ? and is_show = 0", version).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CardConfVersionList() (record []*MbCardConfVersion, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbCardConfVersion{}).
		Where("is_show = 0").
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
