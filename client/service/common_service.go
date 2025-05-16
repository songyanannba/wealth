package service

import (
	"client/protoc/pbs"
	"sync"
)

type commonService struct {
	syncMutex  *sync.Mutex
	HandlerMap map[int32]func(msg *pbs.NetMessage)
}

var CommonService = commonService{
	syncMutex:  new(sync.Mutex),
	HandlerMap: make(map[int32]func(msg *pbs.NetMessage)),
}

func (cs *commonService) First() {

}

func (cs *commonService) Start() {
	Test1()
	LoginAck()
	CurrAPInfoAck()
	UserBetAck()
	ReceivedAnimalSortMsg()
	ReceivedCurrPeriodUserWinMsg()
	ReceivedColorSortMsg()
}

func (cs *commonService) RegisterHandlers(typeInt int32, f func(msg *pbs.NetMessage)) {
	cs.syncMutex.Lock()
	defer cs.syncMutex.Unlock()

	if _, ok := cs.HandlerMap[typeInt]; !ok {
		cs.HandlerMap[typeInt] = f
	}
}

func (cs *commonService) GetHandlers(msg *pbs.NetMessage) (value func(msg *pbs.NetMessage), ok bool) {
	cs.syncMutex.Lock()
	defer cs.syncMutex.Unlock()
	value, ok = cs.HandlerMap[msg.MsgId]
	return
}
