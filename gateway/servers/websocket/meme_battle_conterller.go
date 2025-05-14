package websocket

import (
	"context"
	"encoding/json"
	"gateway/common"
	"gateway/config"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"gateway/protoc/pbs"
	"gateway/servers/grpcclient"
	"gateway/utils/cache"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"time"
)

func MtHeartReq(client *Client, cmd string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		userID = client.UserID
		token  = client.Token
	)
	//jwt 验证
	parseJWT, err := helper.ParseJWT(token)
	global.GVA_LOG.Infof("GetFishConfigController ParseJWT err %v ", parseJWT)
	if err != nil {
		code = common.TokenExpiration
		global.GVA_LOG.Infof("GetFishConfigController ParseJWT err ")
		return
	}

	//uid, err := strconv.Atoi(userID)
	if err != nil {
		code = common.TokenExpiration
		global.GVA_LOG.Error("处理数据 cmd ", zap.Any("err", err))
		return
	}

	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      userID,
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     userID,
			Code:    0,
			Message: "",
		},
		ServiceId: config.NatsSlotServer,
		//MsgId:     int32(pbs.Mmb_mtHeartReq),
		Content: message,
	}

	//组装
	NastManager.SendMemeJs(&msgReq)

	return code, "", nil
}

func MemeEntry(message []byte, uid string, msgId int32) (code uint32, msg string, data interface{}) {
	code = common.OK

	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      uid,
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     uid,
			Code:    0,
			Message: "",
		},
		ServiceId: config.NatsSlotServer,
		MsgId:     msgId,
		Content:   message,
	}

	//组装
	NastManager.SendMemeJs(&msgReq)

	return code, "", nil
}

// MemeBattleEntry 入口
func MemeBattleEntry(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	//jwt 验证
	var (
		userID = client.UserID
		token  = client.Token
	)
	parseJWT, err := helper.ParseJWT(token)
	global.GVA_LOG.Infof("GetFishConfigController ParseJWT err %v ", parseJWT)
	if err != nil {
		code = common.TokenExpiration
		global.GVA_LOG.Infof("GetFishConfigController ParseJWT err ")
		return
	}

	//uid, err := strconv.Atoi(userID)
	//cmdToInt, err := strconv.Atoi(cmd)
	if err != nil {
		code = common.TokenExpiration
		global.GVA_LOG.Error("处理数据 cmd ", zap.Any("err", err))
		return
	}
	MemeEntry(message, userID, int32(pbs.Meb_mtHeartReq))
	return
}

// QuickMatchRoom 匹配房间
// @Summary       meme-匹配房间
// @Tags          meme
// @Description   meme-匹配房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.MatchRoomReq                true  "匹配房间"
// @Success  1   {object}        common.JSONResult{data=models.MatchRoomAck} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebQuickMatchRoom   [post]
func QuickMatchRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.MatchRoomReq{}
		resp    = &models.MatchRoomAck{}
		userID  = client.UserID
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("QuickMatchRoom: %v %v", zap.Error(err))
		return
	}

	request.UserID = userID
	global.GVA_LOG.Infof("QuickMatchRoom: %v ,RoomNo : %v", client.UserID, request.RoomNo)

	//尝试获取锁
	cacheKey := userID + "QuickMatchRoom"
	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("QuickMatchRoom Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("QuickMatchRoom Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("QuickMatchRoom 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("QuickMatchRoom 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("QuickMatchRoom  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("QuickMatchRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.MatchRoomReq{
		RoomNo: request.RoomNo,
		UserId: userID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_memeMatchRoom))

	return code, "", resp
}

// MebCancelMatchRoom 取消快速匹配
// @Summary       meme-取消快速匹配
// @Tags          meme
// @Description   meme-取消快速匹配
// @Accept       json
// @Produce      json
// @Param        user  body      models.MatchRoomReq                true  "匹配房间"
// @Success  1   {object}        common.JSONResult{data=models.MatchRoomAck} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebCancelMatchRoom   [post]
func MebCancelMatchRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.MatchRoomReq{}
		resp    = &models.MatchRoomAck{}
		userID  = client.UserID
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebCancelMatchRoom:", zap.Error(err))
		return
	}

	request.UserID = userID
	global.GVA_LOG.Infof("MebCancelMatchRoom: %v ,RoomNo : %v", client.UserID, request.RoomNo)

	cacheKey := userID + "MebCancelMatchRoom"
	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebJoinRoom Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebJoinRoom Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebJoinRoom 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebJoinRoom 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebJoinRoom  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebJoinRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqProto := &pbs.MatchRoomReq{
		RoomNo: request.RoomNo,
		UserId: userID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_cancelMatchRoom))

	return code, "", resp
}

// GetUserState   获取用户当前游戏状态（包括角色信息）
// @Summary       meme-获取用户当前游戏状态
// @Tags          meme
// @Description   meme-获取用户当前游戏状态
// @Accept       json
// @Produce      json
// @Param        user  body      models.UserStateReq                true  "选择人物形象"
// @Success  1   {object}        common.JSONResult{data=models.UserStateResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /getUserState   [post]
func GetUserState(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.UserStateReq{}
		userID  = client.UserID
		resp    = &models.UserStateResp{}
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("GetUserState: %v %v", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("GetUserState %v", request)

	cacheKey := userID + "GetUserState"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("GetUserState QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqProto := &pbs.UserStateReq{
		UserId: userID,
		RoomNo: "",
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_userState))

	return code, "", resp
}

// GetRoomConfig  房间资费配置
// @Summary       meme-房间资费配置
// @Tags          meme
// @Description   meme-房间资费配置
// @Accept       json
// @Produce      json
// @Param        user  body      models.RoomConfigReq                true  "房间资费配置"
// @Success  1   {object}        common.JSONResult{data=models.RoomConfigResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebRoomConfig   [post]
func GetRoomConfig(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.RoomConfigReq{}
		resp    = &models.RoomConfigResp{}
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("GetRoomConfig: %v %v", zap.Error(err))
		return
	}

	//添加对局用户
	clientInfo := GetUserClient(common.AppId10, client.UserID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", client.UserID)
		code = common.NotLogin
		return
	}

	//todo 走grpc 请求

	return code, "", resp
}

// MebCreateRoom 创建房间
// @Summary       meme-创建房间
// @Tags          meme
// @Description   meme-创建房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.CreateRoomReq                true  "匹配房间"
// @Success  1   {object}        common.JSONResult{data=models.CreateRoomResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebCreateRoom   [post]
func MebCreateRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.CreateRoomReq{}
		//resp    = &models.CreateRoomResp{}
		userID = client.UserID
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebCreateRoom: %v %v", zap.Error(err))
		return
	}

	request.UserID = userID
	global.GVA_LOG.Infof("MebCreateRoom %v", string(message))

	cacheKey := userID + "MebCreateRoom"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebCreateRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//添加对局用户
	clientInfo := GetUserClient(common.AppId10, userID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", userID)
		code = common.NotLogin
		return
	}

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.CreateRoomReq{
		UserId:       userID,
		RoomType:     int32(request.RoomType),
		UserNumLimit: int32(request.UserNumLimit),
		RoomTurnNum:  int32(request.RoomTurnNum),
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_createRoom))

	global.GVA_LOG.Infof("MebCreateRoom %v", userID)
	return code, "", nil
}

// MebJoinRoom 加入房间
// @Summary       meme-加入房间
// @Tags          meme
// @Description   meme-加入房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.JoinRoomReq                true  "加入房间"
// @Success  1   {object}        common.JSONResult{data=models.JoinRoomResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebJoinRoom   [post]
func MebJoinRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.JoinRoomReq{}
		resp    = &models.JoinRoomResp{}
		userID  = client.UserID
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebJoinRoom: %v %v", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebJoinRoom %v", request)
	request.UserID = userID

	//加入房间
	clientInfo := GetUserClient(common.AppId10, request.UserID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebJoinRoom MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", request.UserID)
		code = common.NotLogin
		return
	}

	//尝试获取锁
	cacheKey := request.RoomNo + "MebJoinRoom"
	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebJoinRoom Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebJoinRoom Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebJoinRoom 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebJoinRoom 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebJoinRoom  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebJoinRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.JoinRoomReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_joinRoom))

	return code, "", resp
}

// MebReJoinRoom 重新加入房间
// @Summary       meme-重新加入房间
// @Tags          meme
// @Description   meme-重新加入房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.JoinRoomReq                true  "重新加入房间"
// @Success  1   {object}        common.JSONResult{data=models.JoinRoomResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /reJoinRoom   [post]
func MebReJoinRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.JoinRoomReq{}
		resp    = &models.ReJoinRoomResp{}
		userID  = client.UserID
	)

	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebReJoinRoom: %v %v", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebReJoinRoom %v", request)
	request.UserID = userID

	clientInfo := GetUserClient(common.AppId10, request.UserID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebReJoinRoom MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", request.UserID)
		code = common.NotLogin
		return
	}

	//尝试获取锁
	lockCtx := context.Background()
	cacheKey := request.RoomNo + "MebReJoinRoom"
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebReJoinRoom Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebReJoinRoom Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebReJoinRoom 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebReJoinRoom 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebReJoinRoom  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebReJoinRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.JoinRoomReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_reJoinRoom))

	return code, "", resp
}

// MebReady （就绪）
// @Summary       meme-就绪
// @Tags          meme
// @Description   meme-就绪
// @Accept       json
// @Produce      json
// @Param        user  body      models.ReadyReq                true  "加入房间"
// @Success  1   {object}        common.JSONResult{data=models.ReadyResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /ready   [post]
func MebReady(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.ReadyReq{}
		resp    = &models.ReadyResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebReady: %v %v", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebReady %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebReady"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebReady QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.ReadyReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_readyMsg))

	return code, "", resp
}

func MebCancelReady(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.CancelReadyReq{}
		resp    = &models.ReadyResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebCancelReady: %v %v", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebCancelReady %v", request)

	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebCancelReady"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebCancelReady QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.CancelReadyReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_cancelReady))

	return code, "", resp
}

// MebLeaveRoom 离开房间
// @Summary       meme-离开房间
// @Tags          meme
// @Description   meme-离开房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.LeaveRoomReq                true  "离开房间"
// @Success  1   {object}        common.JSONResult{data=models.LeaveRoomResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebLeaveRoom   [post]
func MebLeaveRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.LeaveRoomReq{}
		resp    = &models.LeaveRoomResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebLeaveRoom:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebLeaveRoom %v", request)
	request.UserID = userID

	//添加对局用户
	clientInfo := GetUserClient(common.AppId10, userID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", userID)
		code = common.NotLogin
		return
	}

	//尝试获取锁
	cacheKey := request.UserID + "MebLeaveRoom"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebLeaveRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.LeaveRoomReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_leaveRoom))

	return code, "", resp
}

// MebKickRoom 房主踢人
// @Summary       meme-离开房间
// @Tags          meme
// @Description   meme-离开房间
// @Accept       json
// @Produce      json
// @Param        user  body      models.KickRoomReq                true  "房主踢人"
// @Success  1   {object}        common.JSONResult{data=models.KickRoomResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebKickRoom   [post]
func MebKickRoom(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.KickRoomReq{}
		resp    = &models.KickRoomResp{}
		userID  = client.UserID //房主
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebKickRoom:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebKickRoom %v", request)
	if userID == request.UserID {
		code = common.NotSelfKickSelf
		return
	}

	if len(request.UserID) == 0 || len(request.RoomNo) == 0 {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebKickRoom: %v %v")
		return
	}

	//添加对局用户
	clientInfo := GetUserClient(common.AppId10, userID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebKickRoom 用户没有客户端,用户可能没登陆 UserID:%v ", userID)
		code = common.NotLogin
		return
	}

	//尝试获取锁
	lockCtx := context.Background()
	cacheKey := userID + "MebKickRoom"
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebKickRoom Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebKickRoom Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebKickRoom 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebKickRoom 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebKickRoom  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebKickRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.KickRoomReq{
		UserId:  request.UserID,
		RoomNo:  request.RoomNo,
		OwnerId: userID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_kickRoom))

	return code, "", resp
}

// MebInviteFriend 邀请好友
func MebInviteFriend(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.InviteFriendReq{}
		resp    = &models.InviteFriendResp{}
		userID  = client.UserID //房主
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebInviteFriend:", zap.Error(err))
		return
	}

	global.GVA_LOG.Infof("MebInviteFriend %v", request)

	//添加对局用户
	clientInfo := GetUserClient(common.AppId10, userID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebInviteFriend 用户没有客户端,用户可能没登陆 UserID:%v ", userID)
		code = common.NotLogin
		return
	}

	//被邀请人是否在线
	beInviterFriendClientInfo := GetUserClient(common.AppId10, request.UserID)
	if beInviterFriendClientInfo == nil || beInviterFriendClientInfo.UserID != request.UserID {
		global.GVA_LOG.Infof("MebInviteFriend 用户没有客户端,用户可能没登陆 UserID:%v ", userID)

		return
	}

	global.GVA_LOG.Infof("MebInviteFriend %v", request)
	if userID == request.UserID {
		code = common.NotSelfKickSelf
		return
	}

	if len(request.UserID) == 0 || len(request.RoomNo) == 0 {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebInviteFriend: %v %v")
		return
	}

	//尝试获取锁
	lockCtx := context.Background()
	cacheKey := userID + "MebInviteFriend"
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebInviteFriend Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebInviteFriend Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebInviteFriend 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebInviteFriend 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebInviteFriend  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebInviteFriend QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.InviteFriendReq{
		InviteUserId: request.UserID,
		RoomNo:       request.RoomNo,
		OwnerId:      userID,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_inviteFriend))

	return code, "", resp
}

// MebRoomAlive 房间心跳 房主开始游戏之后
// @Summary       meme-房间心跳
// @Tags          meme
// @Description   meme-房间心跳
// @Accept       json
// @Produce      json
// @Param        user  body      models.RoomAliveReq                true  "房间心跳"
// @Success  1   {object}        common.JSONResult{data=models.RoomAliveResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebRoomAlive   [post]
func MebRoomAlive(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.RoomAliveReq{}
		resp    = &models.RoomAliveResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebRoomAlive:", zap.Error(err))
		return
	}

	global.GVA_LOG.Infof("MebRoomAlive %v", request)
	request.UserID = userID

	clientInfo := GetUserClient(common.AppId10, request.UserID)
	if clientInfo == nil || clientInfo.UserID != client.UserID {
		global.GVA_LOG.Infof("MebJoinRoom MebCreateRoom 用户没有客户端,用户可能没登陆 UserID:%v ", request.UserID)
		code = common.NotLogin
		return
	}

	cacheKey := request.UserID + "MebRoomAlive"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebJoinRoom QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.RoomAliveReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_roomAlive))

	return code, "", resp
}

// MebStartPlay 房主开始对局游戏
// @Summary       meme-房主开始对局游戏
// @Tags          meme
// @Description   meme-房主开始对局游戏
// @Accept       json
// @Produce      json
// @Param        user  body      models.StartPlayReq                true  "房主开始对局游戏"
// @Success  1   {object}        common.JSONResult{data=models.StartPlayResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebStartPlay   [post]
func MebStartPlay(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK

	var (
		request = &models.StartPlayReq{}
		resp    = &models.StartPlayResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebStartPlay:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebStartPlay %v", request)

	request.UserID = userID

	//本地锁
	cacheKey := request.UserID + "MebStartPlay"
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebStartPlay QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.StartPlayReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_startPlay))

	return code, "", resp
}

// MebLoadCompleted 加载完成
// @Summary       meme-加载完成
// @Tags          meme
// @Description   meme-加载完成
// @Accept       json
// @Produce      json
// @Param        user  body      models.LoadCompletedReq                true  "请求骗子牌"
// @Success  1   {object}        common.JSONResult{data=models.LoadCompletedResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebLoadCompleted   [post]
func MebLoadCompleted(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.LoadCompletedReq{}
		resp    = &models.LoadCompletedResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebLoadCompleted:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebLoadCompleted %v", request)

	request.UserID = userID

	//本地锁
	cacheKey := userID + "MebLoadCompleted"
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebLoadCompleted QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.LoadCompletedReq{
		UserId: userID,
		RoomNo: request.RoomNo,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_loadCompleted))

	return code, "", resp
}

// MebOperateCard 出牌
// @Summary       meme-出牌
// @Tags          meme
// @Description   meme-出牌
// @Accept       json
// @Produce      json
// @Param        user  body      models.OperateCardReq                true  "出牌"
// @Success  1   {object}        common.JSONResult{data=models.OperateCardResp} "返回"
// @Failure      400  {object}   common.JSONResult                     "错误提示"
// @Router       /mebOperateCard  [post]
func MebOperateCard(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.OperateCardReq{}
		resp    = &models.OperateCardResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebOperateCard:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebOperateCard %v", request)

	request.UserID = userID

	//尝试获取锁
	lockCtx := context.Background()
	cacheKey := request.UserID + "MebOperateCard"
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebOperateCard Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebOperateCard Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebOperateCard 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebOperateCard 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebOperateCard  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebOperateCard QueueDataKeyMap TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	pbCard := make([]*pbs.Card, 0)
	//出牌
	if request.OpeType == 1 {
		for _, val := range request.Card {
			card := pbs.Card{CardId: int32(val.CardId)}
			pbCard = append(pbCard, &card)
		}
	}

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.OperateCardReq{
		UserId:  userID,
		RoomNo:  request.RoomNo,
		OpeType: int32(request.OpeType),
		EmojiId: request.EmojiId,
		Pitch:   request.Pitch,
		Yaw:     request.Yaw,
		Looking: request.Looking,
		Card:    pbCard,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_operateCards))

	return code, "", resp
}

// MebCardLike 点赞
func MebCardLike(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.CardLikeReq{}
		resp    = &models.CardLikeResp{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebCardLike:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebCardLike %v", request)
	request.UserID = userID

	//不能给自己点赞
	if userID == request.LikeUserID {
		code = common.ParameterIllegal
		return
	}

	//尝试获取锁
	lockCtx := context.Background()
	cacheKey := request.UserID + "MebCardLike"
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebCardLike Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebCardLike Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebCardLike 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebCardLike 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebCardLike  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebCardLike TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	pbCard := make([]*pbs.Card, 0)
	for _, val := range request.Card {
		card := pbs.Card{CardId: int32(val.CardId)}
		pbCard = append(pbCard, &card)
	}

	//uid, _ := strconv.Atoi(userID)
	reqProto := &pbs.LikeCardReq{
		LikeUserId: request.LikeUserID,
		UserId:     userID,
		RoomNo:     request.RoomNo,
		Card:       pbCard,
	}
	protoReq, _ := proto.Marshal(reqProto)
	MemeEntry(protoReq, userID, int32(pbs.Meb_likeCards))
	return code, "", resp
}

// MebHandbookList 图鉴列表
func MebHandbookList(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.HandbookListReq{}
		resp    = &pbs.HandbookListAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebHandbookList:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebHandbookList %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebHandbookList"

	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebHandbookList Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebHandbookList Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebHandbookList 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebHandbookList 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebHandbookList  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebHandbookList TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqProto := &pbs.HandbookListReq{
		UserId: request.UserID,
		LastId: int32(request.LastId),
		Level:  int32(request.Level),
	}

	reqDataMarshal, _ := proto.Marshal(reqProto)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_handbookList),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))

		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebHandbookList  :", zap.Error(err))

		return
	}
	global.GVA_LOG.Infof("MebHandbookList: %v", &resp)

	return code, "", resp
}

// MebUnpackCard 拆包
func MebUnpackCard(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.UnpackCardReq{}
		resp    = &pbs.UnpackCardAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebUnpackCard:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebUnpackCard %v", request)
	request.UserID = userID

	if request.Version <= 0 {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebUnpackCard: 版本错误")
		return
	}

	//if request.Num != 10 {
	//	request.Num = 1
	//}
	if !helper.InArr(request.Num, []int{1, 5}) {
		request.Num = 1
	}

	//尝试获取锁
	cacheKey := request.UserID + "MebUnpackCard"
	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebUnpackCard Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebUnpackCard Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebUnpackCard 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebUnpackCard 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebUnpackCard  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebUnpackCard TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.UnpackCardReq{
		UserId:  request.UserID,
		Version: int32(request.Version),
		Num:     int32(request.Num),
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_unpackCard),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebUnpackCard :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebUnpackCard: %v", &resp)
	return code, "", resp
}

func MebCardVersionList(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.CardVersionListReq{}
		resp    = &pbs.CardVersionListAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebCardVersionList:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebCardVersionList %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebCardVersionList"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebCardVersionList TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.CardVersionListReq{
		UserId: request.UserID,
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_cardVersionList),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))

		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebFriendList :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebCardVersionList: %v", &resp)

	return code, "", resp

}

func MebAuditUserList(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.AuditUserListReq{}
		resp    = &pbs.AuditUserAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebAuditUserList:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAuditUserList %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebAuditUserList"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebAuditUserList TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.AuditUserListReq{
		UserId: request.UserID,
		LastId: int32(request.LastId),
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_auditUserList),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))

		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebAuditUserList  :", zap.Error(err))

		return
	}
	global.GVA_LOG.Infof("MebAuditUserList: %v", &resp)

	return code, "", resp
}

func MebFriendList(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.FriendUserListReq{}
		resp    = &pbs.FriendListAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebFriendList:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebFriendList %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebFriendList"
	//本地锁
	err := global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebFriendList TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.FriendListReq{
		UserId: request.UserID,
		LastId: int32(request.LastId),
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_friendUserList),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebFriendList :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebFriendList: %v", &resp)

	return code, "", resp
}

func MebAddFriend(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.AddFriendReq{}
		resp    = &pbs.AddFriendAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebAddFriend:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAddFriend %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebAddFriend"

	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebAddFriend Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebAddFriend Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebAddFriend 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebAddFriend 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebAddFriend  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebAddFriend TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.AddFriendReq{
		AuditUser:       request.AuditUser,
		ApplicationUser: request.UserID,
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_addFriend),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebAddFriend :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAddFriend: %v", &resp)

	return code, "", resp
}

func MebDelFriend(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.DelFriendReq{}
		resp    = &pbs.DelFriendAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebDelFriend:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebDelFriend %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebDelFriend"

	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebDelFriend Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebDelFriend Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebDelFriend 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebDelFriend 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebDelFriend  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebDelFriend TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.DelFriendReq{
		UserId:   request.UserID,
		FriendId: int32(request.FriendId),
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_delFriend),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebDelFriend :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebDelFriend: %v", &resp)
	return code, "", resp
}

func MebAuthFriend(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.AuthFriendReq{}
		resp    = &pbs.AuthFriendAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebAuthFriend:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAuthFriend %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebAuthFriend"

	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebAuthFriend Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebAuthFriend Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebAuthFriend 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebAuthFriend 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebAuthFriend  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebAuthFriend TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.AuthFriendReq{
		UserId:  request.UserID,
		AuditId: int32(request.AuditId),
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_authFriend),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebAuthFriend :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAuthFriend: %v", &resp)
	return code, "", resp
}

func MebUserDetail(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.UserDetailReq{}
		resp    = &pbs.UserDetailAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebUserDetail:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebUserDetail %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebUserDetail"
	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebUserDetail Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebUserDetail Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebUserDetail 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebUserDetail 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebUserDetail  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebUserDetail TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.UserDetailReq{
		UserId: request.UserID,
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_userDetail),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebUserDetail :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebUserDetail: %v", &resp)
	return code, "", resp
}

func MebGetCoinExperience(client *Client, seq string, message []byte) (code uint32, msg string, data interface{}) {
	code = common.OK
	var (
		request = &models.CoinExperienceReq{}
		resp    = &pbs.CoinExperienceAck{}
		userID  = client.UserID
	)
	if err := json.Unmarshal(message, request); err != nil {
		code = common.ParameterIllegal
		global.GVA_LOG.Error("MebGetCoinExperience:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebGetCoinExperience %v", request)
	request.UserID = userID

	//尝试获取锁
	cacheKey := request.UserID + "MebGetCoinExperience"

	lockCtx := context.Background()
	lock := cache.NewRedisLock(cacheKey, 10*time.Second)
	acquired, err := lock.Acquire(lockCtx)
	if err != nil {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebGetCoinExperience Acquire %v", cacheKey)
		return
	}
	if !acquired {
		code = common.SysBusy
		global.GVA_LOG.Infof("MebGetCoinExperience Acquire %v", cacheKey)
		return
	}
	defer func() {
		released, err := lock.Release(lockCtx)
		if err != nil {
			global.GVA_LOG.Error("MebGetCoinExperience 释放锁失败", zap.Error(err), zap.Any("UserId", cacheKey))
		}
		if released {
			global.GVA_LOG.Infof("MebGetCoinExperience 成功释放锁 %v", cacheKey)
		} else {
			global.GVA_LOG.Infof("MebGetCoinExperience  锁已经被其他客户端占用，无法释放 %v", cacheKey)
		}
	}()

	//本地锁
	err = global.QueueDataKeyMap.TryAdd(cacheKey, helper.LocalTime().Unix())
	if err != nil {
		code = common.DuplicateRequests
		global.GVA_LOG.Infof("MebGetCoinExperience TryAdd%v", cacheKey)
		return
	}
	defer global.QueueDataKeyMap.Del(cacheKey)

	reqData := &pbs.CoinExperienceReq{
		UserId: request.UserID,
	}
	reqDataMarshal, _ := proto.Marshal(reqData)

	// 调用 gRPC 方法
	msgReq := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead:   &pbs.AckHead{},
		ServiceId: "",
		MsgId:     int32(pbs.Meb_coinExperience),
		Content:   reqDataMarshal,
	}

	response, err := grpcclient.GetMebClient().CallMebMethod(&msgReq)
	if response != nil && response.AckHead.Code != pbs.Code_OK {
		code = uint32(response.AckHead.Code)
		return
	}
	if err != nil || response == nil {
		global.GVA_LOG.Error("could not call method:", zap.Error(err))
		return
	}

	respData := response.Content
	err = proto.Unmarshal(respData, resp)
	if err != nil {
		global.GVA_LOG.Error("Unmarshal MebAuthFriend :", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("MebAuthFriend: %v", &resp)
	return code, "", resp
}
