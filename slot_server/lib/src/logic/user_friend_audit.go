package logic

import (
	"go.uber.org/zap"
	"slot_server/lib/common"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
)

func AuditUserAndIsHaveNextPage(userId string, lastId int) ([]*models.AuditUser, bool) {
	isHaveNextPage := false
	pageCount, err := table.UserFriendAuditPageIsNext(lastId)
	if err != nil {
		global.GVA_LOG.Error("AuditUserAndIsHaveNextPage:", zap.Error(err))
	}
	if pageCount > 10 {
		isHaveNextPage = true
	}

	return AuditUserList(userId, lastId), isHaveNextPage
}

func AuditUserList(userId string, lastId int) []*models.AuditUser {
	//todo 加缓存
	auditUser := make([]*models.AuditUser, 0)

	configs, err := table.GetUserFriendAuditByLastId(userId, lastId)
	if err != nil {
		global.GVA_LOG.Error("UserFriendList GetUserFriendByLastId", zap.Error(err))
		return auditUser
	}

	for _, config := range configs {
		//获取好友昵称 todo
		auditUser = append(auditUser, &models.AuditUser{
			ApplicationUser: config.ApplicationUser,
			Nickname:        config.ApplicationUser + "_" + "敬请期待",
			AuditId:         config.ID,
		})
	}

	return auditUser
}

func UserFriendAuditByAuditAndApplication(auditUser, applicationUser string) int {

	application, err := table.UserFriendAuditByAuditAndApplication(auditUser, applicationUser)

	if err != nil {
		global.GVA_LOG.Error("UserFriendAuditByAuditAndApplication:", zap.Error(err))
	}

	//没有的申请
	if application.ID > 0 {
		//审核中
		if application.IsAgree == 0 {
			return common.AuditIng
		}
		//审核通过 并且有好友关系
		if application.IsAgree == 1 {
			record, err := table.GetUserFriendByUserIdAndFriendId(application.AuditUser, application.ApplicationUser)
			if err != nil {
				global.GVA_LOG.Error("UserFriendAuditByAuditAndApplication CreateUserFriend:", zap.Error(err))
			}
			if record.ID > 0 {
				return common.HaveFriend
			}
		}
		//虽然通过 但是没有好友关系 说明已经删除
	}
	return common.OK
}

func CreateUserFriendAudit(auditUser, applicationUser string) {

	err := table.CreateUserFriendAudit(&table.UserFriendAudit{
		ApplicationUser: applicationUser,
		AuditUser:       auditUser,
		IsAgree:         0,
		DateTime:        helper.LocalTime(),
	})
	if err != nil {
		global.GVA_LOG.Error("CreateUserFriendAudit:", zap.Error(err))
	}
}
