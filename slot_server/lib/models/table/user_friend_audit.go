package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
	"time"
)

//CREATE TABLE `user_friend_audit` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`application_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '申请人',
//`audit_user` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '审核人，被添加好友的人',
//`is_agree` int NOT NULL DEFAULT '0' COMMENT '0=未同意 1=同意',
//`date_time` datetime DEFAULT NULL,
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT '0',
//PRIMARY KEY (`id`),
//KEY `audit_user` (`audit_user`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户好友申请表';

type UserFriendAudit struct {
	GVA_MODEL
	ApplicationUser string     `json:"application_user" form:"application_user" gorm:"column:application_user;comment:申请人;"`
	AuditUser       string     `json:"audit_user" form:"audit_user" gorm:"column:audit_user;comment:审核人，被添加好友的人 主动输入的用户ID;"`
	IsAgree         int        `json:"is_agree" form:"is_agree" gorm:"column:is_agree;default:0;comment:0=未同意 1=同意"`
	DateTime        *time.Time `json:"date_time" form:"date_time" gorm:"column:date_time;comment:时间;"`
}

// TableName 房间配置
func (o *UserFriendAudit) TableName() string {
	return "user_friend_audit"
}

func UserFriendAuditById(aId int) (record *UserFriendAudit, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriendAudit{}).
		Where("id = ?", aId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func UserFriendAuditByAuditAndApplication(auditUser, applicationUser string) (record *UserFriendAudit, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriendAudit{}).
		Where("audit_user = ? and application_user = ? ", auditUser, applicationUser).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

// UserFriendAuditPageIsNext
func UserFriendAuditPageIsNext(lastId int) (count int64, err error) {
	db := global.GVA_SLOT_SERVER_DB.Model(UserFriendAudit{})
	if lastId != 0 {
		db.Where(" id > ?", lastId)
	}
	err = db.Count(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return count, err
	}
	return count, nil
}

func GetUserFriendAuditByLastId(userId string, lastId int) (record []*UserFriendAudit, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriendAudit{}).
		Where("audit_user = ? and id > ? ", userId, lastId).
		Where("is_agree = 0").
		Limit(10).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetUserFriendAuditByUserId(auditUser string) (record []*UserFriendAudit, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriendAudit{}).
		Where("audit_user = ?", auditUser).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func CreateUserFriendAudit(record *UserFriendAudit) error {
	err := global.GVA_SLOT_SERVER_DB.Model(UserFriendAudit{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserFriendAudit error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveUserFriendAudit(record *UserFriendAudit) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserFriendAudit{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserFriendAudit error: %s", zap.Error(err))
		return err
	}
	return nil
}
