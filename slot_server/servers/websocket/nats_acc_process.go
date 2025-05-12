package websocket

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/protoc/pbs"
	"sync"
)

// DisposeFunc 处理函数

//type DisposeProtoFunc func(msgId int32, message []byte) (respMsgId int32, code uint32, data []byte)

type DisposeProtoFunc func(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte)

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

func getHandlersProto(key int32) (value DisposeProtoFunc, ok bool) {
	ProtoHandlersRWMutex.RLock()
	defer ProtoHandlersRWMutex.RUnlock()
	value, ok = ProtoHandlers[key]
	return
}

func ProcessData(netMessage *pbs.NetMessage) {
	global.GVA_LOG.Infof("ProcessData 处理数据 :%v", netMessage)
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("ProcessData 处理数据 stop", zap.Any("", r))
		}
	}()
	if netMessage.MsgId <= 0 {
		global.GVA_LOG.Error("ProcessData 协议号错误 不存在")
	}
	//优先proto协议
	global.GVA_LOG.Infof("ProcessDataNew 请求 proto协议接口 msgId:%v", netMessage.MsgId)
	// 采用 map 注册的方式
	if value, ok := getHandlersProto(netMessage.MsgId); ok {
		ackMsgId, ackCode, contentByte := value(netMessage)
		//ackMsgId, ackCode, contentByte := value(netMessage.MsgId, netMessage.Content)
		global.GVA_LOG.Infof("ackMsgId %v, ackCode %v, contentByte %v", ackMsgId, ackCode, contentByte)
		//ackMsg := common.GetErrorMessage(ackCode, "")
		//netMessageResp := &pbs.NetMessage{
		//	ReqHead: netMessage.ReqHead,
		//	AckHead: &pbs.AckHead{
		//		Uid:     netMessage.ReqHead.Uid,
		//		Code:    pbs.Code(ackCode),
		//		Message: ackMsg,
		//	},
		//	ServiceId: netMessage.ServiceId,
		//	MsgId:     ackMsgId,
		//	Content:   contentByte,
		//}
		//netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		//global.GVA_LOG.Infof("magic_tower send headByte:%v ", string(contentByte))
		//NastManager.Producer(netMessageRespMarshal)
		return
	} else {
		global.GVA_LOG.Error("proto协议接口,处理数据,路由不存在", zap.Any("MsgId", netMessage.MsgId))
	}
	return
}
