// Package routers 路由
package routers

import (
	"gateway/servers/websocket"
)

// WebsocketInit Websocket 路由
func WebsocketInit() {
	websocket.Register("login", websocket.LoginController)
	websocket.Register("heartbeat", websocket.HeartbeatController)
	//websocket.Register("ping", websocket.PingController)

	//proto 协议路由
	websocket.RegisterProto(1, websocket.ProtoTestController)
	websocket.RegisterProto(2, websocket.ProtoTest2Controller)
	//websocket.Register("3", websocket.MTStatus)

	//meme nats消息返回处理方法
	memeBattleRouters.WayRouterInit()

	//meme websocket 入口
	//MemeBattle.WayRouterInit()

	SlotRouter.SlotRouterInit()

}
