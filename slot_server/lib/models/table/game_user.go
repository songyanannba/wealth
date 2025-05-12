package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `game_user` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '自定义用户ID',
//`nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '昵称',
//`king_coin` decimal(12,2) DEFAULT NULL COMMENT '冗余国王积分',
//`token` varchar(255) NOT NULL DEFAULT '' COMMENT '用户token',
//`platform` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '平台',
//`unique_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '唯一标识',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT NULL,
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB  COMMENT='用户表';

type GameUser struct {
	GVA_MODEL
	UserId    string  `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	Nickname  string  `json:"nickname" form:"nickname" gorm:"column:nickname;comment:昵称;"`
	Platform  string  `json:"platform" form:"platform" gorm:"column:platform;comment:platform;"`
	UniqueId  string  `json:"unique_id" form:"unique_id" gorm:"column:unique_id;comment:unique_id;"`
	KingCoin  float64 `json:"king_coin" form:"king_coin" gorm:"column:king_coin;comment:冗余国王积分;"`
	Token     string  `json:"token" form:"token" gorm:"column:token;comment:token;"`
	AvatarUrl string  `json:"avatar_url" form:"avatar_url" gorm:"column:avatar_url;comment:avatar_url;"`
	OpenId    string  `json:"open_id" form:"open_id" gorm:"column:open_id;comment:平台唯一标识;"`
}

// TableName 投注配置
func (o *GameUser) TableName() string {
	return "game_user"
}

func GetGameUserByUid(uid string) (userInfo *GameUser, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(GameUser{}).
		Where("user_id = ?", uid).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetGameUserByUid", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}
