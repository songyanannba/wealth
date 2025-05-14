package websocket

import (
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
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

	res := &pbs.CurrAPInfoAck{
		RoomNo:        "",
		CurrPeriod:    "",
		GameStartTime: 0,
		GameTurnState: 0,
		APRoomInfos:   nil,
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
			res.APRoomInfos = append(res.APRoomInfos, &pbs.APRoomInfos{
				UserId:    uid,
				BetZoneId: int32(uInfo.UserProperty.BetZoneId),
				Bet:       float32(uInfo.UserProperty.Bet),
			})
		}

		ptAck, _ := proto.Marshal(res)

		netMessageResp.Content = ptAck
		NatsSendAimUserMsg(roomSpaceInfo, netMessageResp, request.UserId)
	}
	return
}
