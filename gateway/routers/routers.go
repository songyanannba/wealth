// Package routers 路由
package routers

import (
	"gateway/servers/websocket"
)

// WebsocketInit Websocket 路由
func WebsocketInit() {

	//proto 协议路由
	websocket.RegisterProto(1, websocket.ProtoTestController)
	websocket.RegisterProto(2, websocket.ProtoTest2Controller)
	//websocket.Register("3", websocket.MTStatus)

	SlotRouter.SlotRouterInit()

}
