package meme_serve

import (
	"sync"
)

type CommonRoomManager struct {
	Sync      *sync.RWMutex //读写锁
	MatchLock *sync.Mutex   // 匹配锁
	CloseRoom chan []byte
	Broadcast chan []byte //广播类型的消息 消息中需要有房间号
}

func GetCommonRoomManager() *CommonRoomManager {
	comRoomMgr := &CommonRoomManager{
		Sync:      new(sync.RWMutex),
		MatchLock: new(sync.Mutex),
		Broadcast: make(chan []byte),
		CloseRoom: make(chan []byte),
	}
	return comRoomMgr
}
