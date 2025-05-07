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

}
