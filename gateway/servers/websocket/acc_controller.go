// Package websocket 处理
package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"gateway/common"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"gateway/utils/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
)

// PingController ping
func PingController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	fmt.Println("webSocket_request ping接口", client.Addr, seq, message)
	data = "pong"
	return
}

// LoginController 用户登录
func LoginController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.Login{}

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

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("LoginController 用户登录 解析数据失败", zap.Any("", seq), zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("LoginController, {%v} seq:%v", *request, seq)

	//environment := global.GVA_VP.GetString("app.environment")
	//if environment != "dev" {
	//	userInfo, err := logic.GetUserInfo(request.Token)
	//	if err != nil || userInfo.Data.Id <= 0 {
	//		code = common.Unauthorized
	//		global.GVA_LOG.Error("LoginController GetUserInfo %v", zap.Any("userInfo", userInfo))
	//		global.GVA_LOG.Infof("LoginController GetUserInfo %v", zap.Any("userInfo", userInfo))
	//		return
	//	}
	//
	//	global.GVA_LOG.Infof("LoginController GetUserInfo %v", userInfo)
	//	request.UserID = strconv.Itoa(userInfo.Data.Id)
	//	request.Nickname = userInfo.Data.Nickname
	//	if request.UserID == "" || len(request.UserID) >= 20 {
	//		code = common.UnauthorizedUserID
	//		global.GVA_LOG.Infof("LoginController 用户登录 非法的用户 %v , {%v}", seq, request.UserID)
	//		return
	//	}
	//
	//	if userInfo.Data.IsAuthentication != 1 {
	//		global.GVA_LOG.Infof("LoginController 用户登录 非法的用户 未实名 %v , {%v}", seq, request.UserID)
	//		code = common.NotAuthentication
	//		return
	//	}
	//} else {
	//	if len(request.UserID) == 0 {
	//		userInfo, err := logic.GetUserInfo(request.Token)
	//		if err != nil || userInfo.Data.Id <= 0 {
	//			code = common.Unauthorized
	//			global.GVA_LOG.Error("LoginController GetUserInfo %v", zap.Any("userInfo", userInfo))
	//			global.GVA_LOG.Infof("LoginController GetUserInfo %v", zap.Any("userInfo", userInfo))
	//			return
	//		}
	//
	//		global.GVA_LOG.Infof("LoginController GetUserInfo %v", userInfo)
	//		request.UserID = strconv.Itoa(userInfo.Data.Id)
	//		request.Nickname = userInfo.Data.Nickname
	//		if request.UserID == "" || len(request.UserID) >= 20 {
	//			code = common.UnauthorizedUserID
	//			global.GVA_LOG.Infof("LoginController 用户登录 非法的用户 %v , {%v}", seq, request.UserID)
	//			return
	//		}
	//
	//		if userInfo.Data.IsAuthentication != 1 {
	//			global.GVA_LOG.Infof("LoginController 用户登录 非法的用户 未实名 %v , {%v}", seq, request.UserID)
	//			code = common.NotAuthentication
	//			return
	//		}
	//	}
	//}

	global.GVA_LOG.Infof("LoginController 用户登录成功 请求业务:%v", *client)

	//默认是1
	if request.AppID == 0 {
		request.AppID = common.AppId1 //相当于游戏项目ID
	}
	if !InAppIDs(request.AppID) {
		code = common.Unauthorized
		global.GVA_LOG.Infof("LoginController 用户登录 不支持的平台 %v , %v", seq, request.UserID)
		return
	}

	if client.IsLogin() {
		global.GVA_LOG.Infof("LoginController 用户登录 用户已经登录 %v , %v,%v", client.AppID, client.UserID, seq)
		code = common.ReLogin
		data = models.LogicRespAck{
			ServiceToken: client.Token,
			UserID:       request.UserID,
			Nickname:     request.Nickname,
		}
		return
	}

	oldClient := GetUserClient(common.AppId3, request.UserID)
	if oldClient != nil && oldClient.UserID == request.UserID {
		//发送广播 顶号
		msgData := models.RepLoginMsg{
			ProtoNum:  models.RepLogin,
			Timestamp: time.Now().Unix(),
		}
		responseHead := models.NewResponseHead("", models.RepLogin, common.OK, "", msgData)
		responseHeadByte, _ := json.Marshal(responseHead)
		oldClient.SendMsg(responseHeadByte)
	}

	token, _ := helper.GenerateJWT(request.UserID, "", 48)
	client.Login(request.AppID, request.UserID, currentTime, request.Nickname, token)

	global.GVA_LOG.Infof("LoginController 用户登录成功 client:%v", *client)
	// 存储数据
	userOnline := models.UserLogin(serverIp, serverPort, request.AppID, request.UserID, client.Addr, currentTime, request.Nickname)
	err := cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerRedisError
		global.GVA_LOG.Error("LoginController 用户登录 SetUserOnlineInfo", zap.Any(seq, err))
		return
	}

	idempotent, err := cache.GetIdempotent(request.UserID, "Login")
	if err != nil {
		code = common.ServerError
		return
	}
	if len(idempotent) > 0 {
		code = common.DuplicateRequests
		return
	}
	err = cache.SetIdempotentNx(request.UserID, "Login", "Login")
	if err != nil {
		code = common.ServerError
		return
	}

	//保存用户信息
	//err = dao.SaveGameUser(&table.GameUser{
	//	UserId:   request.UserID,
	//	Nickname: request.Nickname,
	//	KingCoin: KingCoin,
	//	Token:    "",
	//})

	// 用户登录
	login := &login{
		AppID:  request.AppID,
		UserID: request.UserID,
		Client: client,
	}
	clientManager.Login <- login

	//骗子酒馆
	//if request.AppID == common.AppId3 {
	//	dao.InitTavernUsersRoom(request.UserID)
	//	logic.InitTavernUserCharacter(request.UserID)
	//}
	//
	////幽影魔塔
	//if request.AppID == common.AppId4 {
	//	logic.MTInitGameInfo(request.UserID, request.Nickname)
	//}

	global.GVA_LOG.Infof("LoginController 用户登录成功 login:%v", *login)
	data = models.LogicRespAck{
		ServiceToken: token,
		UserID:       request.UserID,
		Nickname:     request.Nickname,
	}

	global.GVA_LOG.Infof("LoginController 用户登录成功 %v , %v ,%v", seq, client.Addr, request.UserID)
	return
}

// HeartbeatController 心跳接口
func HeartbeatController(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	currentTime := uint64(time.Now().Unix())
	request := &models.HeartBeat{}
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("心跳接口 解析数据失败", zap.Any(seq, err))
		return
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
			code = common.NotLoggedIn
			global.GVA_LOG.Error("心跳接口 用户未登录", zap.Any(client.UserID, seq))
			return
		} else {
			code = common.ServerError
			global.GVA_LOG.Error("心跳接口 GetUserOnlineInfo", zap.Any("seq", seq), zap.Any(client.UserID, err))
			return
		}
	}

	if client.AppID == common.AppId10 {
		MtHeartReq(client, "", message)
	}

	global.GVA_LOG.Infof("心跳接口:更新前 Addr:%v client.AppID %v, client.UserID %v currentTime:%v HeartbeatTime:%v", client.Addr, client.AppID, client.UserID, currentTime, helper.TimeIntToStr(int64(client.HeartbeatTime)))
	client.Heartbeat(currentTime)
	global.GVA_LOG.Infof("心跳接口:更新后  Addr:%v client.AppID %v, client.UserID %v heartbeatTime:%v", client.Addr, client.AppID, client.UserID, helper.TimeIntToStr(int64(client.HeartbeatTime)))

	//userOnline.Heartbeat(currentTime)

	err = cache.SetUserOnlineInfo(client.GetKey(), userOnline)
	if err != nil {
		code = common.ServerError
		global.GVA_LOG.Error("心跳接口 SetUserOnlineInfo", zap.Any(client.UserID, err), zap.Any("", seq))
		return
	}
	return
}
