package logic

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
)

func UserFriendAndIsHaveNextPage(userId string, lastId int) ([]*models.UserFriend, bool) {
	isHaveNextPage := false
	pageCount, err := table.UserFriendPageIsNext(lastId)
	if err != nil {
		global.GVA_LOG.Error("UserFriendAndIsHaveNextPage:", zap.Error(err))
	}
	if pageCount > 10 {
		isHaveNextPage = true
	}

	return UserFriendList(userId, lastId), isHaveNextPage
}

func UserFriendList(userId string, lastId int) []*models.UserFriend {
	//todo 加缓存
	userFriend := make([]*models.UserFriend, 0)

	configs, err := table.GetUserFriendByLastId(userId, lastId)
	if err != nil {
		global.GVA_LOG.Error("UserFriendList GetUserFriendByLastId", zap.Error(err))
		return userFriend
	}

	for _, config := range configs {

		//获取好友昵称 todo
		userFriend = append(userFriend, &models.UserFriend{
			FriendUserId: config.FriendUserId,
			Nickname:     config.FriendUserId + "_" + "敬请期待",
			FriendId:     config.ID,
		})
	}

	return userFriend
}

func CreateUserFriend(auditUser, applicationUser string) {
	err := table.CreateUserFriend(&table.UserFriend{
		Type:         0,
		UserId:       auditUser,
		FriendUserId: applicationUser,
	})
	if err != nil {
		global.GVA_LOG.Error("CreateUserFriend:", zap.Error(err))
	}
}
