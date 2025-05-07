package websocket

import (
	"fmt"
	"gateway/common"
	"gateway/protoc/pbs"

	"github.com/golang/protobuf/proto"
)

func ProtoTestController(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	reqData := &pbs.Test1Req{}
	err := proto.Unmarshal(message, reqData)
	fmt.Println(reqData, err)
	ack := &pbs.Test1Ack{
		UserId: "999",
	}

	ackMarshal, _ := proto.Marshal(ack)

	//uID, _ := strconv.Atoi(client.UserID)
	//netMessageResp := &pbs.NetMessage{
	//	AckHead: &pbs.AckHead{
	//		Uid:     int32(uID),
	//		Code:    common.WebOK,
	//		Message: ackMsg,
	//	},
	//	ServiceId: "",
	//	MsgId:      msgId + 1,
	//	Content:   ackMarshal,
	//}

	//ackMarshal, _ := proto.Marshal(ack)
	return msgId + 1, common.OK, ackMarshal
}

func ProtoTest2Controller(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	reqData := &pbs.Test2Req{}

	err := proto.Unmarshal(message, reqData)

	fmt.Println(reqData, err)

	ack := pbs.Test2Ack{UserId: "1243"}

	ackMarshal, _ := proto.Marshal(&ack)

	return msgId + 1, 23123, ackMarshal
}
