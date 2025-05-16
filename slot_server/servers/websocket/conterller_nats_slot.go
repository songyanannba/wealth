package websocket

import (
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/protoc/pbs"
)

func CurrAPInfos(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.NatsCurrAPInfo{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("CurrAPInfos:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("CurrAPInfos %v", request)

	netMessageResp := helper.NewNetMessage(request.UserId, request.UserId, int32(pbs.ProtocNum_CurrAPInfoAck), config.SlotServer)

	aPRoomInfos := &pbs.APRoomInfos{
		UserBetInfos: make([]*pbs.UserBetInfos, 0),
		ColorConfig:  make([]*pbs.ColorConfig, 0),
		AnimalConfig: make([]*pbs.AnimalConfig, 0),
	}

	res := &pbs.CurrAPInfoAck{
		RoomNo:        "",
		CurrPeriod:    "",
		GameStartTime: 0,
		GameTurnState: 0,
		APRoomInfos:   aPRoomInfos,
		BetZoneConf:   make([]*pbs.BetZoneConfig, 0),
	}

	for _, betZoneFigure := range GetBetZoneFigure() {
		var colorIdArr []int32
		if len(betZoneFigure.ColorId) > 0 {
			for _, colorId := range betZoneFigure.ColorId {
				colorIdArr = append(colorIdArr, int32(colorId))
			}
		}
		res.BetZoneConf = append(res.BetZoneConf, &pbs.BetZoneConfig{
			Seat:     int32(betZoneFigure.Seat),
			AnimalId: int32(betZoneFigure.AnimalId),
			ColorId:  colorIdArr,
			Size:     int32(betZoneFigure.Size),
			BetRate:  float32(betZoneFigure.BetRate),
		})
	}

	//房间是否存活
	roomSpaceInfo, err := SlotRoomManager.GetCurrRoomSpace()
	global.GVA_LOG.Infof("CurrAPInfos roomSpaceInfo %v", &roomSpaceInfo)

	if err != nil {
		//返回数据，没有房间信息
		ptAck, _ := proto.Marshal(res)
		netMessageResp.Content = ptAck
		//返回的用户id
		netMessageResp.AckHead.Uid = request.UserId
		netMessageResp.AckHead.Code = pbs.Code(pbs.ErrCode_NotRoom)
		global.GVA_LOG.Infof("CurrAPInfos LikeUserId:{%v} 给客户端发消息:{%v}", request.UserId, res)
		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)

	} else {
		//房间没有被销毁
		//获取一下用户维度的数据 （用户维度会在用户加入房间的时候 保存最近一场的数据，用户离开也会清理数据）

		res.RoomNo = roomSpaceInfo.RoomInfo.RoomNo
		res.CurrPeriod = roomSpaceInfo.RoomInfo.Period
		res.GameStartTime = roomSpaceInfo.ComRoomSpace.GetGameStartTime()

		for uid, uInfo := range roomSpaceInfo.ComRoomSpace.UserInfos {
			res.APRoomInfos.UserBetInfos = append(res.APRoomInfos.UserBetInfos, &pbs.UserBetInfos{
				UserId:    uid,
				BetZoneId: int32(uInfo.UserProperty.BetZoneId),
				Bet:       float32(uInfo.UserProperty.Bet),
			})
		}

		for _, colorConfig := range roomSpaceInfo.ColorConfigs {
			res.APRoomInfos.ColorConfig = append(res.APRoomInfos.ColorConfig, &pbs.ColorConfig{
				Seat:    int32(colorConfig.Seat),
				ColorId: int32(colorConfig.ColorId),
			})
		}

		if len(roomSpaceInfo.ComRoomSpace.CurrAnimalConfigs) > 0 {
			for _, animalConfig := range roomSpaceInfo.ComRoomSpace.CurrAnimalConfigs {
				res.APRoomInfos.AnimalConfig = append(res.APRoomInfos.AnimalConfig, &pbs.AnimalConfig{
					Seat:     int32(animalConfig.Seat),
					AnimalId: int32(animalConfig.AnimalId),
				})
			}
		}

		ptAck, _ := proto.Marshal(res)

		netMessageResp.Content = ptAck
		NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, request.UserId)
	}
	return
}

func UserBetReq(netMessage *pbs.NetMessage) (respMsgId int32, code uint32, data []byte) {
	//解析请求参数
	request := &pbs.UserBetReq{}
	if err := proto.Unmarshal(netMessage.Content, request); err != nil {
		global.GVA_LOG.Error("UserBetReq:", zap.Error(err))
		return
	}
	global.GVA_LOG.Infof("UserBetReq %v", request)

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_betAck), config.SlotServer)

	//获取当前的对局

	//是否是押注时间段

	res := &pbs.UserBetAck{}

	//房间是否存活
	roomSpaceInfo, err := SlotRoomManager.GetCurrRoomSpace()
	global.GVA_LOG.Infof("UserBetReq roomSpaceInfo %v", &roomSpaceInfo)

	if err != nil {
		//返回数据，没有房间信息
		ptAck, _ := proto.Marshal(res)
		netMessageResp.Content = ptAck
		//返回的用户id
		netMessageResp.AckHead.Uid = netMessage.ReqHead.Uid
		netMessageResp.AckHead.Code = pbs.Code(pbs.ErrCode_NotRoom)
		global.GVA_LOG.Infof("UserBetReq LikeUserId:{%v} 给客户端发消息:{%v}", netMessage.ReqHead.Uid, res)
		netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
		NastManager.Producer(netMessageRespMarshal)
	} else {

		gState := roomSpaceInfo.ComRoomSpace.GetGameState()
		currGameStartTime := roomSpaceInfo.ComRoomSpace.GetGameStartTime()
		global.GVA_LOG.Infof("UserBetReq currTime - currGameStartTime:%v 执行, gState:%v ,currGameStartTime:%v", helper.LocalTime().Unix()-currGameStartTime, gState, currGameStartTime)

		//押注期间
		if gState != BetIng {
			//不是押注时间
			netMessageResp.AckHead.Uid = netMessage.ReqHead.Uid
			netMessageResp.AckHead.Code = pbs.Code(pbs.ErrCode_NotBetPeriod)
			global.GVA_LOG.Infof("UserBetReq LikeUserId:{%v} 给客户端发消息:{%v}", netMessage.ReqHead.Uid, res)
			netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
			NastManager.Producer(netMessageRespMarshal)
			return
		} else {
			//保留押注用户信息
			userInfo, _ := roomSpaceInfo.ComRoomSpace.GetUserInfo(netMessage.ReqHead.Uid)
			if userInfo == nil {
				//保存用户信息
				user := models.NewUserInfo(netMessage.ReqHead.Uid, "", models.NewUserProperty(0, 0, false, float64(request.Bet)), models.UserExt{
					RoomNo: roomSpaceInfo.RoomInfo.RoomNo,
				})
				userInfo = &user
				global.GVA_LOG.Infof("UserBetReq 押注 :%v", userInfo)

				roomSpaceInfo.ComRoomSpace.AddUserInfos(netMessage.ReqHead.Uid, userInfo) //JoinRoom
			}

			userInfo, _ = roomSpaceInfo.ComRoomSpace.GetUserInfo(netMessage.ReqHead.Uid)
			if userInfo == nil {
				global.GVA_LOG.Error("押注保留用户信息错误 UserBetReq userInfo")
				return
			}

			//保留押注
			roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(int(request.BetZoneId), request.Bet, userInfo)
			ptAck, _ := proto.Marshal(res)
			netMessageResp.Content = ptAck
			NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, netMessage.ReqHead.Uid)
		}
	}

	return
}
