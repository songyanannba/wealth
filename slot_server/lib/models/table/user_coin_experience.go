package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `user_coin_experience` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户',
//`coin_num` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '积分总数量',
//`experience` decimal(12,2) NOT NULL DEFAULT '0.00' COMMENT '经验总数量',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT NULL,
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB  COMMENT='用户的金币和经验';

type UserCoinExperience struct {
	GVA_MODEL
	UserId     string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	CoinNum    float64    `json:"coin_num" form:"coin_num" gorm:"column:coin_num;comment:积分总数量;"`
	Experience float64    `json:"experience" form:"experience" gorm:"column:experience;comment:经验总数量;"`
	DateTime   *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

func (o *UserCoinExperience) TableName() string {
	return "user_coin_experience"
}

func GetUserCoinExperience(userId string) (record *UserCoinExperience, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserCoinExperience{}).
		Where("user_id = ?", userId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CreateUserCoinExperience(record *UserCoinExperience) error {
	err := global.GVA_SLOT_SERVER_DB.Model(UserCoinExperience{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserCoinExperience error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveUserCoinExperience(record *UserCoinExperience) error {
	err := global.GVA_SLOT_SERVER_DB.Model(UserCoinExperience{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql UserCoinExperience error: %s", zap.Error(err))
		return err
	}
	return nil
}

func UpdateUserCoinExperience(uid string, values map[string]interface{}) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserCoinExperience{}).
		Where("user_id = ? ", uid).
		Updates(values).
		Error
	if err != nil {
		global.GVA_LOG.Error("UserCoinExperience error: %s", zap.Error(err))
		return err
	}
	return nil
}
