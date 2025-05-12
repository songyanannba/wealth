package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `login_state` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '自定义用户ID',
//`session_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '平台返回 如抖音 ;会话密钥，如果请求时有 code 参数才会返回',
//`anonymous_openid` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '抖音:匿名用户在当前小游戏的 ID，如果请求时有 anonymous_code 参数才会返回',
//`unionid` varchar(255) NOT NULL DEFAULT '' COMMENT '用户在小游戏平台的唯一标识符，请求时有 code 参数才会返回。如果开发者拥有多个小游戏，可通过 unionid 来区分用户的唯一性。',
//`open_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT '平台唯一标识 ；平台返回;用户在当前小游戏的 ID，如果请求时有 code 参数才会返回',
//`is_on_line` int unsigned NOT NULL DEFAULT '0' COMMENT '0=显示 1=在线',
//`device` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机设备号',
//`dy_get_access_token` json DEFAULT NULL COMMENT 'getAccessToken接口返回 access_token 是小游戏的全局唯一调用凭据',
//`system` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '手机系统，安卓或者ios',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT NULL,
//PRIMARY KEY (`id`),
//UNIQUE KEY `user_id` (`user_id`) USING BTREE
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';

type LoginState struct {
	GVA_MODEL
	UserId           string     `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户id;"`
	IsOnLine         int        `json:"is_on_line" form:"is_on_line" gorm:"column:is_on_line;default:0;comment:0=不在 1=在线"`
	Device           string     `json:"device" form:"device" gorm:"column:device;comment:device;"`
	SessionKey       string     `json:"session_key" form:"session_key" gorm:"column:session_key;comment:平台返回;"`
	OpenId           string     `json:"open_id" form:"open_id" gorm:"column:open_id;comment:平台返回;"`
	AnonymousOpenid  string     `json:"anonymous_openid" form:"anonymous_openid" gorm:"column:anonymous_openid;comment:平台返回;"`
	System           string     `json:"system" form:"system" gorm:"column:system;comment:system;"`
	DyGetAccessToken string     `json:"dy_get_access_token" form:"dy_get_access_token" gorm:"column:dy_get_access_token;comment:getAccessToken接口返回 access_token 是小游戏的全局唯一调用凭据;"`
	DateTime         *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 投注配置
func (o *LoginState) TableName() string {
	return "login_state"
}

func CreateLoginState(record *LoginState) error {
	err := global.GVA_SLOT_SERVER_DB.Model(LoginState{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql LoginState error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveLoginState(record *LoginState) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(LoginState{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql LoginState error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetLoginStateByUid(uid string) (userInfo *LoginState, err error) {
	err = global.GVA_SLOT_SERVER_DB.Model(LoginState{}).
		Where("user_id = ?", uid).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetLoginStateByUid", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}
