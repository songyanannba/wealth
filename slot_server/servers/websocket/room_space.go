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
	"slot_server/lib/src/logic"
	"slot_server/protoc/pbs"
	"sort"
	"strconv"
	"sync"
	"time"
)

type MemeDisposeFunc func(message []byte, trs *RoomSpace) (resMessage []byte, err error)
type GameStateFunc func(trs *RoomSpace)

type RoomSpace struct {
	//房间信息
	RoomInfo *table.AnimalPartyRoom

	AnimalConfigs []*AnimalConfig

	ColorConfigs []*ColorConfig

	//房间配置
	MemeRoomConfig *table.MemeRoomConfig

	//通用字段
	ComRoomSpace *ComRoomSpace

	//基础问题
	RoomIssueConfig []*table.MbIssueConfig

	//当前的卡版本
	CardVersion []int

	//基础卡
	RoomBaseCard []*table.MbCardConfig

	//版本卡
	RoomVersionCard map[int][]*table.MbCardConfig

	//加载确认
	LoadComps map[string]*models.UserInfo

	//是否全部加载完成
	IsAllLoadComps bool

	//游戏状态
	GameState GameState

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
		MemeRoomConfig:        &table.MemeRoomConfig{},
		RoomBaseCard:          make([]*table.MbCardConfig, 0),
		RoomVersionCard:       make(map[int][]*table.MbCardConfig),
		FuncMap:               make(map[string]MemeDisposeFunc),
		GameStateMap:          make(map[GameTurnState]GameStateFunc),
		FuncMapMutex:          new(sync.RWMutex),
		GameStateFuncMapMutex: new(sync.RWMutex),
		LoadComps:             make(map[string]*models.UserInfo),
		IsAllLoadComps:        false,
		CardVersion:           make([]int, 0),
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
		currTime = helper.LocalTime().Unix()
	)

	//不同状态 触发不同方法

	//gState := trs.ComRoomSpace.GetGameState()
	//global.GVA_LOG.Infof("GameTurnStateCheck 游戏状态%v", gState)

	//1 -30 押注期间
	//30-60 理论上的计算期间
	//计算期间不让押注
	//计算完 推送计算结果
	//计算完 自动开始下一轮

	//如果当前没有押注的用户 不往下执行
	if len(trs.ComRoomSpace.UserInfos) <= 0 {
		global.GVA_LOG.Infof("没有人押注")
		trs.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间
		return
	}

	gState := trs.ComRoomSpace.GetGameState()
	currGameStartTime := trs.ComRoomSpace.GetGameStartTime()
	global.GVA_LOG.Infof("GameTurnStateCheck 执行,currTime:%v 房间编号:%v ,gState:%v ,currGameStartTime:%v", currTime, trs.RoomInfo.RoomNo, gState, currGameStartTime)

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

// GameTurnStateCheck 每3秒检查一下状态变化
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

	//自动判断 进入下一个状态
	//trs.ExecAutoNextTurnState(gState)

	//机器人用户行为
	trs.RobotAction()

	//超时自动出牌
	trs.OutTimePlayHand()
}

//func (trs *RoomSpace) GetAnimalConfigsBySeat(seat int) *AnimalConfig {
//	res := &AnimalConfig{}
//	for _, animalConfigs := range trs.AnimalConfigs {
//		if animalConfigs.Seat == seat {
//			res = animalConfigs
//			break
//		}
//	}
//	return res
//}

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

// OutTimePlayHand 超时出牌（托管）
func (trs *RoomSpace) OutTimePlayHand() {
	//用户掉线后 后台进程进行托管 帮用户出牌
	if !trs.IsAllLoadComps {
		//超时load
		gameStartTime := trs.ComRoomSpace.GetGameStartTime()
		if gameStartTime > 0 && gameStartTime-helper.LocalTime().Unix() < 0 {
			//系统加载
			for _, userInfo := range trs.ComRoomSpace.UserInfos {
				trs.AddLoadComps(userInfo.UserID, userInfo)
			}

			//第一次全部加载完成
			if len(trs.LoadComps) == trs.RoomInfo.UserNumLimit {
				trs.IsAllLoadComps = true
				global.GVA_LOG.Infof("LoadCompleted 房间{%v},全部加载完成，开始发牌", trs.RoomInfo.RoomNo)
				//全部加载改变状态 应该是 游戏中 -> 加载中；过度
				trs.ComRoomSpace.GameStateTransition(EnGameStartIng, EnLoadExec)
			}
		}
		return
	}

	//1 获取当前的游戏状态 如果大于某个阶段的的执行时间 执行委托操作
	//委托阶段有哪些
	// 1 出牌
	// 2 点赞

	//是否托管 系统帮忙执行
	var isSysExec bool
	//0=游戏阶段
	var gameStatus int

	gameStatus, isSysExec = trs.ServerSimplifyGetStateAndTime()

	if !isSysExec {
		return
	}

	//1 遍历用户，获取出牌用户出牌倒计时 如果大于40秒没出牌
	for _, userInfo := range trs.ComRoomSpace.UserInfos {

		//  随牌阶段 + 出牌阶段 托管
		if gameStatus == int(CliOutCard) {
			//是否出牌
			cards := trs.ComRoomSpace.GetUserOutEdCards(userInfo.UserID)
			if len(cards) > 0 {
				//过滤已经出过牌的用户
				continue
			}

			currCards, err := trs.ComRoomSpace.GetCurrCard(userInfo.UserID)
			if err != nil {
				global.GVA_LOG.Error("超时出牌（托管） 获取前一个用户当前的牌 错误", zap.Error(err))
			}
			if len(currCards) <= 0 {
				global.GVA_LOG.Error("超时出牌（托管）用户手里没牌", zap.Error(err))
				continue
			}

			var reqCards []*models.Card
			reqCards = append(reqCards, currCards[0])
			global.GVA_LOG.Infof("超时出牌 出牌 reqCards %v ,currCards  %v ,userInfo.UserID  %v", reqCards, currCards, userInfo.UserID)
			trs.OutCart(reqCards, currCards, userInfo.UserID) //托管出牌
			global.GVA_LOG.Infof("超时出牌（托管）结束,用户:%v", userInfo.UserID)
		}

		//点赞托管
		if gameStatus == int(CliLikePage) {
			userId := userInfo.UserID
			likeUserId := ""
			likeCard := models.LikeCard{}
			outCards := make([]*models.Card, 0)
			likeCards := make([]*models.Card, 0)

			//那个用户没点赞
			likeUserInfo := trs.ComRoomSpace.GetLikeUserInfo(userInfo.UserID)
			if len(likeUserInfo) > 0 {
				//该用户已经给别人点过赞
				global.GVA_LOG.Infof("超时点赞 托管 该用户已经给别人点过赞 userID %v", userInfo.UserID)
				continue
			}

			//每个人出一个牌 取第一个就行
			outCards = trs.ComRoomSpace.GetUserOutEdCardExcludeUser(userInfo.UserID)
			sort.Slice(outCards, func(i, j int) bool {
				return outCards[i].Level > outCards[j].Level
			})

			if len(outCards) <= 0 {
				continue
			}

			//找到牌等级最高的一张 点赞
			//每轮每次只出一张牌
			isOutCard := false
			for k, outCard := range outCards {
				if k == 1 {
					//只取第0个
					break
				}
				likeCard = models.LikeCard{
					CardId:     outCard.CardId,
					LikeUserId: outCard.UserID,
					Level:      outCard.Level,
					AddRate:    outCard.AddRate,
				}
				isOutCard = true
				likeUserId = outCard.UserID
				likeCards = append(likeCards, outCard)
			}

			if !isOutCard {
				continue
			}
			trs.DoLikeCard(userId, likeUserId, likeCard, likeCards)
		}
	}
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

		tavernRoomUsers := table.NewRoomUsers(userInfo.UserID, trs.RoomInfo.RoomNo, userInfo.Nickname, index+1, 1, isOwner, trs.MemeRoomConfig.Bet, 1)
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
			Bet:          trs.MemeRoomConfig.Bet,
			IsRobot:      userInfo.UserProperty.IsRobot,
		}

		roomUserLists = append(roomUserLists, roomUser)
		index++

		//提前扣减积分
		//err = logic.AddUserScore(userInfo.UserID, trs.RoomInfo.RoomNo, 0, payPrice, models.TavernStoryQuickMatch, "快速匹配-扣减积分")
		//if err != nil {
		//	payFailUser = append(payFailUser, userInfo)
		//	global.GVA_LOG.Error("MatchSuccUserList AddUserScore 快速匹配扣减积分", zap.Error(err), zap.Any(".RoomNo,", trs.RoomInfo.RoomNo))
		//} else {
		//	havePayUser = append(havePayUser, userInfo)
		//}
	}

	//说明存在扣减积分失败的情况
	//if len(payFailUser) > 0 {
	//	for _, userInfo := range havePayUser {
	//		logic.ReturnUserCoin(userInfo.UserID, trs.RoomInfo.RoomNo, payPrice, "快速匹配返还积分")
	//	}
	//	//用户回归匹配池子 todo
	//	//这里先不做处理
	//	//SlotRoomManager.MatchIngUser.MatchIng2UsersMap[roomType][roomLevel] = newMatchIng2Users
	//	return errors.New("MatchSuccUserList 存在扣减积分失败情况 不能开始游戏")
	//}

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

		//payPrice, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", trs.MemeRoomConfig.Bet+trs.MemeRoomConfig.AdmissionPrice), 64)
		//logic.ReturnUserCoin(userInfo.UserID, trs.RoomInfo.RoomNo, payPrice, "非房主离开房间")

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

// ReLoadCompleted 断线重连加载
func (trs *RoomSpace) ReLoadCompleted(userId string, currIssue *models.Issue) {
	rs := trs.ComRoomSpace
	roomUserLists := make([]models.MemeRoomUser, 0)

	//知道每个用户的牌情况
	allUserCardState := make([]models.UserCartState, 0)

	//游戏阶段 阶段倒计时 秒
	gameStatus, timeDown := trs.CurrGameTurnStateAndDownTime()
	global.GVA_LOG.Infof("CurrGameTurnStateAndDownTime 断线重连加载 OperateCards gameStatus %v, turnTime %v", gameStatus, timeDown)
	for _, userItem := range rs.UserInfos {
		//用户当前牌
		cards := rs.GetTurnCards(userItem.UserID)
		userCards := models.UserCartState{
			UserID:     userItem.UserID,
			OutCardNum: 4 - len(cards),
			CardNum:    len(cards),
		}

		if userId != userItem.UserID {
			allUserCardState = append(allUserCardState, userCards)
		}

		userCards.Card = cards
		memeRoomUser := models.MemeRoomUser{
			UserID:       userItem.UserID,
			Nickname:     userItem.Nickname,
			Turn:         userItem.UserProperty.Turn,
			IsLeave:      userItem.UserProperty.IsLeave,
			IsOwner:      userItem.UserProperty.IsOwner,
			IsReady:      userItem.UserProperty.IsReady,
			Seat:         userItem.UserProperty.Seat,
			UserLimitNum: userItem.UserProperty.UserLimitNum,
			WinPrice:     userItem.UserProperty.WinPrice,
			Bet:          userItem.UserProperty.Bet,
			IsRobot:      userItem.UserProperty.IsRobot,
		}

		if userId == userItem.UserID {
			memeRoomUser.UserCards = userCards
			//断线重连 需要重置托管时间 todo
			//userItem.SetOutCardCountDown(models.GetOutCardCountDownTimeInt(OutCardCountDownTimeInt))
		}
		roomUserLists = append(roomUserLists, memeRoomUser)
	}

	//点赞页面的卡
	likeCards := rs.GetCurrTurnLikeCards()

	//获取本轮房间所有用户的牌，发送给房间所有人
	allOutCard := trs.ComRoomSpace.GetAllUserOutEdCards()

	//发送广播
	msgData := models.LoadMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_loadCompleted)),
		Timestamp: time.Now().Unix(),
		RoomCom: models.RoomCom{
			Turn:         trs.ComRoomSpace.GetTurn(),
			RoomNo:       trs.RoomInfo.RoomNo,
			UserId:       trs.RoomInfo.Owner,
			RoomName:     trs.RoomInfo.Name,
			Status:       trs.RoomInfo.IsOpen, //游戏状态
			UserNumLimit: trs.RoomInfo.UserNumLimit,
			RoomType:     int(trs.RoomInfo.RoomType),
			RoomLevel:    int(trs.RoomInfo.RoomLevel),
			CurrIssue:    currIssue,
			GameStatus:   gameStatus,
			TimeDown:     timeDown,
		},
		RoomUserList:   roomUserLists,
		OtherUserCards: allUserCardState,
		LikeCards:      likeCards,
		OutCards:       allOutCard,
	}
	//给用户消息
	responseHeadByte, _ := json.Marshal(msgData)
	NatsSendAimUserMsg(trs, helper.GetNetMessage("", "", int32(pbs.Meb_loadCompleted), config.SlotServer, responseHeadByte), userId)
}

func (trs *RoomSpace) LoadCompletedFirst(userId string, currIssue *models.Issue) {
	roomUserLists, _ := trs.ComRoomSpace.UserInfoToRoomUser()
	//发送广播
	msgData := models.LoadMsg{
		ProtoNum:  strconv.Itoa(int(pbs.Meb_loadCompleted)),
		Timestamp: time.Now().Unix(),
		RoomCom: models.RoomCom{
			Turn:         trs.ComRoomSpace.GetTurn(),
			RoomNo:       trs.RoomInfo.RoomNo,
			UserId:       trs.RoomInfo.Owner,
			RoomName:     trs.RoomInfo.Name,
			Status:       trs.RoomInfo.IsOpen,
			UserNumLimit: trs.RoomInfo.UserNumLimit,
			RoomType:     int(trs.RoomInfo.RoomType),
			RoomLevel:    int(trs.RoomInfo.RoomLevel),
			CurrIssue:    currIssue,
		},
		RoomUserList: roomUserLists,
	}
	//给用户消息
	global.GVA_LOG.Infof("LoadCompletedFirst 加载广播: %v", msgData)
	responseHeadByte, _ := json.Marshal(msgData)
	NatsSendAimUserMsg(trs, helper.GetNetMessage("", "", int32(pbs.Meb_loadCompleted), config.SlotServer, responseHeadByte), userId)

	//查找自己的牌并赋值
	cardConfigByIds := logic.GetUserOwnCards(userId)
	trs.ComRoomSpace.AddUserOwnCards(userId, cardConfigByIds)

	//第一次全部加载完成
	if len(trs.LoadComps) == trs.RoomInfo.UserNumLimit {

		trs.IsAllLoadComps = true
		global.GVA_LOG.Infof("LoadCompleted 房间{%v},全部加载完成，开始发牌", trs.RoomInfo.RoomNo)

		//全部加载改变状态 应该是 游戏中 -> 加载中；过度
		trs.ComRoomSpace.GameStateTransition(EnGameStartIng, EnLoadExec)
	}
}

// DoLikeCard : userId：点赞用户 likeUserId：被点赞用户 likeCard：被点赞的卡
func (trs *RoomSpace) DoLikeCard(userId, likeUserId string, likeCard models.LikeCard, likeCards []*models.Card) {
	//纪录被点赞的卡
	//内存纪录
	likeCard.UserID = userId

	trs.ComRoomSpace.SetLikeCardsCard(&likeCard)
	//数据库纪录 todo

	trs.ComRoomSpace.SetLikeUserInfo(userId, &likeCard)

	netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_likeCards), config.SlotServer)
	protoNum := strconv.Itoa(int(pbs.Meb_likeCards))
	msgData := models.LikeCardsMsg{
		ProtoNum:   protoNum,
		Timestamp:  time.Now().Unix(),
		LikeUserId: likeUserId,
		UserId:     userId,
		Card:       likeCards,
	}

	//给客户消息
	global.GVA_LOG.Infof("OperateCards 表情: %v", msgData)
	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	NatsSendAllUserMsg(trs, netMessageResp)

	//如果都点赞了
	if trs.ComRoomSpace.IsAllUserLikeCard() {
		//全部点赞
		if trs.RoomInfo.RoomTurnNum == trs.ComRoomSpace.GetTurn() {
			//结束游戏 结束游戏需要计算统计数据
			if !trs.ComRoomSpace.GameStateTransition(EnLikeCardIng, EnCalculateExec) {
				global.GVA_LOG.Infof("从点赞结束的状态 向全部点赞完成的状态转变 失败 RoomNo:{%v}", trs.RoomInfo.RoomNo)
			} else {
				global.GVA_LOG.Infof("从点赞结束的状态 向全部点赞完成的状态转变 成功 RoomNo:{%v}", trs.RoomInfo.RoomNo)
			}
		} else {
			//进入下一轮
			if !trs.ComRoomSpace.GameStateTransition(EnLikeCardIng, EnNextTurnExec) {
				global.GVA_LOG.Infof("从点赞结束的状态 向全部点赞完成的状态转变 失败 RoomNo:{%v}", trs.RoomInfo.RoomNo)
			} else {
				global.GVA_LOG.Infof("从点赞结束的状态 向全部点赞完成的状态转变 成功 RoomNo:{%v}", trs.RoomInfo.RoomNo)
			}
		}
	}
}

// OutCart 出牌
func (trs *RoomSpace) OutCart(reqCards, cards []*models.Card, userId string) ([]*models.Card, uint32) {
	//请求的牌数量
	var reqCardVerify map[int]int
	reqCardVerify = make(map[int]int)

	for k, _ := range reqCards {
		reqCard := reqCards[k]
		_, ok := reqCardVerify[reqCard.CardId]
		if ok {
			reqCardVerify[reqCard.CardId]++
		} else {
			reqCardVerify[reqCard.CardId] = 1
		}
	}

	//本来手里的牌数量
	var cardVerify map[int]int
	cardVerify = make(map[int]int)

	for k, _ := range cards {
		card := cards[k]
		_, ok := cardVerify[card.CardId]
		if ok {
			cardVerify[card.CardId]++
		} else {
			cardVerify[card.CardId] = 1
		}
	}

	//比较是否一致
	for cId, cartNum := range reqCardVerify {
		num := cardVerify[cId]
		if num < cartNum {
			return cards, common.NotCanOutCards
		}
	}

	//剩下的牌
	var newCards []*models.Card
	//要出的牌
	var outCards []*models.Card

	for idx, _ := range cards {
		card := cards[idx]
		isInHand := false

		idNum, ok := reqCardVerify[card.CardId]
		if ok && idNum > 0 {
			//有相同的牌
			isInHand = true
			reqCardVerify[card.CardId]--
		}

		if isInHand {
			outCards = append(outCards, card)
		} else {
			newCards = append(newCards, card)
		}
	}

	protoNum := strconv.Itoa(int(pbs.Meb_outCards))
	netMessageResp := helper.NewNetMessage("", "", int32(pbs.Meb_outCards), config.SlotServer)
	//非出牌用户得到的消息
	msgData := models.OperateCardsMsg{
		ProtoNum:   protoNum,
		Timestamp:  time.Now().Unix(),
		UserId:     userId,
		OutCardNum: len(outCards),
		CardNum:    len(newCards),
	}
	responseHeadByte, _ := json.Marshal(msgData)
	netMessageResp.Content = responseHeadByte
	global.GVA_LOG.Infof("出牌的广播: %v", string(responseHeadByte))
	//非出牌用户得到的消息
	trs.ComRoomSpace.SendExcludeUserMsg(netMessageResp, userId)

	//===

	netMessageResp1 := helper.NewNetMessage("", "", int32(pbs.Meb_outCards), config.SlotServer)
	//出牌用户得到的消息
	aimMsgData := models.OperateCardsMsg{
		ProtoNum:   protoNum,
		Timestamp:  time.Now().Unix(),
		UserId:     userId,
		OutCardNum: len(outCards),
		CardNum:    len(newCards),
		Card:       newCards,
	}
	global.GVA_LOG.Infof("出牌的广播: %v", aimMsgData)
	aimMsgDataByte, _ := json.Marshal(aimMsgData)
	netMessageResp1.Content = aimMsgDataByte

	//出牌用户得到的消息
	NatsSendAimUserMsg(trs, netMessageResp1, userId)

	//ReMakeCurrCard 重置当前手中的卡
	trs.ComRoomSpace.ReMakeCurrCard(userId, newCards, outCards)

	//都出过牌的时候
	if trs.ComRoomSpace.IsAllUserOutCart() {
		global.GVA_LOG.Infof("都出牌 进入点赞阶段 %v", trs.ComRoomSpace.GetGameState())
		//trs.ComRoomSpace.ChangeGameState(EnLikePage)
		// 这个应该是从出牌状态 到点赞页面状态 但是没有做随牌到出牌的状态改变 所以这里 状态直接设置为 去点赞
		// 应该是 出牌状态（随牌状态） -> 点赞
		trs.ComRoomSpace.GameStateTransition(RemakeCardIng, EnLikePageExec)
	}

	return newCards, common.OK
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
