package table

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slot_server/lib/global"
)

//CREATE TABLE `user_friend` (
//`id` int unsigned NOT NULL AUTO_INCREMENT,
//`type` int NOT NULL DEFAULT '0' COMMENT '0=双向好友',
//`user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户',
//`friend_user_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '用户朋友',
//`createtime` int DEFAULT '0',
//`updatetime` int DEFAULT NULL,
//PRIMARY KEY (`id`),
//KEY `user_id` (`user_id`) USING BTREE,
//KEY `friend_user_id` (`friend_user_id`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户好友表';

type UserFriend struct {
	GVA_MODEL
	Type         int    `json:"type" form:"type" gorm:"column:type;default:0;comment:0=双向好友"`
	UserId       string `json:"user_id" form:"user_id" gorm:"column:user_id;comment:用户;"`
	FriendUserId string `json:"friend_user_id" form:"friend_user_id" gorm:"column:friend_user_id;comment:用户朋友;"`
}

// TableName 房间配置
func (o *UserFriend) TableName() string {
	return "user_friend"
}

func GetUserFriendByUserIdAndFriendId(userId string, friendUserId string) (record *UserFriend, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("user_id = ? and friend_user_id = ? ", userId, friendUserId).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetUserFriendByUserId(userId string) (record []*UserFriend, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("user_id = ?", userId).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func GetUserFriendById(id int) (record *UserFriend, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("id = ? ", id).
		First(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}

func DelUserFriendById(id int) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("id = ? ", id).
		Delete(UserFriend{}).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql DelUserFriendById error: %s", zap.Error(err))
		return err
	}
	return nil
}

func DelUserFriendByUserIdAndFriendUserId(userId, friendUserId string) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("user_id = ? and friend_user_id = ? ", userId, friendUserId).
		Delete(UserFriend{}).
		Error
	if err != nil {
		global.GVA_LOG.Error("sql DelUserFriendByUserIdAndFriendUserId error: %s", zap.Error(err))
		return err
	}
	return nil
}

func CreateUserFriend(record *UserFriend) error {
	err := global.GVA_SLOT_SERVER_DB.Model(UserFriend{}).
		Create(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserFriend error: %s", zap.Error(err))
		return err
	}
	return nil
}

func SaveUserFriend(record *UserFriend) error {
	err := global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("id = ?", record.ID).
		Save(&record).
		Error
	if err != nil {
		global.GVA_LOG.Error("insert sql UserFriend error: %s", zap.Error(err))
		return err
	}
	return nil
}

// UserFriendPageIsNext
func UserFriendPageIsNext(lastId int) (count int64, err error) {
	db := global.GVA_SLOT_SERVER_DB.Model(UserFriend{})

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

func GetUserFriendByLastId(userId string, lastId int) (record []*UserFriend, err error) {
	err = global.GVA_SLOT_SERVER_DB.
		Model(UserFriend{}).
		Where("user_id = ? and id > ? ", userId, lastId).
		Limit(10).
		Find(&record).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GVA_LOG.Error(err.Error())
		return record, err
	}
	return record, nil
}
