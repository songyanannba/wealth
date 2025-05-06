package table

import (
	"gateway/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//CREATE TABLE `b_user` (
//`id` bigint unsigned NOT NULL AUTO_INCREMENT,
//`uuid` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'UUID',
//`user_name` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '用户名',
//`password` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '密码',
//`nick_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '昵称',
//`header_img` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '头像',
//`phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '手机号',
//`email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '邮箱',
//`amount` decimal(10,2) DEFAULT '0.00' COMMENT '金额',
//`ip` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '注册IP',
//`last_ip` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '最后登录IP',
//`status` tinyint unsigned DEFAULT '1' COMMENT '状态 1正常 2冻结',
//`online` tinyint unsigned DEFAULT '2' COMMENT '是否在线',
//`merchant_id` int unsigned DEFAULT '0' COMMENT '商户ID',
//`city` varchar(255) DEFAULT NULL COMMENT '城市',
//`currency` varchar(12) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT 'USD' COMMENT '货币',
//`type` tinyint unsigned DEFAULT '1' COMMENT '类型',
//`indicate_num` tinyint unsigned DEFAULT '0' COMMENT 'ab测 随机编号',
//`created_at` datetime DEFAULT NULL,
//`updated_at` datetime DEFAULT NULL,
//`deleted_at` datetime DEFAULT NULL,
//PRIMARY KEY (`id`) USING BTREE,
//KEY `idx_sys_users_deleted_at` (`deleted_at`) USING BTREE,
//KEY `idx_sys_users_uuid` (`uuid`) USING BTREE,
//KEY `idx_sys_users_username` (`username`) USING BTREE,
//KEY `idx_b_user_deleted_at` (`deleted_at`) USING BTREE,
//KEY `idx_b_user_uuid` (`uuid`) USING BTREE,
//KEY `idx_b_user_username` (`username`) USING BTREE
//) ENGINE=InnoDB COMMENT='用户表';

type BUser struct {
	GVA_MODEL
	Uuid        string  `json:"uuid" form:"uuid" gorm:"column:uuid;comment:UUID;"`
	UserName    string  `json:"user_name" form:"user_name" gorm:"column:user_name;comment:用户名;"`
	Nickname    string  `json:"nick_name" form:"nick_name" gorm:"column:nick_name;comment:昵称;"`
	Password    string  `json:"password" form:"password" gorm:"column:password;comment:密码;"`
	HeaderImg   string  `json:"header_img" form:"header_img" gorm:"column:header_img;comment:头像;"`
	Phone       string  `json:"phone" form:"phone" gorm:"column:phone;comment:手机号;"`
	Email       string  `json:"email" form:"email" gorm:"column:email;comment:邮箱;"`
	Amount      float64 `json:"amount" form:"amount" gorm:"column:amount;comment:金额;"`
	Status      int     `json:"status" form:"status" gorm:"column:status;comment:状态 1正常 2冻结;"`
	Online      int     `json:"online" form:"online" gorm:"column:online;comment:是否在线;"`
	MerchantId  int     `json:"merchant_id" form:"merchant_id" gorm:"column:merchant_id;comment:商户ID;"`
	Type        int     `json:"type" form:"type" gorm:"column:type;comment:类型;"`
	IndicateNum int     `json:"indicate_num" form:"indicate_num" gorm:"column:indicate_num;comment:ab测 随机编号;"`
	City        string  `json:"city" form:"city" gorm:"column:city;comment:城市;"`
	Currency    string  `json:"currency" form:"currency" gorm:"column:currency;comment:货币;"`
	Ip          string  `json:"ip" form:"ip" gorm:"column:ip;comment:注册IP;"`
	LastIp      string  `json:"last_ip" form:"last_ip" gorm:"column:last_ip;comment:最后登录IP;"`
}

func (o *BUser) TableName() string {
	return "b_user"
}

func SaveBUser(record *BUser) error {
	err := global.GVA_USER_DB.Model(BUser{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("save sql SaveBUser error: %s", zap.Error(err))
		return err
	}
	return nil
}

func CreateBUser(record *BUser) error {
	err := global.GVA_USER_DB.Model(BUser{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql CreateBUser error: %s", zap.Error(err))
		return err
	}
	return nil
}

func GetBUserByUsername(username string) (userInfo *BUser, err error) {
	err = global.GVA_USER_DB.Model(BUser{}).
		Where("user_name = ?", username).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetBUserByUsername", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}

func GetBUserByUsernameAndPassword(username, password string) (userInfo *BUser, err error) {
	err = global.GVA_USER_DB.Model(BUser{}).
		Where("user_name = ? and password = ? ", username, password).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetBUserByUsername", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}

func GetAllBUser() (userList []*BUser, err error) {
	err = global.GVA_USER_DB.Model(BUser{}).
		Find(&userList).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return userList, err
	}
	return userList, nil
}

func GetBUserByUid(uid string) (userInfo *BUser, err error) {
	err = global.GVA_USER_DB.Model(BUser{}).
		Where("user_id = ?", uid).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetGameUserByUid", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}

func GetBUserByUUid(uid string) (userInfo *BUser, err error) {
	err = global.GVA_USER_DB.Model(BUser{}).
		Where("uuid = ?", uid).
		First(&userInfo).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error("GetBUserByUUid", zap.Error(err))
		return userInfo, err
	}
	return userInfo, nil
}
