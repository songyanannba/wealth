package table

import (
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `mb_card_config` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`level` int NOT NULL DEFAULT '0' COMMENT '等级 1=等级1 ',
//`version` int NOT NULL DEFAULT '0' COMMENT '卡牌包版本 0=基础牌 1=V1;2=V2',
//`class` int NOT NULL DEFAULT '0' COMMENT '分类 0=图片',
//`name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称(唯一标识)',
//`suffix_name` varchar(255) NOT NULL DEFAULT '' COMMENT '后缀',
//`f_img` varchar(255) NOT NULL DEFAULT '' COMMENT '链接',
//`add_rate` decimal(10,2) NOT NULL DEFAULT '1.00' COMMENT '加成比例',
//`date_time` datetime DEFAULT NULL,
//`desc` varchar(300) DEFAULT NULL COMMENT '介绍',
//`state` int NOT NULL DEFAULT '0' COMMENT '状态 11=删除',
//`unique` varchar(255) NOT NULL DEFAULT '' COMMENT '系统唯一标识',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `level` (`level`) USING BTREE,
//KEY `version` (`version`) USING BTREE,
//KEY `class` (`class`) USING BTREE,
//KEY `name` (`name`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='卡包配置';

type MbCardConfig struct {
	GVA_MODEL
	Level      int        `json:"level" form:"level" gorm:"column:level;default:0;comment:等级 1=流辉级 2=幻彩级 3=璀璨"`
	Version    int        `json:"version" form:"version" gorm:"column:version;default:0;comment:卡牌包版本 0=基础牌 1=V1;2=V2"`
	Class      int        `json:"class" form:"class" gorm:"column:class;default:0;comment:分类 0=图片"`
	Name       string     `json:"name" form:"name" gorm:"column:name;comment:名称(唯一标识);"`
	Unique     string     `json:"unique" form:"unique" gorm:"column:unique;comment:系统唯一标识;"`
	SuffixName string     `json:"suffix_name" form:"suffix_name" gorm:"column:suffix_name;comment:后缀 png;"`
	FImg       string     `json:"f_img" form:"f_img" gorm:"column:f_img;comment:链接;"`
	AddRate    float64    `json:"add_rate" form:"add_rate" gorm:"column:add_rate;comment:加成比例;"`
	DateTime   *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
	State      int        `json:"state" form:"state" gorm:"column:state;default:0;comment:状态 11=删除"`
	Desc       string     `json:"desc" form:"date" gorm:"column:desc;comment:介绍;"`
}

func (o *MbCardConfig) TableName() string {
	return "mb_card_config"
}

func GetMbCardConfigByVersion(version int) (record []*MbCardConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbCardConfig{}).
		Where("version = ?", version).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetMbCardConfigByLastId(lastId, level int) (record []*MbCardConfig, err error) {
	db := global.GVA_SLOT_SERVER_DB.
		Model(MbCardConfig{})

	if level == 0 {
		db.Where(" id > ? ", lastId)
	} else {
		db.Where(" id > ? and level = ? ", lastId, level)
	}

	err = db.Limit(20).Find(&record).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CardConfigPageIsNext(lastId, level int) (count int64, err error) {
	db := global.GVA_SLOT_SERVER_DB.Model(MbCardConfig{})

	if lastId != 0 {
		if level == 0 {
			db.Where(" id > ? ", lastId)
		} else {
			db.Where(" id > ? and level = ? ", lastId, level)
		}
	}

	err = db.Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return count, err
	}
	return count, nil
}

func CardConfigCount(level int) (count int64, err error) {
	db := global.GVA_SLOT_SERVER_DB.Model(MbCardConfig{})
	if level > 0 {
		db.Where("level = ? ", level)
	}
	err = db.Count(&count).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return count, err
	}
	return count, nil
}

func GetMbCardConfigByIds(ids []int) (record []*MbCardConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbCardConfig{}).
		Where("id IN (?)", ids).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
