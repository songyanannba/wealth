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
		reqData := &pbs.LoginAck{}
		err := proto.Unmarshal(msg.Content, reqData)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("LoginAck === ", reqData)
	})
}

func CurrAPInfoAck() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_CurrAPInfoAck), func(msg *pbs.NetMessage) {
		reqData := &pbs.CurrAPInfoAck{}
		err := proto.Unmarshal(msg.Content, reqData)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("CurrAPInfoAck === ", reqData, string(msg.Content))

	})
}

func UserBetAck() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_betAck), func(msg *pbs.NetMessage) {
		reqData := &pbs.UserBetAck{}

		err := proto.Unmarshal(msg.Content, reqData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("UserBetAck === ", reqData)
	})
}

func ReceivedAnimalSortMsg() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_AnimalSortMsg), func(msg *pbs.NetMessage) {
		msgData := &pbs.AnimalSortMsg{}
		err := proto.Unmarshal(msg.Content, msgData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ReceivedAnimalSortMsg === ", msgData, string(msg.Content))
	})
}

func ReceivedCurrPeriodUserWinMsg() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_CurrPeriodUserWinMsg), func(msg *pbs.NetMessage) {
		msgData := &pbs.CurrPeriodUserWinMsg{}
		err := proto.Unmarshal(msg.Content, msgData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ReceivedCurrPeriodUserWinMsg === ", msgData, string(msg.Content))
	})
}

func ReceivedColorSortMsg() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_ColorSortMsg), func(msg *pbs.NetMessage) {
		msgData := &pbs.ColorSortMsg{}
		err := proto.Unmarshal(msg.Content, msgData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("ReceivedColorSortMsg === ", msgData, string(msg.Content))
	})
}

func OnLineUserListAck() {
	CommonService.RegisterHandlers(int32(pbs.ProtocNum_OnLineUserListAck), func(msg *pbs.NetMessage) {
		msgData := &pbs.OnLineUserListAck{}
		err := proto.Unmarshal(msg.Content, msgData)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("OnLineUserListAck === ", msgData, string(msg.Content))
	})
}
