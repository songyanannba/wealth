package router

import (
	"slot_server/protoc/pbs"
	"slot_server/servers/websocket"
)

func InitRouters() {
	MemeBattleRouters()

}

func MemeBattleRouters() {
	//初始化用户信息

	//好友列表
	RegisterProto(int32(pbs.Meb_friendUserList), websocket.FriendUserList)

	//好友申请列表
	RegisterProto(int32(pbs.Meb_auditUserList), websocket.AuditUserList)

	//图鉴列表
	RegisterProto(int32(pbs.Meb_handbookList), websocket.HandbookListController)

	//版本列表
	RegisterProto(int32(pbs.Meb_cardVersionList), websocket.CardVersionListController)

	//拆包
	RegisterProto(int32(pbs.Meb_unpackCard), websocket.UnpackCardController)

	//添加朋友
	RegisterProto(int32(pbs.Meb_addFriend), websocket.AddFriendController)

	//删除好友
	RegisterProto(int32(pbs.Meb_delFriend), websocket.DelFriendController)

	//用户资料详情
	RegisterProto(int32(pbs.Meb_userDetail), websocket.UserDetailController)

	//审核朋友
	RegisterProto(int32(pbs.Meb_authFriend), websocket.AuthFriendController)

	//经验个积分
	RegisterProto(int32(pbs.Meb_coinExperience), websocket.CoinExperienceController)

}
