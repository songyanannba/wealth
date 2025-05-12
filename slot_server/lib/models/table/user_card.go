package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `user_card` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户',
//`card_conf_id` int NOT NULL DEFAULT '0' COMMENT '拥有的开包版本',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `card_conf_id` (`card_conf_id`) USING BTREE,
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB  COMMENT='玩家卡牌（基础卡牌+开包卡牌）';

type UserCard struct {
	GVA_MODEL
	UserId     string    `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	CardConfId int       `json:"card_conf_id" form:"card_conf_id" gorm:"column:card_conf_id;comment:卡牌配置ID;"`
	DateTime   time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

func (o *UserCard) TableName() string {
	return "user_card"
}

func GetUserCards() (record []*UserCard, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserCard{}).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetUserCardByUserId(userId string) (record []*UserCard, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserCard{}).
		Where("userId = ?", userId).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CreateUserCard(record *UserCard) error {
	err := global.GVA_SLOT_SERVER_DB.Model(UserCard{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserCard error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveUserCard(record *UserCard) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserCard{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserCard error: %s", zap.Error(err))
		return err
	}
	return nil
}
