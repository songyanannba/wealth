// Package routers 路由
package routers

import (
	"gateway/servers/websocket"
)

// WebsocketInit Websocket 路由
func WebsocketInit() {
	//websocket.Register("login", websocket.LoginController)
	//websocket.Register("heartbeat", websocket.HeartbeatController)

	//proto 协议路由
	websocket.RegisterProto(1, websocket.ProtoTestController)
	websocket.RegisterProto(2, websocket.ProtoTest2Controller)

	SlotRouter.SlotRouterInit()

}
