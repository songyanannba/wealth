package websocket

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/protoc/pbs"
)

// APRobotAction 动物派对机器人行为
func (trs *RoomSpace) APRobotAction() {
	//机器人押注前端表现

	var (
		//押注金额 Bet
		bet = helper.RandInt(100) + 1
		//游戏ID GameId
		gameId = 1
		//押注区域
		betZoneId = helper.RandInt(18)
		//用户ID
		robotUid = uuid.New().String()
	)

	res := &pbs.UserBetAck{
		Bet:       float32(bet),
		GameId:    int32(gameId),
		BetZoneId: int32(betZoneId),
		UserId:    robotUid,
		IsRobot:   true,
	}

	userInfo, _ := trs.ComRoomSpace.GetUserInfo(robotUid)
	if userInfo == nil {
		//保存用户信息
		user := models.NewUserInfo(robotUid, "xxx", models.NewUserProperty(0, 0, false, float64(bet)), models.UserExt{})
		user.UserProperty.IsRobot = 1
		userInfo = &user
		global.GVA_LOG.Infof("APRobotAction 机器人押注 :%v", userInfo)
		trs.ComRoomSpace.AddUserInfos(robotUid, userInfo) //JoinRoom
	}

	trs.ComRoomSpace.APRobotActionCount++

	//保留押注
	//trs.ComRoomSpace.AddBetZoneUserInfoMap(betZoneId, bet, userInfo.Copy())
	//测试 多压几个 todo
	//{
	//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(0, 1, userInfo.Copy())
	//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(1, 2, userInfo.Copy())
	//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(3, 3, userInfo.Copy())
	//	roomSpaceInfo.ComRoomSpace.AddBetZoneUserInfoMap(8, 4, userInfo.Copy())
	//}

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_betAck), config.SlotServer)
	ptAck, _ := proto.Marshal(res)
	netMessageResp.Content = ptAck
	NatsSendAimUserMsg(trs, netMessageResp, "")
}
