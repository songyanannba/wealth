package table

import (
	"gateway/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//CREATE TABLE `game_user` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) NOT NULL DEFAULT '' COMMENT '会员ID',
//`nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '昵称',
//`king_coin` decimal(12,2) DEFAULT NULL COMMENT '冗余国王积分',
//`token` varchar(255) NOT NULL DEFAULT '' COMMENT '用户token',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT NULL,
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB  COMMENT='用户表';

type GameUser struct {
	GVA_MODEL
	UserId   string  `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	Nickname string  `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	KingCoin float64 `json:"king_coin" form:"king_coin" gorm:"column:king_coin;comment:冗余国王积分;"`
	Token    string  `json:"token" form:"token" gorm:"column:token;comment:token;"`
}

// TableName 投注配置
func (o *GameUser) TableName() string {
	return "game_user"
}

func SaveGameUser(record *GameUser) error {
	global.GVA_LOG.Infof("SaveGameUser %v", zap.Any("record", *record))
	err := global.GVA_USER_DB.Model(GameUser{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("save sql SaveGameUser error: %s", zap.Error(err))
		return err
	}
	return nil
}

func CreateGameUser(record *GameUser) error {
	global.GVA_LOG.Infof("CreateGameUser %v", zap.Any("CreateGameUser", record))
	err := global.GVA_USER_DB.Model(GameUser{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateGameUser error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetAllUser() (userList []*GameUser, err error) {
	err = global.GVA_USER_DB.Model(GameUser{}).
		Find(&userList).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return userList, err
	}
	return userList, nil
}

func GetGameUserByUid(uid string) (userInfo *GameUser, err error) {
	err = global.GVA_USER_DB.Model(GameUser{}).
		Where("user_id = ?", uid).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetGameUserByUid", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}
