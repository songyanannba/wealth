package routers

import "gateway/servers/websocket"

type memeBattle struct {
}

var MemeBattle = memeBattle{}

func (rm *memeBattle) WayRouterInit() {

	//meme websocket 入口
	websocket.Register("memeBattleEntry", websocket.MemeBattleEntry)

	// 快速匹配房间 后台定时任务自动匹配
	websocket.Register("mebQuickMatchRoom", websocket.QuickMatchRoom)

	//取消快速匹配
	websocket.Register("mebCancelMatchRoom", websocket.MebCancelMatchRoom)

	////获取币余额
	//websocket.Register("getCoinBalance", websocket.TavernCoinBalanceCtl)

	//获取当前游戏状态
	websocket.Register("mebUserState", websocket.GetUserState)

	//房间资费配置
	websocket.Register("mebRoomConfig", websocket.GetRoomConfig)

	//1:创建房间
	websocket.Register("mebCreateRoom", websocket.MebCreateRoom)

	//邀请好友
	websocket.Register("mebInviteFriend", websocket.MebInviteFriend)

	//3:加入房间
	websocket.Register("mebJoinRoom", websocket.MebJoinRoom)

	//4:准备就绪（就绪） 玩家发起
	websocket.Register("mebReady", websocket.MebReady)

	//取消就绪
	websocket.Register("mebCancelReady", websocket.MebCancelReady)

	//重新加入房间 是在用户掉线前已经加入过房间，重新连接后，才能用：也就是mebUserState接口有返回房间数据的时候调用
	websocket.Register("mebReJoinRoom", websocket.MebReJoinRoom)

	//5 离开房间
	websocket.Register("mebLeaveRoom", websocket.MebLeaveRoom)

	//6 房主踢人
	websocket.Register("mebKickRoom", websocket.MebKickRoom)

	//7:房主开始对局游戏
	websocket.Register("mebStartPlay", websocket.MebStartPlay)

	//8:房间心跳 房主创建房间完成之后
	websocket.Register("mebRoomAlive", websocket.MebRoomAlive)

	//9:加载完成 （进入房间后每个用户都要告诉服务器）
	websocket.Register("mebLoadCompleted", websocket.MebLoadCompleted)

	//10:操作牌 （0:看牌 1:出牌 2:表情 3:重随）
	websocket.Register("mebOperateCard", websocket.MebOperateCard)

	//点赞
	websocket.Register("mebCardLike", websocket.MebCardLike)

	// === 下面是 grpc 接口 ===

	//图鉴列表
	websocket.Register("mebHandbookList", websocket.MebHandbookList)

	//版本列表
	websocket.Register("mebCardVersionList", websocket.MebCardVersionList)

	//拆包
	websocket.Register("mebUnpackCard", websocket.MebUnpackCard)

	//添加好友
	websocket.Register("mebAddFriend", websocket.MebAddFriend)

	//删除好友
	websocket.Register("mebDelFriend", websocket.MebDelFriend)

	//查看资料
	websocket.Register("mebUserDetail", websocket.MebUserDetail)

	//我的好友列表
	websocket.Register("mebFriendList", websocket.MebFriendList)

	//申请好友列表
	websocket.Register("mebAuditUserList", websocket.MebAuditUserList)

	//审核同意成为好友
	websocket.Register("mebAuthFriend", websocket.MebAuthFriend)

	//获取用户的经验和积分
	websocket.Register("mebGetCoinExperience", websocket.MebGetCoinExperience)

}
