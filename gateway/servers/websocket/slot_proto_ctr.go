package websocket

import (
	"errors"
	"gateway/common"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"gateway/models/table"
	"gateway/protoc/pbs"
	"gateway/utils/cache"
	"github.com/golang/protobuf/proto"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

func Login(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())

	reqData := &pbs.Login{}
	err := proto.Unmarshal(message, reqData)
	if err != nil {
		global.GVA_LOG.Error("proto.Unmarshal error", zap.String("err", err.Error()))
	}

	//服务维护阶段
	//gameServiceConf, err := dao.GetGameServiceConf(1)
	//if err != nil {
	//	code = common.ServerError
	//	global.GVA_LOG.Error("LoginController", zap.Error(err))
	//	return
	//}
	//if gameServiceConf.Maintenance == 1 {
	//	code = common.Maintenance
	//	global.GVA_LOG.Infof("LoginController %v", *gameServiceConf)
	//	return
	//}
	global.GVA_LOG.Infof("Login 用户登录成功 请求业务:%v", *client)

	//jwt 验证
	parseJWT, err := helper.ParseJWT(client.Token)
	global.GVA_LOG.Infof(" ParseJWT %v ", parseJWT)
	if err != nil {
		global.GVA_LOG.Infof(" ParseJWT err ")
		return int32(pbs.ProtocNum_LoginAck), uint32(pbs.ErrCode_TokenExpiration), []byte{}
	}
	userUuId := parseJWT["user_id"].(string)
	if len(userUuId) == 0 {
		return int32(pbs.ProtocNum_LoginAck), uint32(pbs.ErrCode_NotLogin), []byte{}
	}

	if client.IsLogin() {
		global.GVA_LOG.Infof("Login 用户登录 用户已经登录 %v , %v", client.AppID, client.UserID)
		loginAck := pbs.LoginAck{
			UserName: client.Nickname,
		}
		ackMarshal, _ := proto.Marshal(&loginAck)
		return int32(pbs.ProtocNum_LoginAck), uint32(pbs.Code_OK), ackMarshal
	}

	userInfo, err := table.GetBUserByUUid(userUuId)
	if err != nil {
		return int32(pbs.ProtocNum_LoginAck), uint32(pbs.ErrCode_ServerError), []byte{}
	}
	if userInfo.ID <= 0 {
		return int32(pbs.ProtocNum_LoginAck), uint32(pbs.ErrCode_NotRegister), []byte{}
	}
	client.Login(uint32(reqData.AppId), userInfo.Uuid, currentTime, userInfo.UserName, client.Token)

	global.GVA_LOG.Infof("用户登录成功 client:%v", *client)
	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, uint32(reqData.AppId), userInfo.Uuid, client.Addr, currentTime, userInfo.UserName)
	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerRedisError
		global.GVA_LOG.Error("用户登录 SetUserOnlineInfo", zap.Any("", err))
		return
	}

	// 用户登录
	login := &login{
		AppID:  uint32(reqData.AppId),
		UserID: userInfo.Uuid,
		Client: client,
	}
	clientManager.Login <- login

	global.GVA_LOG.Infof("LoginController 用户登录成功 login:%v", *login)
	loginAck := pbs.LoginAck{
		UserName: "",
		City:     "",
		Amount:   0,
	}
	ackMarshal, _ := proto.Marshal(&loginAck)
	global.GVA_LOG.Infof("LoginController 用户登录成功, %v ,%v", client.Addr, string(ackMarshal))
	return int32(pbs.ProtocNum_LoginAck), uint32(pbs.Code_OK), ackMarshal
}

// Heartbeat 心跳接口
func Heartbeat(client *Client, msgId int32, message []byte) (respMsgId int32, code uint32, data []byte) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())

	reqData := &pbs.HeartBeat{}
	err := proto.Unmarshal(message, reqData)
	if err != nil {
		global.GVA_LOG.Error("proto.Unmarshal error", zap.String("err", err.Error()))
	}

	global.GVA_LOG.Infof("心跳接口 webSocket_request client.AppID %v, client.UserID %v", client.AppID, client.UserID)
	if !client.IsLogin() {
		global.GVA_LOG.Error("心跳接口 用户未登录", zap.Any("AppID", client.AppID), zap.Any("UserID", client.UserID))
		code = common.NotLoggedIn
		return
	}

	userOnline, err := cache.GetUserOnlineInfo(client.GetKey())
	if err != nil {
		if errors.Is(err, redis.Nil) {
			global.GVA_LOG.Error("心跳接口 用户未登录", zap.Any(client.UserID, reqData))
			return int32(pbs.ProtocNum_HeartAck), uint32(pbs.ErrCode_NotLogin), []byte{}
		} else {
			global.GVA_LOG.Error("心跳接口 GetUserOnlineInfo", zap.Any("seq", reqData), zap.Any(client.UserID, err))
			return int32(pbs.ProtocNum_HeartAck), uint32(pbs.ErrCode_ServerError), []byte{}
		}
	}

	//if client.AppID == common.AppId10 {
	//	MtHeartReq(client, "", message)
	//}

	global.GVA_LOG.Infof("心跳接口:更新前 Addr:%v client.AppID %v, client.UserID %v currentTime:%v HeartbeatTime:%v", client.Addr, client.AppID, client.UserID, currentTime, helper.TimeIntToStr(int64(client.HeartbeatTime)))
	client.Heartbeat(currentTime)
	global.GVA_LOG.Infof("心跳接口:更新后  Addr:%v client.AppID %v, client.UserID %v heartbeatTime:%v", client.Addr, client.AppID, client.UserID, helper.TimeIntToStr(int64(client.HeartbeatTime)))

	userOnline.Heartbeat(currentTime)

	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		global.GVA_LOG.Error("心跳接口 SetUserOnlineInfo", zap.Any(client.UserID, err), zap.Any("", reqData))
		return int32(pbs.ProtocNum_HeartAck), uint32(pbs.ErrCode_ServerError), []byte{}
	}

	return int32(pbs.ProtocNum_HeartAck), uint32(pbs.Code_OK), []byte{}
}
