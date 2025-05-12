package router

import (
	"context"
	"slot_server/protoc/pbs"
	"sync"
)

type DisposeProtoFunc func(ctx context.Context, req *pbs.NetMessage) (*pbs.NetMessage, error)

var (
	ProtoHandlers        = make(map[int32]DisposeProtoFunc)
	ProtoHandlersRWMutex sync.RWMutex
)

func RegisterProto(key int32, value DisposeProtoFunc) {
	ProtoHandlersRWMutex.Lock()
	defer ProtoHandlersRWMutex.Unlock()
	ProtoHandlers[key] = value
	return
}

func GetHandlersProto(key int32) (value DisposeProtoFunc, ok bool) {
	ProtoHandlersRWMutex.RLock()
	defer ProtoHandlersRWMutex.RUnlock()
	value, ok = ProtoHandlers[key]
	return value, ok
}
