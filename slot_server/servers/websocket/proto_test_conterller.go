package websocket

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"slot_server/lib/common"
	"slot_server/protoc/pbs"
)

func ProtoTestController(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {

	reqData := &pbs.Test1Req{}

	err := proto.Unmarshal(message, reqData)

	fmt.Println(reqData, err)

	ack := &pbs.Test1Ack{
		UserId: "999",
	}

	ackMarshal, _ := proto.Marshal(ack)

	return msgId + 1, common.WebOK, ackMarshal
}

func ProtoTest2Controller(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	reqData := &pbs.Test2Req{}

	err := proto.Unmarshal(message, reqData)

	fmt.Println(reqData, err)

	ack := pbs.Test2Ack{UserId: "1243"}

	ackMarshal, _ := proto.Marshal(&ack)

	return msgId + 1, 23123, ackMarshal
}
