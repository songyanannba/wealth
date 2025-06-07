package router

import (
	"slot_server/protoc/pbs"
	"slot_server/servers/websocket"
)

func NatsSlotInit() {
	//nats.RegisterProto(3, nats.ProtoMTTest)

	websocket.RegisterProto(int32(pbs.ProtocNum_CurrAPInfoReq), websocket.CurrAPInfos)

	websocket.RegisterProto(int32(pbs.ProtocNum_betReq), websocket.UserBetReq)

}
