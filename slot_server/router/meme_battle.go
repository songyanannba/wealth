package router

import (
	"slot_server/protoc/pbs"
	"slot_server/servers/websocket"
)

// NatsRouterInit nats 路由；收到网关的消息 - 找到对应的处理方法
func NatsRouterInit() {
	//nats.RegisterProto(3, nats.ProtoMTTest)

	websocket.RegisterProto(int32(pbs.Meb_mtHeartReq), websocket.Heart400)
	//匹配
	websocket.RegisterProto(int32(pbs.Meb_memeMatchRoom), websocket.MatchRoomController)
	//取消匹配
	websocket.RegisterProto(int32(pbs.Meb_cancelMatchRoom), websocket.CancelMatchRoomController)
	//创建房间
	websocket.RegisterProto(int32(pbs.Meb_createRoom), websocket.CreateRoomController)

	websocket.RegisterProto(int32(pbs.Meb_inviteFriend), websocket.InviteFriendController)

	websocket.RegisterProto(int32(pbs.Meb_joinRoom), websocket.JoinRoomRoomController)

	websocket.RegisterProto(int32(pbs.Meb_readyMsg), websocket.ReadyRoomRoomController)

	websocket.RegisterProto(int32(pbs.Meb_cancelReady), websocket.CancelReadyRoomRoomController)

	websocket.RegisterProto(int32(pbs.Meb_leaveRoom), websocket.LeaveRoomController)

	websocket.RegisterProto(int32(pbs.Meb_userState), websocket.UserStateController)

	websocket.RegisterProto(int32(pbs.Meb_reJoinRoom), websocket.ReJoinRoomController)

	websocket.RegisterProto(int32(pbs.Meb_kickRoom), websocket.KickRoomController)

	websocket.RegisterProto(int32(pbs.Meb_startPlay), websocket.StartPlayController)

	websocket.RegisterProto(int32(pbs.Meb_loadCompleted), websocket.LoadCompletedController)

	websocket.RegisterProto(int32(pbs.Meb_operateCards), websocket.OperateCardController)

	websocket.RegisterProto(int32(pbs.Meb_likeCards), websocket.LikeCardsController)

	websocket.RegisterProto(int32(pbs.Meb_roomAlive), websocket.RoomAliveController)

}
