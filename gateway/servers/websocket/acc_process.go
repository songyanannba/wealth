// Package websocket 处理
package websocket

import (
	"encoding/json"
	"gateway/common"
	"gateway/config"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"gateway/protoc/pbs"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"sync"
)

// DisposeFunc 处理函数
type DisposeFunc func(client *Client, seq string, message []byte) (code uint32, msg string, data interface{})

type DisposeProtoFunc func(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte)
type DisposeProtoResp func(msgId int32, message []byte) (respMsgId uint32, code uint32, data interface{})

var (
	handlers        = make(map[string]DisposeFunc)
	handlersRWMutex sync.RWMutex

	ProtoHandlers        = make(map[int32]DisposeProtoFunc)
	ProtoHandlersRWMutex sync.RWMutex

	NatsProtoHandlers        = make(map[int32]DisposeProtoResp)
	NatsProtoHandlersRWMutex sync.RWMutex
)

func RegisterNatsProtoResp(key int32, value DisposeProtoResp) {
	NatsProtoHandlersRWMutex.Lock()
	defer NatsProtoHandlersRWMutex.Unlock()
	NatsProtoHandlers[key] = value
	return
}

func getNatsProtoResp(key int32) (value DisposeProtoResp, ok bool) {
	NatsProtoHandlersRWMutex.RLock()
	defer NatsProtoHandlersRWMutex.RUnlock()
	value, ok = NatsProtoHandlers[key]
	return
}

// Register 注册
func Register(key string, value DisposeFunc) {
	handlersRWMutex.Lock()
	defer handlersRWMutex.Unlock()
	handlers[key] = value
	return
}

func getHandlers(key string) (value DisposeFunc, ok bool) {
	handlersRWMutex.RLock()
	defer handlersRWMutex.RUnlock()
	value, ok = handlers[key]
	return
}

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

// ProcessData 处理数据
//func ProcessData(client *Client, message []byte) {
//	global.GVA_LOG.Infof("处理数据 %v:%v", client.Addr, string(message))
//	defer func() {
//		if r := recover(); r != nil {
//			global.GVA_LOG.Error("处理数据 stop", zap.Any("", r))
//		}
//	}()
//	request := &models.Request{}
//	if err := json.Unmarshal(message, request); err != nil {
//		global.GVA_LOG.Error("处理数据 json Unmarshal", zap.Any("err", err))
//		client.SendMsg([]byte("数据不合法"))
//		return
//	}
//	requestData, err := json.Marshal(request.Data)
//	if err != nil {
//		global.GVA_LOG.Error("处理数据 json Marshal", zap.Any("err", err))
//		client.SendMsg([]byte("处理数据失败"))
//		return
//	}
//	seq := request.Seq
//	cmd := request.Cmd
//	var (
//		code uint32
//		msg  string
//		data interface{}
//	)
//
//	// request
//	global.GVA_LOG.Infof("gate_way_response %v:%v", cmd, client.Addr)
//
//	// 采用 map 注册的方式
//	if value, ok := getHandlers(cmd); ok {
//		code, msg, data = value(client, seq, requestData)
//	} else {
//		code = common.RoutingNotExist
//		global.GVA_LOG.Error("处理数据 路由不存在", zap.Any("Addr", client.Addr), zap.Any("cmd", cmd))
//	}
//	msg = common.GetErrorMessage(code, msg)
//	responseHead := models.NewResponseHead(seq, cmd, code, msg, data)
//	headByte, err := json.Marshal(responseHead)
//	if err != nil {
//		global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
//		return
//	}
//	client.SendMsg(headByte)
//	global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", client.Addr, client.AppID, client.UserID, cmd, code)
//	return
//}

func ProcessDataNew(client *Client, message []byte) {
	global.GVA_LOG.Infof("处理数据 %v:%v", client.Addr, string(message))
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("处理数据 stop", zap.Any("", r))
		}
	}()

	if client.ProtocType == 1 {
		//先解析proto
		netMessage := &pbs.NetMessage{}
		err := proto.Unmarshal(message, netMessage)
		if err == nil && netMessage.MsgId > 0 {
			//优先proto协议
			global.GVA_LOG.Infof("ProcessDataNew 请求 proto协议接口 msgId:%v:%v", netMessage.MsgId, client.Addr)
			// 采用 map 注册的方式
			if value, ok := getHandlersProto(netMessage.MsgId); ok {
				//if netMessage.MsgId == int32(pbs.ProtocNum_LoginReq) {
				//	client.Token = netMessage.ReqHead.Token
				//}
				ackMsgId, ackCode, contentByte := value(client, netMessage.MsgId, netMessage.Content)
				ackMsg := common.GetErrorMessage(ackCode, "")
				//client.SendMsg(headByte)
				netMessageResp := &pbs.NetMessage{
					ReqHead: netMessage.ReqHead,
					AckHead: &pbs.AckHead{
						Uid:     netMessage.ReqHead.Uid,
						Code:    pbs.Code(ackCode),
						Message: ackMsg,
					},
					ServiceId: netMessage.ServiceId,
					MsgId:     ackMsgId,
					Content:   contentByte,
				}
				netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
				global.GVA_LOG.Infof("gate_way_response send headByte:%v ", string(contentByte))
				client.SendMsg(netMessageRespMarshal)
				return
			} else {
				global.GVA_LOG.Error("proto协议接口,处理数据,路由不存在", zap.Any("Addr", client.Addr), zap.Any("MsgId", netMessage.MsgId))
			}
			return
		}
	} else {
		request := &models.Request{}
		err := json.Unmarshal(message, request)
		if err != nil {
			global.GVA_LOG.Error("处理数据 json Unmarshal", zap.Any("err", err))
			client.SendMsg([]byte("数据不合法"))
			return
		}
		requestData, err := json.Marshal(request.Data)
		if err != nil {
			global.GVA_LOG.Error("处理数据 json Marshal", zap.Any("err", err))
			client.SendMsg([]byte("处理数据失败"))
			return
		}
		seq := request.Seq
		cmd := request.Cmd
		serviceId := request.ServiceId
		var (
			code uint32
			msg  string
			data interface{}
		)

		if serviceId == config.NatsSlotServer {
			if value, ok := getHandlers(cmd); ok {
				beforeCode := BeforeHandler(client, cmd)
				if beforeCode != common.OK {
					code = uint32(beforeCode)
				} else {
					code, msg, data = value(client, cmd, requestData)
				}
			} else {
				code = common.RoutingNotExist
				global.GVA_LOG.Error("处理数据 路由不存在", zap.Any("Addr", client.Addr), zap.Any("NatsMemeBattle", cmd))
			}
			msg = common.GetErrorMessage(code, msg)
			responseHead := models.NewResponseHead(seq, cmd, code, msg, data)
			headByte, err := json.Marshal(responseHead)
			if err != nil {
				global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
				return
			}
			client.SendMsg(headByte)
			global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", client.Addr, client.AppID, client.UserID, cmd, code)
		} else {
			// request
			global.GVA_LOG.Infof("gate_way_response %v:%v", cmd, client.Addr)
			// 采用 map 注册的方式
			if value, ok := getHandlers(cmd); ok {
				beforeCode := BeforeHandler(client, cmd)
				if beforeCode != common.OK {
					code = uint32(beforeCode)
				} else {
					code, msg, data = value(client, seq, requestData)
				}
			} else {
				code = common.RoutingNotExist
				global.GVA_LOG.Error("处理数据 路由不存在", zap.Any("Addr", client.Addr), zap.Any("cmd", cmd))
			}
			msg = common.GetErrorMessage(code, msg)
			responseHead := models.NewResponseHead(seq, cmd, code, msg, data)
			headByte, err := json.Marshal(responseHead)
			if err != nil {
				global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
				return
			}
			client.SendMsg(headByte)
			global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", client.Addr, client.AppID, client.UserID, cmd, code)
		}

	}
	return
}

func BeforeHandler(client *Client, cmd string) int {
	if cmd != "login" {
		parseJWT, err := helper.ParseJWT(client.Token)
		global.GVA_LOG.Infof("ParseJWT err %v ", parseJWT)
		if err != nil {
			return common.TokenExpiration
		}
	}
	return common.OK
}

//func NatsSend(request *models.Request) {
//	cmdToInt, err := strconv.Atoi(request.Cmd)
//	if err != nil {
//		global.GVA_LOG.Error("处理数据 cmd ", zap.Any("err", err))
//	}
//
//	msgReq := pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      11,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead: &pbs.AckHead{
//			Uid:     11,
//			Code:    0,
//			Message: "",
//		},
//		ServiceId: NatsMemeBattle,
//		MsgId:     int32(cmdToInt),
//		Content:   []byte{11},
//	}
//
//	//组装
//
//	//获取用户ID
//	clientInfo := GetUserClient(common.AppId3, request.UserID)
//
//	//NastManager.Send("meme_battle", &msgReq)
//	NastManager.SendMemeJs(&msgReq)
//}
