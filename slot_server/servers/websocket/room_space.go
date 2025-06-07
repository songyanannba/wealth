package websocket

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"runtime/debug"
	"slot_server/lib/common"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/lib/src/dao"
	"slot_server/protoc/pbs"
	"strconv"
	"sync"
	"time"
)

type MemeDisposeFunc func(message []byte, trs *RoomSpace) (resMessage []byte, err error)
type GameStateFunc func(trs *RoomSpace)

type RoomSpace struct {
	//房间信息
	RoomInfo *table.AnimalPartyRoom
	//通用字段
	ComRoomSpace *ComRoomSpace

	//最外圈的动物排序 这里是固定的 需要旋转特殊处理
	AnimalConfigs []*AnimalConfig

	//颜色排序 每轮开始前变化
	ColorConfigs []*ColorConfig

	//押大小配置
	BigOrSmallConfig []*BigOrSmallConfig

	//加载确认
	LoadComps map[string]*models.UserInfo

	//是否全部加载完成
	IsAllLoadComps bool

	//方法集合
	FuncMapMutex *sync.RWMutex
	FuncMap      map[string]MemeDisposeFunc

	//方法集合
	GameStateFuncMapMutex *sync.RWMutex

	//状态机制管理器
	GameStateMap map[GameTurnState]GameStateFunc
}

func GetRoomSpace() *RoomSpace {
	trSpace := &RoomSpace{
		RoomInfo:              &table.AnimalPartyRoom{},
		AnimalConfigs:         []*AnimalConfig{},
		FuncMap:               make(map[string]MemeDisposeFunc),
		GameStateMap:          make(map[GameTurnState]GameStateFunc),
		FuncMapMutex:          new(sync.RWMutex),
		GameStateFuncMapMutex: new(sync.RWMutex),
		LoadComps:             make(map[string]*models.UserInfo),
		IsAllLoadComps:        false,
	}
	trSpace.ComRoomSpace = GetComRoomSpace()

	return trSpace
}

func (trs *RoomSpace) Start() {
	//注册函数
	trs.InItFunc()

	//状态机函数
	trs.InItTurnStateFunc()

	//定时器
	//打印一些存活房间的信息
	sTimer := time.NewTicker(time.Second * 20)
	defer sTimer.Stop()

	//清理异常房间/结算等
	clearRobotRoomTimer := time.NewTicker(time.Second * 60)
	defer clearRobotRoomTimer.Stop()

	//房间逻辑处理 ：在预定的时间点发送广播
	//托管
	serviceTimer := time.NewTicker(time.Second * 6)
	defer serviceTimer.Stop()

	//游戏轮状态检测
	GameTurnStateTimer := time.NewTicker(time.Second * 5)
	defer GameTurnStateTimer.Stop()

	defer func() {
		if err := recover(); err != nil {
			global.GVA_LOG.Error("游戏结束 RoomSpace:", zap.Any("recover ", err), zap.Any("recover ", string(debug.Stack())))
			global.GVA_LOG.Infof("游戏结束 RoomSpace %v %v  recover roomNo:%v UserOwner%v ",
				zap.Any("recover ", err), zap.Any("recover ", string(debug.Stack())), trs.RoomInfo.RoomNo, trs.ComRoomSpace.UserOwner.UserID)
		}
		global.GVA_LOG.Infof("游戏结束房间号:%v ", trs.RoomInfo.RoomNo)
	}()

	for {
		select {
		case <-sTimer.C:
			global.GVA_LOG.Infof("对局房间 ======###====== 存活:%v", trs.RoomInfo.RoomNo)
		case <-trs.ComRoomSpace.Close:
			global.GVA_LOG.Infof("游戏结束 关闭房间 ============ %v 中", trs.RoomInfo)
			//没有开始的房间
			if trs.RoomInfo.IsOpen == table.RoomStatusIng {
				//回收资源 开始下一期
				trs.AnalyzeFinalIncome()
			}

			//trs.AnalyzeFinalIncome()

			//关闭资源
			close(trs.ComRoomSpace.Broadcast)
			close(trs.ComRoomSpace.ReceiveMsg)

			//清理房间管理器的房间
			//_, ok := SlotRoomManager.Rooms[trs.RoomInfo.RoomNo]
			//if ok {
			//	SlotRoomManager.DelRoom(trs.RoomInfo.RoomNo)
			//}
			return
		case <-GameTurnStateTimer.C:
			//游戏轮状态检测
			trs.GameTurnStateCheck()
		case <-serviceTimer.C:
			//逻辑包含：超时自动出牌
			//trs.AnalyzeRoom()
		case message := <-trs.ComRoomSpace.ReceiveMsg:
			// 接收待处理的事件
			global.GVA_LOG.Infof("ReceiveMsg message %v ", message)
			//trs.ProcessData(message)
		case <-clearRobotRoomTimer.C:
			//如果长时间 用户不进入下一轮 判定为见好就收
			global.GVA_LOG.Infof("clearRobotRoomTimer 房间编号 %v", trs.RoomInfo.RoomNo)
			if trs.ClearRoom() {
				global.GVA_LOG.Infof("超过2分钟没有进入房间下一轮 执行清理房间逻辑 ClearRoom 成功")
				return
			}
		}
	}
}

func (trs *RoomSpace) ProcessData(message []byte) {
	global.GVA_LOG.Infof("tavernRoomSpace 处理数据:%v", string(message))
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("tavernRoomSpace 处理数据 stop", zap.Any("", r))
		}
	}()

	comMsg := &models.ComMsg{}
	err := json.Unmarshal(message, comMsg)
	if err != nil {
		global.GVA_LOG.Error("tavernRoomSpace 处理数据 json Unmarshal", zap.Any("err", err))
		return
	}

	var data []byte
	// 采用 map 注册的方式
	if value, ok := trs.GetHandlers(comMsg.MsgId); ok {
		data, err = value(comMsg.Data, trs)
	} else {
		global.GVA_LOG.Error("tavernRoomSpace 处理数据 路由不存在", zap.Any("MsgId", comMsg.MsgId))
		return
	}

	global.GVA_LOG.Infof("RoomNo{%v},处理 comMsg.MsgId %v 返回数据data:%v ", trs.RoomInfo.RoomNo, comMsg.MsgId, string(data))
	return
}

func (trs *RoomSpace) Register(key string, memeDisposeFunc func(message []byte, trs *RoomSpace) (resMessage []byte, err error)) {
	trs.FuncMapMutex.Lock()
	defer trs.FuncMapMutex.Unlock()
	trs.FuncMap[key] = memeDisposeFunc
	return
}

func (trs *RoomSpace) GetHandlers(key string) (value MemeDisposeFunc, ok bool) {
	trs.FuncMapMutex.Lock()
	defer trs.FuncMapMutex.Unlock()
	value, ok = trs.FuncMap[key]
	return
}

func (trs *RoomSpace) GameTurnStateCheck() {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	var (
		currTime          = helper.LocalTime().Unix()
		gState            = trs.ComRoomSpace.GetGameState()
		currGameStartTime = trs.ComRoomSpace.GetGameStartTime()
	)

	//如果当前没有押注的用户 不往下执行
	if len(trs.ComRoomSpace.UserInfos) < 1 {
		global.GVA_LOG.Infof("没有人押注")
		trs.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间
		return
	}
	global.GVA_LOG.Infof("GameTurnStateCheck 执行,currTime:%v 房间编号:%v ,gState:%v ,currGameStartTime:%v", currTime, trs.RoomInfo.RoomNo, gState, currGameStartTime)

	if gState == BetIng && trs.ComRoomSpace.APRobotActionCount < 1 {
		trs.APRobotAction()
	}

	//押注期间
	if gState == BetIng && currTime-currGameStartTime <= BetIngPeriod {
		return
	}

	if gState == BetIng && currTime-currGameStartTime > BetIngPeriod {
		//告诉客户端开始计算
		trs.ComRoomSpace.ChangeGameState(EnWheelAnimalPartyCalculateExec)
	}

	trs.ExecProcessTurnStateFunc(trs.ComRoomSpace.GetGameState())

}

// GameTurnStateCheck1 每3秒检查一下状态变化
func (trs *RoomSpace) GameTurnStateCheck1() {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	var (
		currTime = time.Now().Unix()
	)
	global.GVA_LOG.Infof("GameTurnStateCheck 执行,currTime:%v 房间编号 %v", currTime, trs.RoomInfo.RoomNo)

	//不同状态 触发不同方法

	gState := trs.ComRoomSpace.GetGameState()
	global.GVA_LOG.Infof("GameTurnStateCheck 游戏状态%v", gState)

	//自动判断执行逻辑
	trs.ExecProcessTurnStateFunc(gState)

}

func (trs *RoomSpace) GetNewAnimalConfigsBySeat(seat int, animalConfigs []*AnimalConfig) *AnimalConfig {
	res := &AnimalConfig{}
	for _, animalConfig := range animalConfigs {
		if animalConfig.Seat == seat {
			res = animalConfig
			break
		}
	}
	return res
}

func (trs *RoomSpace) GetBigOrSmallConfigsBySeat(seat int, bigOrSmallConfigs []*BigOrSmallConfig) *BigOrSmallConfig {
	res := &BigOrSmallConfig{}
	for _, bigOrSmallConfig := range bigOrSmallConfigs {
		if bigOrSmallConfig.Seat == seat {
			res = bigOrSmallConfig
			break
		}
	}
	return res
}

func (trs *RoomSpace) GetColorConfigsBySeat(seat int) *ColorConfig {
	res := &ColorConfig{}
	for _, colorConfigs := range trs.ColorConfigs {
		if colorConfigs.Seat == seat {
			res = colorConfigs
			break
		}
	}
	return res
}

func (trs *RoomSpace) AnalyzeFinalIncome() {
	//var tavernRoomUsers []models.MemeRoomUser

	//继续游戏的房间
	nNo := uuid.New().String()
	//trs.IntoNextRoom(nNo)
	newRoomSpace := trs.IntoNextRoom(nNo)

	if newRoomSpace == nil {
		return
	}

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.ProtocNum_ColorSortMsg), config.SlotServer)
	//发送广播
	msgData := &pbs.ColorSortMsg{}
	for _, colorConfig := range newRoomSpace.ColorConfigs {
		msgData.ColorConfig = append(msgData.ColorConfig, &pbs.ColorConfig{
			Seat:    int32(colorConfig.Seat),
			ColorId: int32(colorConfig.ColorId),
		})
	}

	//给用户消息
	global.GVA_LOG.Infof("游戏结束 下一轮开始: %v", msgData)
	responseHeadByte, _ := proto.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAimUserMsg(trs, netMessageResp, "")

}

func (trs *RoomSpace) IntoNextRoom(roomNo string) *RoomSpace {
	trs.RoomInfo.IsOpen = table.RoomStatusStop
	err := table.SaveMemeRoom(trs.RoomInfo)
	if err != nil {
		global.GVA_LOG.Error("LeaveRoom SaveMemeRoom", zap.Error(err))
	}

	period := helper.Itoa(helper.Atoi(trs.RoomInfo.Period) + 1)

	roomInfo := table.NewAnimalPartyRoom("", "", roomNo, trs.RoomInfo.Name, "第"+period+"期", helper.Itoa(helper.Atoi(trs.RoomInfo.Period)+1),
		table.TavernRoomOpen, trs.RoomInfo.RoomType, trs.RoomInfo.RoomLevel, table.RoomClassMatch, trs.RoomInfo.RoomTurnNum, trs.RoomInfo.UserNumLimit)
	roomInfo.IsGoOn = 1
	err = table.CreateMemeRoom(roomInfo)
	if err != nil {
		global.GVA_LOG.Error("IntoNextRoom CreateMemeRoom:{%v},roomInfo:%v", zap.Error(err))
		return nil
	}

	//1 先创建对局空间
	roomSpaceInfo := GetRoomSpace()
	roomSpaceInfo.RoomInfo = roomInfo
	roomSpaceInfo.ComRoomSpace.ChangeGameState(BetIng)
	roomSpaceInfo.ComRoomSpace.IsStartGame = true
	roomSpaceInfo.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间
	roomSpaceInfo.ColorConfigs = GetColorWheel()

	//2 添加到全局房间管理器
	SlotRoomManager.ReRoomSpace(roomInfo.Name, roomSpaceInfo)

	//每个小房间是一个 协成
	go roomSpaceInfo.Start()

	return roomSpaceInfo
}

func (trs *RoomSpace) ClearRoom() bool {
	trs.ComRoomSpace.Sync.Lock()
	defer trs.ComRoomSpace.Sync.Unlock()

	//turnMate := trs.ComRoomSpace.TurnMateInfo
	// CountdownTime  == 0 是初始化
	//if turnMate.CountdownTime <= 0 {
	//	return false
	//}

	currTime := helper.LocalTime().Unix()
	//超过2分钟没有进入房间 清理房间
	global.GVA_LOG.Infof("tavernRoomSpace ClearRoom 执行,currTime:%v ", currTime)

	//是否都掉线了
	//notHeartbeat := 0
	//for _, uClt := range trs.ComRoomSpace.UserInfos {
	//	if int64(uClt.HeartbeatTime)+30 < currTime {
	//		notHeartbeat += 1
	//	}
	//}

	//是否都掉线了
	//isAllNotHeartbeat := false
	//if notHeartbeat == len(trs.ComRoomSpace.UserInfos) {
	//	isAllNotHeartbeat = true
	//}
	//
	//isAllNotHeartbeat = false //todo 先不管心跳

	//if !isAllNotHeartbeat {
	//	//当前时间 和客户端上次请求的存活时间大于15分钟 释放房间
	//	if currTime-trs.ComRoomSpace.CurrentOpTime < RoomAlive {
	//		return false
	//	}
	//}

	if currTime-trs.ComRoomSpace.CurrentOpTime < RoomAlive {
		return false
	}

	close(trs.ComRoomSpace.Broadcast)
	close(trs.ComRoomSpace.ReceiveMsg)

	global.GVA_LOG.Infof("清理房间 ClearRoom DelRoom,currTime:%v RoomNo:%v", currTime, trs.RoomInfo.RoomNo)

	//如果被清理的房间没有开始，并且里面加入的有人 需要返回积分
	if !trs.ComRoomSpace.IsStartGame {
		//for _, uInfo := range trs.ComRoomSpace.UserInfos {
		//	//是否已经扣减过积分
		//	payPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", trs.MemeRoomConfig.Bet+trs.MemeRoomConfig.AdmissionPrice), 64)
		//	logic.ReturnUserCoin(uInfo.UserID, trs.RoomInfo.RoomNo, payPrice, "清理房间,未开始的房间返还用户积分")
		//}
	}

	//清理房间管理器的房间
	_, ok := SlotRoomManager.Rooms[trs.RoomInfo.RoomNo]
	if ok {
		global.GVA_LOG.Infof("清理房间 ClearRoom DelRoom,currTime:%v RoomNo:%v", currTime, trs.RoomInfo.RoomNo)

		//更新房间状态
		trs.RoomInfo.IsOpen = table.RoomStatusAbnormal

		if trs.ComRoomSpace.IsMatchClear {
			trs.RoomInfo.IsOpen = table.RoomStatusAbnormal
		}

		//更新数据库 房间状态
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("ClearRoom DelRoomByTime LeaveRoomNotStartGame", zap.Error(err))
		}

		SlotRoomManager.DelRoom(trs.RoomInfo.RoomNo)
	}

	return true
}

func (trs *RoomSpace) MatchSuccUserList() ([]models.MemeRoomUser, error) {
	if trs.ComRoomSpace.TurnMateInfo.Turn != 0 {
		return nil, errors.New("MatchSuccUserList 不能开始")
	}
	roomUserLists := make([]models.MemeRoomUser, 0)

	// 快速匹配 直接进入游戏
	// 加入房间
	// 就绪
	// 房主开始游戏
	// 加载完成 | 客户端调用
	index := 0

	//给用户发消息,游戏开始  并告诉谁先开始说话
	for _, userInfo := range trs.ComRoomSpace.UserInfos {
		if trs.RoomInfo.UserId == userInfo.UserID {
			userInfo.SetUserIsOwner(true)
		}

		//添加房间用户
		isOwner := 0
		if userInfo.GetUserIsIsOwner() {
			isOwner = 1
		}

		tavernRoomUsers := table.NewRoomUsers(userInfo.UserID, trs.RoomInfo.RoomNo, userInfo.Nickname, index+1, 1, isOwner, 1, 1)
		err := table.CreateRoomUsers(tavernRoomUsers)
		if err != nil {
			global.GVA_LOG.Error("MatchSuccUserList CreateRoomUsers", zap.Error(err), zap.Any(".RoomNo,", trs.RoomInfo.RoomNo))
			return roomUserLists, err
		}

		if !userInfo.UserIsRobot() {
			//用户维度状态 ｜ 快速匹配 | 创建房间
			tavernUserRoomData := table.NewUserRoom(tavernRoomUsers.UserId, tavernRoomUsers.RoomNo, tavernRoomUsers.Nickname, tavernRoomUsers.IsLeave, tavernRoomUsers.IsKilled,
				tavernRoomUsers.IsOwner, tavernRoomUsers.Turn, tavernRoomUsers.Seat, tavernRoomUsers.IsRobot, tavernRoomUsers.IsReady)
			err = dao.CreateOrUpdateUsersRoom(tavernUserRoomData)
			if err != nil {
				global.GVA_LOG.Error("MatchSuccUserList", zap.Any("CreateOrUpdateUsersRoom", err), zap.Any(".RoomNo,", trs.RoomInfo.RoomNo))
			}
		}

		//设置用户座位
		userInfo.UserProperty.Seat = tavernRoomUsers.Seat
		userInfo.UserProperty.Turn = tavernRoomUsers.Turn
		userInfo.UserExt.RoomNo = trs.RoomInfo.RoomNo

		//就绪
		userInfo.SetUserIsReady(1)

		roomUser := models.MemeRoomUser{
			UserID:       userInfo.UserID,
			Nickname:     userInfo.Nickname,
			Turn:         tavernRoomUsers.Turn,
			IsOwner:      userInfo.UserProperty.IsOwner,
			IsReady:      userInfo.UserProperty.IsReady,
			Seat:         tavernRoomUsers.Seat,
			UserLimitNum: userInfo.UserProperty.UserLimitNum,
			UserCards:    models.UserCartState{},
			Bet:          0,
			IsRobot:      userInfo.UserProperty.IsRobot,
		}

		roomUserLists = append(roomUserLists, roomUser)
		index++
	}

	return roomUserLists, nil
}

// LeaveRoomNotStartGame 没开始游戏前离开房间
func (trs *RoomSpace) LeaveRoomNotStartGame(userInfo *models.UserInfo) uint32 {

	userId := userInfo.UserID

	if userInfo.UserID == trs.ComRoomSpace.UserOwner.UserID {
		//todo 返回值优化
		return trs.OwnerLeaveRoomNotStartGame(userInfo)
	} else {
		//非房主
		//从数据库里面删除
		err := table.DelRoomUsers(userInfo.UserExt.RoomNo, userId)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom SaveMemeRoom CreateRoomUsers", zap.Error(err))
			return common.ModelDeleteError
		}

		//更新老房主用户状态
		updateMap := dao.MakeUpdateData("room_no", "")
		updateMap["is_leave"] = 1
		updateMap["is_owner"] = 0
		updateMap["seat"] = 0
		dao.UpdateUsersRoomRoomNo(userId, updateMap)

		//获取房间人数  发送广播 离开房间
		msgData := models.LeaveRoomMsg{
			ProtoNum: strconv.Itoa(int(pbs.Meb_leaveRoom)),
			UserId:   userId,
			RoomNo:   userInfo.UserExt.RoomNo,
		}

		responseHeadByte, _ := json.Marshal(msgData)
		//uIdInt, _ := strconv.Atoi(userId)
		netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_leaveRoom), config.SlotServer)
		netMessageResp.Content = responseHeadByte

		global.GVA_LOG.Infof("LeaveRoom 离开房间的广播: %v", string(responseHeadByte))
		NatsSendAllUserMsg(trs, netMessageResp) //LeaveRoom

		//从房间里面把用户删除
		trs.ComRoomSpace.DelUserInfoAndUserClient(userId)

		if trs.RoomInfo.IsOpen == table.RoomStatusFill {
			trs.RoomInfo.IsOpen = table.RoomStatusOpen
			//更新数据库 房间状态
			err := table.SaveMemeRoom(trs.RoomInfo)
			if err != nil {
				global.GVA_LOG.Error("LeaveRoom", zap.Error(err))
				return common.ModelAddError
			}
		}
	}
	return common.OK
}

func (trs *RoomSpace) OwnerLeaveRoomNotStartGame(userInfo *models.UserInfo) uint32 {
	userId := userInfo.UserID

	//游戏没开始 房主离开房间 ；如果房间就房主一个人 直接解散
	records, err := table.GetRoomUsers(userInfo.UserExt.RoomNo)
	if err != nil {
		global.GVA_LOG.Error("LeaveRoom GetTavernRoomUsers", zap.Error(err))
		return common.ModelDeleteError
	}

	//更新用户维度数据 老房主用户状态
	updateMap := dao.MakeUpdateData("room_no", "")
	updateMap["is_leave"] = 1
	updateMap["is_owner"] = 0
	updateMap["seat"] = 0
	dao.UpdateUsersRoomRoomNo(userId, updateMap)

	//房间就一个人 离开房间 就是解散房间
	if len(records) == 1 && records[0].UserId == userId {

		//更新数据库 房间状态
		trs.RoomInfo.IsOpen = table.RoomStatusDissolve
		err := table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom LeaveRoomNotStartGame", zap.Error(err))
			return common.ModelAddError
		}

		netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_leaveRoom), config.SlotServer)
		//获取房间人数  发送广播 离开房间
		msgData := models.LeaveRoomMsg{
			ProtoNum: strconv.Itoa(int(pbs.Meb_leaveRoom)),
			UserId:   userId,
			RoomNo:   userInfo.UserExt.RoomNo,
		}
		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp.Content = responseHeadByte

		global.GVA_LOG.Infof("LeaveRoom 离开房间的广播: %v", string(responseHeadByte))
		NatsSendAllUserMsg(trs, netMessageResp) //LeaveRoom

		//删除房间用户
		trs.ComRoomSpace.DelRoomAllUserAndUserClient(userInfo.UserExt.RoomNo)

		//通知房间管理器 删除房间
		trs.CloseRoom(userInfo.UserExt.RoomNo, table.RoomStatusStop)

	} else {
		//如果房间还有其他人-选择最新进来的人成为房主

		//从数据库里面删除
		err := table.DelRoomUsers(userInfo.UserExt.RoomNo, userId)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom SaveMemeRoom CreateRoomUsers", zap.Error(err))
			return common.ModelDeleteError
		}

		//payPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", trs.MemeRoomConfig.Bet+trs.MemeRoomConfig.AdmissionPrice), 64)
		//logic.ReturnUserCoin(userInfo.UserID, trs.RoomInfo.RoomNo, payPrice, "房主离开房间-并选择新房主")

		//在获取一次房间用户
		records, err = table.GetRoomUsers(userInfo.UserExt.RoomNo)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom GetTavernRoomUsers", zap.Error(err))
			return common.ModelDeleteError
		}

		if len(records) <= 0 {
			//没有其他用户
			//没开始游戏的房间 房主离开 房间解散
			trs.RoomInfo.IsOpen = table.RoomStatusDissolve

			//更新数据库 房间状态
			err := table.SaveMemeRoom(trs.RoomInfo)
			if err != nil {
				global.GVA_LOG.Error(err.Error())
				return common.ModelAddError
			}

			//通知房间管理器 删除房间
			trs.CloseRoom(userInfo.UserExt.RoomNo, table.RoomStatusStop)

			global.GVA_LOG.Error("LeaveRoom GetTavernRoomUsers", zap.Error(err))
			return common.OK
		}

		//新房主
		newOwnerRoomUser := records[0]
		newOwnerUserInfo, err := trs.ComRoomSpace.GetUserInfo(newOwnerRoomUser.UserId)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom GetUserInfo NewOwnerUserInfo", zap.Error(err))
			return common.ModelDeleteError
		}

		//更新房间数据
		newOwnerUserInfo.SetUserIsOwner(true)

		trs.ComRoomSpace.UserOwner = newOwnerUserInfo

		//更新房间用户数据
		trs.RoomInfo.Owner = newOwnerUserInfo.UserID
		err = table.SaveMemeRoom(trs.RoomInfo)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom SaveMemeRoom", zap.Error(err))
			return common.ModelDeleteError
		}

		//更新新房主标识
		err = dao.UpdateRoomOwner(newOwnerUserInfo.UserID, userInfo.UserExt.RoomNo, table.BeOwner)
		if err != nil {
			global.GVA_LOG.Error("LeaveRoom UpdateTavernRoomUsers", zap.Error(err))
			return common.ModelDeleteError
		}

		//更新用户维度 房主状态
		newUpdateMap := dao.MakeUpdateData("room_no", userInfo.UserExt.RoomNo)
		newUpdateMap["is_owner"] = 1
		dao.UpdateUsersRoomRoomNo(newOwnerUserInfo.UserID, newUpdateMap)

		//发送广播 离开房间
		msgData := models.LeaveRoomMsg{
			ProtoNum:     strconv.Itoa(int(pbs.Meb_leaveRoom)),
			UserId:       userId,
			RoomNo:       userInfo.UserExt.RoomNo,
			IsOwnerLeave: true,
			NewOwner:     newOwnerUserInfo.UserID,
			Timestamp:    time.Now().Unix(),
		}

		responseHeadByte, _ := json.Marshal(msgData)
		netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_leaveRoom), config.SlotServer)
		netMessageResp.Content = responseHeadByte

		NatsSendAllUserMsg(trs, netMessageResp)

		//从房间里面把用户删除
		trs.ComRoomSpace.DelUserInfoAndUserClient(userId)

		if trs.RoomInfo.IsOpen == table.RoomStatusFill {
			//更新数据库 房间状态
			trs.RoomInfo.IsOpen = table.RoomStatusOpen
			err := table.SaveMemeRoom(trs.RoomInfo)
			if err != nil {
				global.GVA_LOG.Error("LeaveRoom", zap.Error(err))
				return common.ModelAddError
			}
		}
	}

	return common.OK
}

func (trs *RoomSpace) CloseRoom(roomNo string, typ int) {
	closeInfo := models.CloseRoom{
		RoomNo: roomNo,
		IsStop: true,
		//Type:   typ,
	}
	global.GVA_LOG.Infof("质疑后 游戏结束, closeInfo:%v userInfo:%v", closeInfo)
	closeInfoMs, _ := json.Marshal(closeInfo)
	SlotRoomManager.CommonRoomManager.CloseRoom <- closeInfoMs
}

// IsAllCompleted 是否都确认了
func (trs *RoomSpace) IsAllCompleted() bool {
	return trs.IsAllLoadComps
}

func (trs *RoomSpace) AddLoadComps(userId string, userInfo *models.UserInfo) {
	//调用前有加锁
	_, ok := trs.LoadComps[userId]
	if !ok {
		trs.LoadComps[userId] = userInfo
	}
}
