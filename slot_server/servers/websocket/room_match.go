package websocket

import (
	"encoding/json"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/protoc/pbs"
	"strconv"
	"time"
)

func (trMgr *roomManager) DealMatchIng1User() {
	match4 := 4
	matchLimits := 12 //三组以上

	//获取双排人数
	matchIngRoom1User := trMgr.MatchIngRoom.MatchIngRoom1User

	if len(matchIngRoom1User) < match4 {
		//不够4个人
		return
	}

	//考虑锁的时间 太长 每次匹配最多10组
	//和快速匹配会竞争锁
	global.GVA_LOG.Infof("DealMatchIng1User 单排 {%v}人场匹配期...", match4)

	matchIngRoom1UserGroup := make([]*MatchIngRoomInfo, 0)

	if len(matchIngRoom1User) > matchLimits {
		//每次取前10个房间
		matchIngRoom1UserGroup = matchIngRoom1User[0:matchLimits]
		trMgr.MatchIngRoom.MatchIngRoom1User = matchIngRoom1User[matchLimits:]
	} else {
		if len(matchIngRoom1User)%4 == 0 {
			matchIngRoom1UserGroup = matchIngRoom1User
			trMgr.MatchIngRoom.MatchIngRoom1User = []*MatchIngRoomInfo{}
		} else {
			matchIngRoom1UserGroup = matchIngRoom1User[0:match4]
			trMgr.MatchIngRoom.MatchIngRoom1User = matchIngRoom1User[match4:]
		}
	}

	// 前两个用户匹配成功
	matchUser := &MatchGroupRoomInfo{}

	//匹配成功 多少对
	matchRoomUserArr := make([]*MatchGroupRoomInfo, 0)

	for _, matchRoomInfo := range matchIngRoom1UserGroup {
		if len(matchRoomInfo.UserInfoArr) != 1 {
			//没有足够的匹配用户
			global.GVA_LOG.Infof("DealMatchIng1User 没有足够的匹配用户 RoomNo:%v, UserInfoMapLen:%v", matchRoomInfo.RoomNo, len(matchRoomInfo.UserInfoArr))
			continue
		}

		matchUser.DelRoomNo = append(matchUser.DelRoomNo, matchRoomInfo.RoomNo)

		//匹配用户 顺便过滤已经退出的用户
		for k, _ := range matchRoomInfo.UserInfoArr {
			uInfo := matchRoomInfo.UserInfoArr[k]
			//检查一下用户维度的数据

			//tavernUserStatus, err := table.GetUsersRoomByUid(uInfo.UserID)
			//if err != nil {
			//	global.GVA_LOG.Error("DealMatchIng4User GetUsersRoomByUid: %v %v", zap.Error(err))
			//	continue
			//}

			//if tavernUserStatus.RoomNo !=  matchRoomInfo.RoomNo {
			//	global.GVA_LOG.Infof("DealMatchIng2User 用户{%v}在房间中 ", uInfo.UserID)
			//	continue
			//}

			if len(matchUser.UserInfoArr) < match4 {
				matchUser.UserInfoArr = append(matchUser.UserInfoArr, uInfo)
			}

			if len(matchUser.UserInfoArr) == match4 {
				newMatchUser := MatchGroupRoomInfo{
					RoomNo:      "",
					DelRoomNo:   matchUser.DelRoomNo,
					UserInfoArr: matchUser.UserInfoArr,
				}
				matchRoomUserArr = append(matchRoomUserArr, &newMatchUser)
				//重置
				matchUser = &MatchGroupRoomInfo{}
			}

		}
	}

	trMgr.MatchGroupStart(matchRoomUserArr) //单排
	global.GVA_LOG.Infof("DealMatchIng1User 单排 {%v}人场匹配 结束...", match4)
}

func (trMgr *roomManager) DealMatchIng2User() {
	match4 := 4
	matchLimits := 10
	//groupLimits := 5

	//获取双排人数
	matchIngRoom2User := trMgr.MatchIngRoom.MatchIngRoom2User

	if len(matchIngRoom2User) < 2 {
		//global.GVA_LOG.Infof("DealMatchIng4User 匹配人数不够")
		return
	}

	//考虑锁的时间 太长 每次匹配最多10组

	//和快速匹配会竞争锁
	global.GVA_LOG.Infof("DealMatchIng2User 双排 {%v}人场匹配期...", match4)

	matchIngRoom2UserGroup := []*MatchIngRoomInfo{}

	if len(matchIngRoom2User) > matchLimits {
		//每次取前10个房间
		matchIngRoom2UserGroup = matchIngRoom2User[0:matchLimits]
		trMgr.MatchIngRoom.MatchIngRoom2User = matchIngRoom2User[matchLimits:]
	} else {
		if len(matchIngRoom2User)%2 == 0 {
			matchIngRoom2UserGroup = matchIngRoom2User
			trMgr.MatchIngRoom.MatchIngRoom2User = []*MatchIngRoomInfo{}
		} else {
			matchIngRoom2UserGroup = matchIngRoom2User[0 : len(matchIngRoom2User)-1]
			trMgr.MatchIngRoom.MatchIngRoom2User = matchIngRoom2User[len(matchIngRoom2User)-1:]
		}
	}

	// 前两个用户匹配成功
	matchUser := &MatchGroupRoomInfo{}

	//匹配成功 多少对
	matchRoomUserArr := []*MatchGroupRoomInfo{}

	for _, matchRoomInfo := range matchIngRoom2UserGroup {
		if len(matchRoomInfo.UserInfoArr) != 2 {
			//没有足够的匹配用户
			global.GVA_LOG.Infof("DealMatchIng4User 没有足够的匹配用户 RoomNo:%v, UserInfoMapLen:%v", matchRoomInfo.RoomNo, len(matchRoomInfo.UserInfoArr))
			continue
		}

		matchUser.DelRoomNo = append(matchUser.DelRoomNo, matchRoomInfo.RoomNo)

		//匹配用户 顺便过滤已经退出的用户
		for k, _ := range matchRoomInfo.UserInfoArr {
			uInfo := matchRoomInfo.UserInfoArr[k]
			//检查一下用户维度的数据

			//tavernUserStatus, err := table.GetUsersRoomByUid(uInfo.LikeUserId)
			//if err != nil {
			//	global.GVA_LOG.Error("DealMatchIng4User GetUsersRoomByUid: %v %v", zap.Error(err))
			//	continue
			//}

			//if tavernUserStatus.RoomNo !=  matchRoomInfo.RoomNo {
			//	global.GVA_LOG.Infof("DealMatchIng2User 用户{%v}在房间中 ", uInfo.LikeUserId)
			//	continue
			//}

			if len(matchUser.UserInfoArr) < match4 {
				matchUser.UserInfoArr = append(matchUser.UserInfoArr, uInfo)
			}

			if len(matchUser.UserInfoArr) == match4 {
				newMatchUser := MatchGroupRoomInfo{
					RoomNo:      "",
					DelRoomNo:   matchUser.DelRoomNo,
					UserInfoArr: matchUser.UserInfoArr,
				}
				matchRoomUserArr = append(matchRoomUserArr, &newMatchUser)
				//重置
				matchUser = &MatchGroupRoomInfo{}
			}

		}
	}

	trMgr.MatchGroupStart(matchRoomUserArr) //双排
	global.GVA_LOG.Infof("DealMatchIng2User 双排  {%v}人场匹配 结束...", match4)
}

// MatchGroupStart 匹配组开始游戏
func (trMgr *roomManager) MatchGroupStart(matchUserArr []*MatchGroupRoomInfo) {
	matchUserLimit := 4 //匹配人数限制
	turnNum := 5        //匹配轮数

	//每一对都启用一个房间
	global.GVA_LOG.Infof("匹配组开始游戏 MatchGroupStart  本次一共{%v}组", len(matchUserArr))

	for _, matchUserGroup := range matchUserArr {

		//1 先创建对局空间
		roomSpaceInfo := GetRoomSpace()
		//添加对局用户
		userInfoArr := matchUserGroup.UserInfoArr

		if len(userInfoArr) != matchUserLimit {
			matchUserGroupMarshal, _ := json.Marshal(matchUserGroup)
			global.GVA_LOG.Infof("MatchGroupStart 加入房间用户数量不对%v", string(matchUserGroupMarshal))
			continue
		}

		for k, _ := range userInfoArr {
			userInfo := userInfoArr[k]
			userInfo.UserExt.RoomNo = ""
			roomSpaceInfo.ComRoomSpace.AddUserInfos(userInfo.UserID, userInfo)

			//机器人直接加入确认里面
			if userInfo.UserIsRobot() {
				roomSpaceInfo.AddLoadComps(userInfo.UserID, userInfo)
			}
		}

		//todo 没有押注相关信息
		//roomSpaceInfo.MemeRoomConfig = tavernRoomConfig

		//给用户创建房间 并发送游戏开始的广播
		rNo := uuid.New().String()
		rName := userInfoArr[0].Nickname + "'s" + " lobby"
		memeRoom := table.NewAnimalPartyRoom(userInfoArr[0].UserID, userInfoArr[0].UserID, rNo, rName, "匹配房间", "", table.TavernRoomOpen, table.RoomTypeMatch, 0, 0, turnNum, matchUserLimit)
		err := table.CreateMemeRoom(memeRoom)
		if err != nil {
			global.GVA_LOG.Error("MatchGroupStart:{%v},roomInfo:%v", zap.Error(err), zap.Any("tavernRoom", memeRoom.RoomNo))
			continue
		}
		roomSpaceInfo.RoomInfo = memeRoom

		//房间用户列表
		roomUserList, err := roomSpaceInfo.MatchSuccUserList()
		if err != nil {
			memeRoom.IsOpen = table.RoomStatusAbnormal
			err := table.SaveMemeRoom(memeRoom)
			global.GVA_LOG.Error("MatchGroupStart", zap.Error(err), zap.Any("tavernRoom", memeRoom.RoomNo))
			continue
		}

		//游戏开始
		roomSpaceInfo.ComRoomSpace.IsStartGame = true

		//更新数据库 房间状态
		roomSpaceInfo.RoomInfo.IsOpen = table.RoomStatusIng
		err = table.SaveMemeRoom(roomSpaceInfo.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("MatchGroupStart JoinRoom ", zap.Error(err))
			return
		}

		//游戏每 小轮状态 游戏开始
		roomSpaceInfo.ComRoomSpace.ChangeGameState(EnGameStartIng)

		//添加到全局房间管理器
		SlotRoomManager.AddRoomSpace(memeRoom.RoomNo, roomSpaceInfo)

		//每个小房间是一个chan
		go roomSpaceInfo.Start()

		netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_matchStart), config.SlotServer)
		//发送广播通知
		msgData := models.MatchSuccResp{
			ProtoNum:  strconv.Itoa(int(pbs.Meb_matchStart)), //快速匹配成功协议 游戏开始
			Timestamp: time.Now().Unix(),
			RoomCom: models.RoomCom{
				UserId:       roomSpaceInfo.RoomInfo.UserId,
				RoomNo:       roomSpaceInfo.RoomInfo.RoomNo,
				RoomId:       roomSpaceInfo.RoomInfo.ID,
				Turn:         1,
				RoomName:     roomSpaceInfo.RoomInfo.Name,
				Status:       roomSpaceInfo.RoomInfo.IsOpen,
				UserNumLimit: roomSpaceInfo.RoomInfo.UserNumLimit,
				RoomType:     int(roomSpaceInfo.RoomInfo.RoomType),
				RoomLevel:    int(roomSpaceInfo.RoomInfo.RoomLevel),
			},
			RoomUserList: roomUserList,
		}

		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp.Content = responseHeadByte

		global.GVA_LOG.Infof("匹配成功 给客户端发消息 MatchGroupStart :%v", string(responseHeadByte))
		NatsSendAllUserMsg(roomSpaceInfo, netMessageResp) //MatchGroupStart

		roomSpaceInfo.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间

		//更新本轮数据
		roomSpaceInfo.ComRoomSpace.AddTurn()
		roomSpaceInfo.ComRoomSpace.UpdateTurnMateInfo(roomSpaceInfo.ComRoomSpace.GetTurn(), time.Now().Unix(), userInfoArr)

		delRoomNos := matchUserGroup.DelRoomNo
		global.GVA_LOG.Infof("MatchGroupStart 需要删除的房间 %v", delRoomNos)

		for _, delRoomNo := range delRoomNos {
			//删除房间用户
			beforeRoomSpaceInfo, err := SlotRoomManager.GetRoomSpace(delRoomNo)
			if err != nil {
				global.GVA_LOG.Error("MatchGroupStart 删除用户匹配前所在的房间 ", zap.Error(err))
				continue
			}

			//重置
			beforeRoomSpaceInfo.ComRoomSpace.UserInfos = map[string]*models.UserInfo{}
			beforeRoomSpaceInfo.ComRoomSpace.IsMatchClear = true

			//更新数据库数据
			//从数据库里面删除
			err = table.DelRoomUsersByRoomNo(delRoomNo)
			if err != nil {
				global.GVA_LOG.Error("MatchGroupStart DelRoomUsersByRoomNo err ", zap.Error(err))
			}

			err = table.DelMemeRoom(delRoomNo)
			if err != nil {
				global.GVA_LOG.Error("MatchGroupStart DelMemeRoom err ", zap.Error(err))
			}

		}
	}
}
