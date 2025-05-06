package test

import (
	"testing"
)

func Test_Marshal(t *testing.T) {
	//req := pbCom.Test1Req{UserId: "1"}
	//reqMarshal, _ := proto.Marshal(&req)
	//
	//netMessage := pbCom.NetMessage{
	//	ReqHead: &pbCom.ReqHead{
	//		Uid:      1,
	//		Token:    "xxx",
	//		Platform: "a",
	//	},
	//	AckHead: &pbCom.AckHead{
	//		Uid:     1,
	//		Code:    0,
	//		Message: "",
	//	},
	//	ServiceId: "1",
	//	MsgId:     1,
	//	Content:   reqMarshal,
	//}
	//netMessageMarshal, _ := proto.Marshal(&netMessage)
	//
	//fmt.Println("netMessageMarshal:", netMessageMarshal)

	websocket.WsClientService.Start()

}
