package table

import (
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//
//CREATE TABLE `mb_issue_config` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`level` int NOT NULL DEFAULT '0' COMMENT '等级',
//`class` int NOT NULL DEFAULT '0' COMMENT '分类',
//`desc` varchar(1000) DEFAULT NULL COMMENT '问题描述',
//`state` int NOT NULL DEFAULT '0' COMMENT '状态 11=删除',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `level` (`level`) USING BTREE,
//KEY `class` (`class`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='问题配置';

type MbIssueConfig struct {
	GVA_MODEL
	Level    int        `json:"level" form:"level" gorm:"column:level;default:0;comment:等级"`
	Class    int        `json:"class" form:"class" gorm:"column:class;default:0;comment:分类"`
	Desc     string     `json:"desc" form:"date" gorm:"column:desc;comment:介绍;"`
	State    int        `json:"state" form:"state" gorm:"column:state;default:0;comment:状态 11=删除"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

func (o *MbIssueConfig) TableName() string {
	return "mb_issue_config"
}

func GetMbIssueConfigs() (record []*MbIssueConfig, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbIssueConfig{}).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
