package service

import (
	"client/protoc/pbs"
	"fmt"
	"sync"
)

type cliHandler struct {
	syncMutex  *sync.Mutex
	HandlerMap map[int32]func(msg *pbs.NetMessage)
}

var CliHandler = &cliHandler{
	syncMutex:  new(sync.Mutex),
	HandlerMap: make(map[int32]func(msg *pbs.NetMessage)),
}

func (ch *cliHandler) Start() {
	//LoginAck()
}

func (ch *cliHandler) DaYin(msg *pbs.NetMessage) {

	fmt.Println("dayin msg == ", msg)

}

func (cs *cliHandler) RegisterHandlers(typeInt int32, f func(msg *pbs.NetMessage)) {
	cs.syncMutex.Lock()
	defer cs.syncMutex.Unlock()

	if _, ok := cs.HandlerMap[typeInt]; !ok {
		cs.HandlerMap[typeInt] = f
	}
}

//func LoginAck() {
//	CommonService.RegisterHandlers(int32(pbs.ProtocNum_LoginAck), func(msg *pbs.NetMessage) {
//		fmt.Println("", msg)
//
//		reqData := &pbs.LoginAck{}
//
//		err := proto.Unmarshal(msg.Content, reqData)
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println("", reqData)
//
//	})
//}
