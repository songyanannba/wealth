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

	aPRoomInfos := &pbs.APRoomInfos{}

	res := &pbs.CurrAPInfoAck{
		RoomNo:        "",
		CurrPeriod:    "",
		GameStartTime: 0,
		GameTurnState: 0,
		APRoomInfos:   aPRoomInfos,
		GameStates:    0,
	}
	getBetZoneFigure := GetBetZoneFigure()

	for _, betZoneFigure := range getBetZoneFigure {
		var colorIdArr []int32
		if betZoneFigure.ColorId != nil && len(betZoneFigure.ColorId) > 0 {
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
		getGameState := roomSpaceInfo.ComRoomSpace.GetGameState()

		//if getGameState == BetIng && helper.LocalTime().Unix()-res.GameStartTime <= BetIngPeriod {
		//	res.GameStates = 1
		//}
		res.GameStates = 2
		if getGameState == BetIng {
			res.GameStates = 1
		}

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

		//todo
		if len(roomSpaceInfo.ComRoomSpace.CurrAnimalWheelSort) > 0 {
			currAnimalWheelSort := roomSpaceInfo.ComRoomSpace.CurrAnimalWheelSort
			for _, animalWheelSort := range currAnimalWheelSort {
				winBetZoneConfig := &pbs.WinBetZoneConfig{
					WinSeat:      int32(animalWheelSort.WinSeat),
					AnimalConfig: make([]*pbs.AnimalConfig, 0),
					WinZoneConf:  make([]*pbs.WinZoneConf, 0),
				}

				for _, animalConf := range animalWheelSort.AnimalConfigs {
					winBetZoneConfig.AnimalConfig = append(winBetZoneConfig.AnimalConfig, &pbs.AnimalConfig{
						Seat:     int32(animalConf.Seat),
						AnimalId: int32(animalConf.AnimalId),
					})
				}

				for _, bigOrSmallConfig := range animalWheelSort.BigOrSmallConfigs {
					winBetZoneConfig.BigSmallConfig = append(winBetZoneConfig.BigSmallConfig, &pbs.BigOrSmallConfig{
						Seat:       int32(bigOrSmallConfig.Seat),
						BigSmallId: int32(bigOrSmallConfig.BigOrSmall),
					})
				}

				//对应位置的颜色
				colorConfigSeat := roomSpaceInfo.GetColorConfigsBySeat(animalWheelSort.WinAnimalConfig.Seat)
				//根据本局赢钱的位置的动物和颜色确定赔率
				betZoneConfig := GetBetZoneConfigByAnimalIdAndColorId(animalWheelSort.WinAnimalConfig.AnimalId, colorConfigSeat.ColorId)
				animalWheelSort.WinBetZoneConfig = betZoneConfig

				for _, bzz := range betZoneConfig {
					winBetZoneConfig.WinZoneConf = append(winBetZoneConfig.WinZoneConf, &pbs.WinZoneConf{
						BetZoneId: int32(bzz.Seat),
						BetRate:   float32(bzz.BetRate),
					})
				}
				res.APRoomInfos.WinBetZoneConfig = append(res.APRoomInfos.WinBetZoneConfig, winBetZoneConfig)
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
	res := &pbs.UserBetAck{
		Bet:       request.Bet,
		GameId:    request.GameId,
		BetZoneId: request.BetZoneId,
		UserId:    netMessage.ReqHead.Uid,
	}

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
			} else {
				userInfo.UserProperty.Bet += float64(request.Bet)
			}

			userInfo, _ = roomSpaceInfo.ComRoomSpace.GetUserInfo(netMessage.ReqHead.Uid)
			if userInfo == nil {
				global.GVA_LOG.Error("押注保留用户信息错误 UserBetReq userInfo")
				return
			}

			//保留押注
			roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(int(request.BetZoneId), request.Bet, userInfo.Copy())

			//测试 多压几个 todo
			//{
			//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(0, 1, userInfo.Copy())
			//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(1, 2, userInfo.Copy())
			//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(3, 3, userInfo.Copy())
			//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(8, 4, userInfo.Copy())
			//}

			ptAck, _ := proto.Marshal(res)
			netMessageResp.Content = ptAck
			NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, "")
		}
	}

	return
}
