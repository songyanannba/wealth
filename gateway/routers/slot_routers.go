package routers

import (
	"gateway/protoc/pbs"
	"gateway/servers/websocket"
)

type slotRouter struct {
}

var SlotRouter = slotRouter{}

func (rm *slotRouter) SlotRouterInit() {
	//登录
	websocket.RegisterProto(int32(pbs.ProtocNum_LoginReq), websocket.Login)

	//心跳
	websocket.RegisterProto(int32(pbs.ProtocNum_HeartReq), websocket.Heartbeat)

	//动物派对全局的房间信息
	websocket.RegisterProto(int32(pbs.ProtocNum_CurrAPInfoReq), websocket.CurrAPInfo)

	//获取所有的在线用户列表
	websocket.RegisterProto(int32(pbs.ProtocNum_OnLineUserListReq), websocket.OnLineUserList)

	//押注
	websocket.RegisterProto(int32(pbs.ProtocNum_betReq), websocket.UserBetReq)
}
