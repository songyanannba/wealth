package websocket

import (
	"gateway/global"
	"gateway/protoc/pbs"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func CurrAPInfo(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	//解析参数
	//client.UserID = "b3703c51-1238-4c7c-ae71-9d73fd4bf1f1"
	global.GVA_LOG.Infof("CurrAPInfo %v", client.UserID)

	reqProto := &pbs.NatsCurrAPInfo{
		UserId: client.UserID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, client.UserID, int32(pbs.ProtocNum_CurrAPInfoReq))

	return int32(pbs.ProtocNum_CurrAPInfoReq), uint32(pbs.Code_OK), []byte{}
}

func OnLineUserList(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	//解析参数
	global.GVA_LOG.Infof("OnLineUserList %v", client.UserID)

	reqProto := &pbs.OnLineUserListAck{}

	for _, uCli := range clientManager.GetUserClients() {
		reqProto.OnlineUser = append(reqProto.OnlineUser, &pbs.OnlineUser{
			UserId:   uCli.UserID,
			UserName: uCli.Nickname,
		})
		if len(reqProto.OnlineUser) == 20 {
			break
		}
	}

	ackMarshal, _ := proto.Marshal(reqProto)
	return int32(pbs.ProtocNum_OnLineUserListAck), uint32(pbs.Code_OK), ackMarshal
}

func UserBetReq(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	//解析参数
	global.GVA_LOG.Infof("UserBetReq %v", client.UserID)

	reqProto := &pbs.UserBetReq{}
	err := proto.Unmarshal(message, reqProto)
	if err != nil {
		global.GVA_LOG.Error("UserBetReq", zap.Error(err))
	}

	global.GVA_LOG.Infof("UserBetReq %v", reqProto)

	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, client.UserID, int32(pbs.ProtocNum_betReq))

	return int32(pbs.ProtocNum_betReq), uint32(pbs.Code_OK), []byte{}
}
