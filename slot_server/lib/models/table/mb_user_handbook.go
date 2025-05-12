package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `mb_user_handbook` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`card_id` int NOT NULL DEFAULT '0' COMMENT '表情包id',
//`user_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT '会员ID',
//`num` int NOT NULL DEFAULT '0' COMMENT '数量',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT NULL,
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB COMMENT='用户表情包图鉴';

type MbUserHandbook struct {
	GVA_MODEL
	CardId   int        `json:"card_id" form:"card_id" gorm:"column:card_id;default:0;comment:表情包ID"`
	Num      int        `json:"num" form:"num" gorm:"column:num;default:0;comment:表情包数量"`
	UserId   string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	DateTime *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

func (o *MbUserHandbook) TableName() string {
	return "mb_user_handbook"
}

func GetUserHandbooks(userId string) (record []*MbUserHandbook, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Where("user_id = ? ", userId).
		Model(MbUserHandbook{}).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

// UserIsHaveHandbook 用户是否有卡
func UserIsHaveHandbook(userId string) (record *MbUserHandbook, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbUserHandbook{}).
		Where("user_id = ? ", userId).
		First(&record).Limit(1).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetUserHandbookByCardId(userId string, cardId int) (record *MbUserHandbook, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbUserHandbook{}).
		Where("user_id = ? and card_id = ? ", userId, cardId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CreateMbUserHandbook(record *MbUserHandbook) error {
	err := global.GVA_SLOT_SERVER_DB.Model(MbUserHandbook{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateMbUserHandbook error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveMbUserHandbook(record *MbUserHandbook) error {
	err := global.GVA_SLOT_SERVER_DB.Model(MbUserHandbook{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql SaveMbUserHandbook error: %s", zap.Error(err))
		return err
	}
	return nil
}

func MbUserHandbookById(id int) (record *MbUserHandbook, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(MbUserHandbook{}).
		Where("id = ?", id).
		First(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return
	}
	return
}
