package service

import (
	"client/protoc/pbs"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func Test1() {
	CommonService.RegisterHandlers(int32(1), func(msg *pbs.NetMessage) {
		//request := &pbs.GameMessage{}
		//request.Do = "sd"
		//request.To = "hhahah"
		//request.Todo = "会面吧"

		//发送

	})
}

func LoginAck() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_LoginAck), func(msg *pbs.NetMessage) {
		fmt.Println("", msg)

		reqData := &pbs.LoginAck{}

		err := proto.Unmarshal(msg.Content, reqData)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("", reqData)

	})
}
