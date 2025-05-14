package websocket

import (
	"gateway/global"
	"gateway/protoc/pbs"
	"github.com/golang/protobuf/proto"
)

func CurrAPInfo(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	//解析参数

	client.UserID = "b3703c51-1238-4c7c-ae71-9d73fd4bf1f1"

	global.GVA_LOG.Infof("CurrAPInfo %v", client.UserID)

	reqProto := &pbs.NatsCurrAPInfo{
		UserId: client.UserID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, client.UserID, int32(pbs.ProtocNum_CurrAPInfoReq))

	return int32(pbs.ProtocNum_CurrAPInfoReq), uint32(pbs.Code_OK), []byte{}
}
