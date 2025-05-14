package websocket

import (
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"slot_server/lib/common"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/lib/src/dao"
	"slot_server/lib/src/logic"
	"slot_server/protoc/pbs"
	"slot_server/servers/meme_serve"
	"strconv"
	"time"
)

//func ProtoMTTest(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
//	reqData := &pbs.Test1Req{}
//
//	err := proto.Unmarshal(netMessage.Content, reqData)
//
//	fmt.Println(reqData, err)
//
//	ack := &pbs.Test1Ack{
//		UserId: "999",
//	}
//
//	ackMarshal, _ := proto.Marshal(ack)
//
//	ackMsg := common.GetErrorMessage(common.WebOK, "")
//	netMessageResp := &pbs.NetMessage{
//		ReqHead: &pbs.ReqHead{
//			Uid:      0,
//			Token:    "",
//			Platform: "",
//		},
//		AckHead: &pbs.AckHead{
//			Uid:     0,
//			Code:    0,
//			Message: ackMsg,
//		},
//		ServiceId: netMessage.ServiceId,
//		MsgId:     netMessage.MsgId + 1,
//		Content:   ackMarshal,
//	}
//	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
//
//	global.GVA_LOG.Infof("magic_tower send headByte:%v ", string(netMessageRespMarshal))
//	NastManager.Producer(netMessageRespMarshal)
//
//	return netMessage.MsgId + 1, common.WebOK, ackMarshal
//}

// Heart400 心跳 是否在房间
func Heart400(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//本协议不需要反解析解析协议
	//不需要发送消息回到网关

	ackMsg := common.GetErrorMessage(common.OK, "")
	netMessageResp := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     netMessage.ReqHead.Uid,
			Code:    0,
			Message: ackMsg,
		},
		ServiceId: netMessage.ServiceId,
		MsgId:     int32(pbs.Meb_mtHeart),
		Content:   nil,
	}

	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)

	global.GVA_LOG.Infof("Heart400 magic_tower send headByte:%v ", string(netMessageRespMarshal))
	NastManager.Producer(netMessageRespMarshal)

	//查看当前期房间是否创建 如果创建需要添加到 当前房间的容器中
	uidStr := netMessage.ReqHead.Uid

	err := meme_serve.MemeRoomManager.AddManager(uidStr, "")
	if err != nil {
		global.GVA_LOG.Error("add manager err:", zap.Any("err", err))
	}

	return netMessage.MsgId + 1, common.OK, nil
}

func MatchRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	code = common.OK
	request := &pbs.MatchRoomReq{}
	err := proto.Unmarshal(netMessage.Content, request)
	if err != nil {
		global.GVA_LOG.Error("MatchRoomController Unmarshal err:", zap.Any("err", err))
		return
	}
	global.GVA_LOG.Infof("MatchRoomController:%v ", request)

	netMessageResp := helper.NewNetMessage(netMessage.ReqHead.Uid, netMessage.ReqHead.Uid, int32(pbs.Meb_memeMatchRoom), netMessage.ServiceId)
	msgData := models.MatchSuccResp{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_memeMatchRoom)), //快速匹配成功协议 游戏开始
		Timestamp: time.Now().Unix(),
	}

	//roomSpace, err := SlotRoomManager.GetRoomSpace(request.RoomNo)
	//if err != nil {
	//	code = common.RoomNotExist
	//}

	//if roomSpace != nil && roomSpace.RoomInfo.Owner != request.UserId {
	//	code = common.NotRoomOwner
	//
	//	if roomSpace.RoomInfo.RoomType != table.RoomTypeMatch {
	//		code = common.NotRoomOwner
	//	}
	//}
	if code != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(code))
		global.GVA_LOG.Infof("MatchRoomController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	//匹配开始 把用户放入匹配中的结构中
	SlotRoomManager.JoinMatchIngRoom(request.RoomNo)

	//返回内容
	msgDataMarshal, _ := json.Marshal(msgData)
	netMessageResp.Content = msgDataMarshal
	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
	NastManager.Producer(netMessageRespMarshal)

	return int32(pbs.Meb_memeMatchRoom), code, nil
}

func CancelMatchRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	request := &pbs.MatchRoomReq{}
	err := proto.Unmarshal(netMessage.Content, request)
	if err != nil {
		global.GVA_LOG.Error("CancelMatchRoomController Unmarshal err:", zap.Any("err", err))
		return
	}
	global.GVA_LOG.Infof("CancelMatchRoomController:%v ", request)

	netMessageResp := helper.NewNetMessage(netMessage.ReqHead.Uid, netMessage.ReqHead.Uid, int32(pbs.Meb_cancelMatchRoom), netMessage.ServiceId)
	msgData := models.MatchSuccResp{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_cancelMatchRoom)), //快速匹配成功协议 游戏开始
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))
		global.GVA_LOG.Infof("CancelMatchRoomController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	//匹配开始 把用户放入匹配中的结构中
	SlotRoomManager.CancelMatchIngUser(request.RoomNo)

	//返回内容
	msgDataMarshal, _ := json.Marshal(msgData)
	netMessageResp.Content = msgDataMarshal
	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
	NastManager.Producer(netMessageRespMarshal)

	return int32(pbs.Meb_cancelMatchRoom), common.OK, nil
}

func CreateRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	netMessageResp := helper.NewNetMessage(netMessage.ReqHead.Uid, netMessage.ReqHead.Uid, int32(pbs.Meb_createRoom), netMessage.ServiceId)

	request := &pbs.CreateRoomReq{}
	err := proto.Unmarshal(netMessage.Content, request)
	if err != nil {
		global.GVA_LOG.Error("CreateRoomController Unmarshal err:", zap.Any("err", err))
		return
	}
	global.GVA_LOG.Infof("Heart400 magic_tower send headByte:%v ", request)

	//创建房间
	isCanCode := logic.IsCanCreateOrJoinRoom(request.UserId)
	if isCanCode != common.OK {
		code = uint32(isCanCode)
		global.GVA_LOG.Error("CreateRoomController IsCanCreateOrJoinRoom 先离开原来的房间 ")

		netMessageResp.AckHead.Code = pbs.Code(int32(code))

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		global.GVA_LOG.Infof("CreateRoomController:%v ", string(netMessageRespMarshal))
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	//创建房间
	roomInfo, err := logic.SaveRoom(request.UserId, int(request.RoomType), int(request.UserNumLimit), 0, int(request.RoomTurnNum), 0)
	if err != nil {
		code = common.TavernCreateRoomErr
		global.GVA_LOG.Error("CreateRoomController SaveRoom: %v %v", zap.Error(err))

		netMessageResp.AckHead.Code = pbs.Code(int32(code))

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		global.GVA_LOG.Infof("CreateRoomController:%v ", string(netMessageRespMarshal))
		NastManager.Producer(netMessageRespMarshal)
		return
	}
	global.GVA_LOG.Infof("CreateRoomController roomInfo:%v", roomInfo)

	//加入房间管理器

	//创建房间以后 房间会进入房间管理器
	//1 先创建对局空间

	roomSpaceInfo := GetRoomSpace()
	roomSpaceInfo.RoomInfo = roomInfo
	roomSpaceInfo.ComRoomSpace.AddTurn()

	//创建房间 并添加用户
	userInfo := &models.UserInfo{
		UserID:   request.UserId,
		Nickname: "",
		UserProperty: models.UserProperty{
			IsOwner: true,
			Seat:    1,
		},
		UserExt: models.UserExt{
			//RoomNo: roomInfo.RoomNo,
		},
	}

	//就绪
	userInfo.SetUserIsReady(1)

	//设置房主
	roomSpaceInfo.ComRoomSpace.UserOwner = userInfo

	//保存用户信息
	roomSpaceInfo.ComRoomSpace.AddUserInfos(request.UserId, userInfo)

	//2 添加到全局房间管理器
	//SlotRoomManager.AddRoomSpace(roomInfo.RoomNo, roomSpaceInfo)
	//
	//roomUserList, _ := dao.GetRoomUser(roomInfo.RoomNo, roomSpaceInfo.ComRoomSpace.GetTurn())
	//
	////发送广播 谁加入房间
	//msgData := models.CreateRoomMsg{
	//	ProtoNum:  strconv.Itoa(int(pbs.Meb_createRoom)),
	//	Timestamp: time.Now().Unix(),
	//	RoomCom: models.RoomCom{
	//		UserId:       request.UserId,
	//		RoomNo:       roomInfo.RoomNo,
	//		RoomName:     roomInfo.Name,
	//		Status:       roomInfo.IsOpen,
	//		UserNumLimit: roomInfo.UserNumLimit,
	//		RoomType:     int(roomInfo.RoomType),
	//		RoomLevel:    int(roomInfo.RoomLevel),
	//	},
	//	RoomUserList: roomUserList,
	//}

	//responseHeadByte, _ := json.Marshal(msgData)

	//给客户消息
	//global.GVA_LOG.Infof("CreateRoomController 加入房间的广播: %v", string(responseHeadByte))
	//
	////每个小房间是一个 协成
	//go roomSpaceInfo.Start()
	//
	//netMessageResp.Content = responseHeadByte
	//netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
	//global.GVA_LOG.Infof("CreateRoomController:%v ", string(netMessageRespMarshal))
	//NastManager.Producer(netMessageRespMarshal)
	return int32(pbs.Meb_createRoom), common.OK, nil
}

func ReadyRoomRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.ReadyReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("ReadyRoomRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("ReadyRoomRoomController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_readyMsg), config.SlotServer)
	msgData := models.ReadyMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_readyMsg)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("NatsSendAimUserMsg LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_readyMsg)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("MebJoinRoom GetRoomSpace ", zap.Error(err))
		return
	}
	return
}

func CancelReadyRoomRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.CancelReadyReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("CancelReadyRoomRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("CancelReadyRoomRoomController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_cancelReadyMsg), config.SlotServer)
	msgData := models.ReadyMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_cancelReadyMsg)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("NatsSendAimUserMsg LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_cancelReady)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("MebJoinRoom GetRoomSpace ", zap.Error(err))
		return
	}
	return
}

func JoinRoomRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.JoinRoomReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("JoinRoomRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("JoinRoomRoomController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_joinRoom), config.SlotServer)
	msgData := models.JoinRoomMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_joinRoom)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := JoinRoomVerifyParas(request.UserId, request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("NatsSendAimUserMsg LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_joinRoom)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("MebJoinRoom GetRoomSpace ", zap.Error(err))
		return
	}
	return
}

func RoomIsExist(roomNo string) int {
	code := common.OK
	_, err := SlotRoomManager.GetRoomSpace(roomNo)
	if err != nil {
		code = common.RoomNotExist
		return code
	}
	return code
}

func JoinRoomVerifyParas(userID, roomNo string) int {
	code := common.OK
	exist := RoomIsExist(roomNo)
	if exist != code {
		return exist
	}

	//先离开原来的房间
	roomDetail := logic.UserIsJoinRoom(userID)
	global.GVA_LOG.Infof("MebJoinRoom roomDetail %v", roomDetail)

	if roomDetail != nil && len(roomDetail.RoomNo) != 0 && roomDetail.RoomNo != roomNo && roomDetail.Status > 0 {
		//通知房间谁必须先离开房间
		code = common.LeavePreRoom
		global.GVA_LOG.Error("MebJoinRoom 先离开原来的房间 ")
		return code
	}

	//查看房间是否存在
	roomRecord, err := table.SlotRoomByRoomNo(roomNo)
	if err != nil {
		code = common.ServerError
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}
	if roomRecord.ID <= 0 {
		code = common.NotData
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}

	//已结束
	if roomRecord.IsOpen == table.RoomStatusDissolve || roomRecord.IsOpen == table.RoomStatusStop {
		code = common.RoomStatusStopErr
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}

	//房间满员
	if roomRecord.IsOpen == table.RoomStatusFill {
		code = common.JoinRoomFull
		global.GVA_LOG.Info("MebJoinRoom JoinRoomFull")
		return code
	}

	return code
}

func ReJoinRoomVerifyParas(userID, roomNo string) int {
	code := common.OK
	exist := RoomIsExist(roomNo)
	if exist != code {
		return exist
	}

	//先离开原来的房间
	roomDetail := logic.UserIsJoinRoom(userID)
	global.GVA_LOG.Infof("MebJoinRoom roomDetail %v", roomDetail)

	if roomDetail != nil && len(roomDetail.RoomNo) != 0 && roomDetail.RoomNo != roomNo && roomDetail.Status > 0 {
		//通知房间谁必须先离开房间
		code = common.LeavePreRoom
		global.GVA_LOG.Error("MebJoinRoom 先离开原来的房间 ")
		return code
	}

	//查看房间是否存在
	roomRecord, err := table.SlotRoomByRoomNo(roomNo)
	if err != nil {
		code = common.ServerError
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}
	if roomRecord.ID <= 0 {
		code = common.NotData
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}

	//已结束
	if roomRecord.IsOpen == table.RoomStatusDissolve || roomRecord.IsOpen == table.RoomStatusStop {
		code = common.RoomStatusStopErr
		global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
		return code
	}
	return code
}

func LeaveRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.LeaveRoomReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("LeaveRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("LeaveRoomController %v", request)

	//查看房间是否存在
	roomRecord, err := table.SlotRoomByRoomNo(request.RoomNo)
	if err != nil {
		code = common.ServerError
		global.GVA_LOG.Error("MebLeaveRoom ", zap.Error(err))
		return
	}
	if roomRecord.ID <= 0 {
		code = common.NotData
		global.GVA_LOG.Error("MebLeaveRoom ", zap.Error(err))
		return
	}

	_, err = SlotRoomManager.GetRoomSpace(roomRecord.RoomNo)
	if err != nil {
		netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_leaveRoom), config.SlotServer)

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(common.RoomNotExist))

		global.GVA_LOG.Infof("NatsSendAimUserMsg LikeUserId:{%v}", request.UserId)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_leaveRoom)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err = SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		code = common.RoomNotExist
		global.GVA_LOG.Error("MebLeaveRoom RoomNotExist ", zap.Error(err))
		return
	}
	return
}

func UserStateController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.UserStateReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("UserStateController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("UserStateController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_userState), config.SlotServer)
	res := &models.UserStateMsg{
		ProtoNum:   strconv.Itoa(int(pbs.Meb_userState)),
		Timestamp:  time.Now().Unix(),
		UserId:     request.UserId,
		IsContinue: false,
		RoomDetail: nil,
	}
	userRoomStatus, err := table.GetUsersRoomByUid(request.UserId)
	if err != nil {
		code = common.DBErr
		global.GVA_LOG.Error("GetUserState GetTavernUsersRoomByUid: %v %v", zap.Error(err))
	}
	request.RoomNo = userRoomStatus.RoomNo

	var isContinue bool
	//是否首次5连抽
	record, _ := table.UserIsHaveHandbook(request.UserId)
	if record.ID > 0 {
		isContinue = true
	}
	res.IsContinue = isContinue

	//房间是否存活
	roomSpaceInfo, err := SlotRoomManager.GetRoomSpace(request.RoomNo)
	global.GVA_LOG.Infof("UserStateController roomSpaceInfo %v", &roomSpaceInfo)

	if err != nil {
		//房间已经被回收 ，房间管理器没有房间
		//但是数据库 房间状态还是在进行中 需要销毁房间
		record, err := table.SlotRoomByRoomNo(request.RoomNo)
		if err != nil {
			global.GVA_LOG.Error("UserStateController: %v %v", zap.Error(err))
		}

		//房间状态: 1=开放中,2=已满员,3=已解散,4=进行中,5=已结束 6=异常房间 7=服务字段清理残存房间
		if record.ID > 0 && helper.InArr(int(record.IsOpen), []int{table.RoomStatusOpen, table.RoomStatusFill, table.RoomStatusIng}) {
			record.IsOpen = table.RoomStatusAbnormal
			err = table.SaveMemeRoom(record)
			if err != nil {
				global.GVA_LOG.Error(" UserStateController MebLeaveRoom", zap.Any("err", err))
			}
		}

		//返回数据，没有房间信息
		userStateRespMarshal, _ := json.Marshal(res)
		netMessageResp.Content = userStateRespMarshal
		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId
		global.GVA_LOG.Infof("InviteFriend LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, res)
		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)

	} else {
		//房间没有被销毁
		//获取一下用户维度的数据 （用户维度会在用户加入房间的时候 保存最近一场的数据，用户离开也会清理数据）
		userStatus, err := table.GetUsersRoomByUid(request.UserId)
		if err != nil {
			global.GVA_LOG.Error("GetUserState GetTavernUsersRoomByUid: %v %v", zap.Error(err))
			netMessageResp.AckHead.Code = pbs.Code(common.DBErr)
			NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, request.UserId)
			return
		}

		roomUserLists, _ := dao.GetRoomUser(userStatus.RoomNo, roomSpaceInfo.ComRoomSpace.GetTurn())
		res = &models.UserStateMsg{
			ProtoNum:   strconv.Itoa(int(pbs.Meb_userStateAck)),
			Timestamp:  time.Now().Unix(),
			UserId:     request.UserId,
			RoomNo:     userStatus.RoomNo,
			IsContinue: isContinue,
			RoomDetail: &models.RoomItem{
				RoomCom: models.RoomCom{
					RoomId: 0,
					Turn:   roomSpaceInfo.ComRoomSpace.GetTurn(),
					RoomNo: userStatus.RoomNo,
					UserId: userStatus.UserId,
					//RoomName:     roomSpaceInfo.RoomInfo.Name,
					//Status:       roomSpaceInfo.RoomInfo.IsOpen,
					//UserNumLimit: roomSpaceInfo.RoomInfo.UserNumLimit,
					//RoomType:     int(roomSpaceInfo.RoomInfo.RoomType),
					//RoomLevel:    int(roomSpaceInfo.RoomInfo.RoomLevel),
				},
				RoomUserList: roomUserLists,
			},
		}

		userStateRespMarshal, _ := json.Marshal(res)
		netMessageResp.Content = userStateRespMarshal
		NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, request.UserId)
	}
	return
}

func ReJoinRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.JoinRoomReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("ReJoinRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("ReJoinRoomController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_reJoinRoom), config.SlotServer)
	msgData := models.JoinRoomMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_reJoinRoom)),
		Timestamp: time.Now().Unix(),
	}

	resCode := common.OK
	roomUsersByRoomNoAndUid, err := table.RoomUsersByRoomNoAndUid(request.RoomNo, request.UserId)
	if err != nil {
		resCode = common.NotRoomOwner
		global.GVA_LOG.Error("ReJoinRoomController RoomUsersByRoomNoAndUid", zap.Error(err))
	}
	if roomUsersByRoomNoAndUid.ID < 0 {
		resCode = common.NotRoomOwner
		global.GVA_LOG.Infof("ReJoinRoomController %v", &roomUsersByRoomNoAndUid)
	}

	vfCode := ReJoinRoomVerifyParas(request.UserId, request.RoomNo)
	if vfCode != common.OK || resCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))
		global.GVA_LOG.Infof("NatsSendAimUserMsg LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)
		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_reJoinRoom)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err = SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("ReJoinRoomController ", zap.Error(err))
		return
	}
	return
}

func KickRoomController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.KickRoomReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("KickRoomController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("KickRoomController %v", request)

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_kickRoom)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("ReJoinRoomController ", zap.Error(err))
		return
	}
	return
}

func InviteFriendController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.InviteFriendReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("InviteFriendInter:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("InviteFriendInter %v", request)

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_inviteFriend)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)
	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("InviteFriendInter ", zap.Error(err))
		return
	}
	return
}

func StartPlayController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.StartPlayReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("StartPlayController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("StartPlayController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_startPlay), config.SlotServer)
	msgData := models.StartPlayMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_startPlay)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("StartPlayController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_startPlay)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)
	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("StartPlayController ", zap.Error(err))
		return
	}
	return
}

func LoadCompletedController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.LoadCompletedReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("LoadCompletedController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("LoadCompletedController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_loadCompleted), config.SlotServer)
	msgData := models.LoadMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_loadCompleted)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal
		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId
		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))
		global.GVA_LOG.Infof("LoadCompletedController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)
		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_loadCompleted)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("LoadCompletedController ", zap.Error(err))
		return
	}

	return
}

func LikeCardsController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.LikeCardReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("LikeCardsController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("LikeCardsController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_likeCards), config.SlotServer)
	msgData := models.OperateCardsMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_likeCards)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("LikeCardsController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_likeCards)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("LoadCompletedController ", zap.Error(err))
	}

	return
}

func OperateCardController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.OperateCardReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("OperateCardController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("OperateCardController %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, "", int32(pbs.Meb_operateCards), config.SlotServer)
	msgData := models.OperateCardsMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_operateCards)),
		Timestamp: time.Now().Unix(),
	}

	vfCode := RoomIsExist(request.RoomNo)
	if vfCode != common.OK {
		//返回内容
		msgDataMarshal, _ := json.Marshal(msgData)
		netMessageResp.Content = msgDataMarshal

		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId

		//返回的code
		netMessageResp.AckHead.Code = pbs.Code(int32(vfCode))

		global.GVA_LOG.Infof("OperateCardController LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, msgData)

		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
		return
	}

	comMsg := &models.ComMsg{
		MsgId: strconv.Itoa(int(pbs.Meb_operateCards)),
		Data:  netMessage.Content,
	}
	comMsgMarshal, _ := json.Marshal(comMsg)

	err := SlotRoomManager.SendMsgToRoomSpace(request.RoomNo, comMsgMarshal)
	if err != nil {
		global.GVA_LOG.Error("LoadCompletedController ", zap.Error(err))
	}

	return
}

func RoomAliveController(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.RoomAliveReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("RoomAliveController:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("RoomAliveController %v", request)

	roomSpaceInfo, err := SlotRoomManager.GetRoomSpace(request.RoomNo)
	if err != nil {
		global.GVA_LOG.Error("RoomAliveController GetRoomSpace ", zap.Error(err))
		return
	}

	//更新当前时间 认为客户端存活
	roomSpaceInfo.ComRoomSpace.CurrentOpTime = time.Now().Unix()
	return
}
