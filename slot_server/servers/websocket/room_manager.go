package websocket

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"sync"
	"time"
)

// roomManager 相当于管理者
type roomManager struct {
	//key房间编号 所有的房间
	Rooms map[string]*RoomSpace

	//item  全部的卡牌 信息
	Cards []models.Card

	//匹配的逻辑
	MatchIngRoom *MatchIngRoom

	//通用字段
	CommonRoomManager *CommonRoomManager

	//房间配置 key 版本 0 基础版本
	RoomConfigMap map[int]*table.MemeRoomConfig

	//基础卡
	RoomBaseCard []*table.MbCardConfig

	//版本卡
	RoomVersionCard map[int][]*table.MbCardConfig

	//基础问题
	RoomIssueConfig []*table.MbIssueConfig

	RobotAiUser []*table.RobotAiUser
}

var SlotRoomManager = roomManager{
	Rooms:             make(map[string]*RoomSpace), //key房间编号
	CommonRoomManager: GetCommonRoomManager(),
	MatchIngRoom:      NewMatchIngRoom(),
	RoomConfigMap:     make(map[int]*table.MemeRoomConfig),
	RoomBaseCard:      make([]*table.MbCardConfig, 0),
	RoomVersionCard:   make(map[int][]*table.MbCardConfig),
}

func (trMgr *roomManager) Start() {
	trMgr.InitDBData()

	// 创建一个计时器
	serviceTimer := time.NewTicker(time.Second * 30)
	defer serviceTimer.Stop() //定时器不用了需要关闭

	serviceRoomTimer := time.NewTicker(time.Second * 300) //todo
	defer serviceRoomTimer.Stop()                         //定时器不用了需要关闭

	stopRoomTimer := time.NewTimer(10 * time.Second)
	defer stopRoomTimer.Stop() //定时器不用了需要关闭

	//匹配房间的定时器
	matchRoomTimer := time.NewTicker(time.Second * 7)
	defer matchRoomTimer.Stop()

	//这个定时器是分析房间游戏是否开始
	defer func() {
		if err := recover(); err != nil {
			global.GVA_LOG.Error("roomManager :", zap.Any("recover ", err))
		}
		global.GVA_LOG.Info("roomManager end")
	}()

	for {
		select {
		case <-serviceTimer.C:
			//打印一些调试信息
			global.GVA_LOG.Infof("房间管理器:计时器触发,房间长度:%d", len(trMgr.Rooms))

		case <-serviceRoomTimer.C:
			global.GVA_LOG.Infof("房间管理器处理一些问题 %d", len(trMgr.Rooms))
			//1 如果房间长时间没有操作
			//trMgr.DelRoomByTime() //房间协程 自己会回收自己

		case <-stopRoomTimer.C:

		case <-matchRoomTimer.C:
			//匹配房间的算法

		case message := <-trMgr.CommonRoomManager.CloseRoom:
			// 关闭房间
			closeInfo := &models.CloseRoom{}
			err := json.Unmarshal(message, closeInfo)
			if err != nil {
				global.GVA_LOG.Error("closeInfo", zap.Error(err))
			}
			if closeInfo.IsStop {
				trMgr.SendCloseRoomMsg(closeInfo)
			}

		case message := <-trMgr.CommonRoomManager.Broadcast:
			// 广播事件
			global.GVA_LOG.Info("房间管理器:广播事件")
			trMgr.SendMsg(message)
		}
	}
}

type MatchIngRoom struct {
	Sync *sync.Mutex

	IsArithmeticIng bool

	//匹配中的房间 key房间好 val:房间用户
	//单排
	MatchIngRoom1User []*MatchIngRoomInfo

	//双排
	MatchIngRoom2User []*MatchIngRoomInfo
}

type MatchIngRoomInfo struct {
	RoomNo      string
	StartTime   *time.Time //匹配的开始时间
	UserInfoArr []*models.UserInfo
}

type MatchGroupRoomInfo struct {
	RoomNo      string
	DelRoomNo   []string
	UserInfoArr []*models.UserInfo
	//被删用户组
}

func NewMatchIngRoom() *MatchIngRoom {
	matchIngRoom := MatchIngRoom{
		Sync:              new(sync.Mutex),
		IsArithmeticIng:   false,
		MatchIngRoom1User: make([]*MatchIngRoomInfo, 0),
		MatchIngRoom2User: make([]*MatchIngRoomInfo, 0),
	}
	return &matchIngRoom
}

func NewMatchIngRoomInfo(roomNo string, uInfo []*models.UserInfo) *MatchIngRoomInfo {
	return &MatchIngRoomInfo{
		RoomNo:      roomNo,
		StartTime:   helper.LocalTime(),
		UserInfoArr: uInfo,
	}
}

func (trMgr *roomManager) JoinMatchIngRoom(roomNo string) {
	trMgr.CommonRoomManager.MatchLock.Lock()
	defer trMgr.CommonRoomManager.MatchLock.Unlock()
	global.GVA_LOG.Infof("JoinMatchIngRoom:roomNo:%v ,time:%v", roomNo, helper.TimeIntToStr(time.Now().Unix()))

	space := trMgr.Rooms[roomNo]

	trMgr.MatchIngRoom.AddMatchIngRoomInfo(space)
}

func (mu *MatchIngRoom) IsCanAddMatchIng(roomNo string, userNum int) bool {
	isInMatchIng := false
	if userNum == 1 {
		for k, _ := range mu.MatchIngRoom1User {
			item := mu.MatchIngRoom1User[k]
			if item.RoomNo == roomNo {
				isInMatchIng = true
				break
			}
		}
	}

	if userNum == 2 {
		for k, _ := range mu.MatchIngRoom2User {
			item := mu.MatchIngRoom2User[k]
			if item.RoomNo == roomNo {
				isInMatchIng = true
				break
			}
		}
	}

	return isInMatchIng
}

func (mu *MatchIngRoom) AddMatchIngRoomInfo(roomSpace *RoomSpace) {
	mu.Sync.Lock()
	defer mu.Sync.Unlock()

	//获取房间 信息
	//房间号 和 房间人
	roomNo := roomSpace.RoomInfo.RoomNo
	userInfoMaps := roomSpace.ComRoomSpace.UserInfos

	uInfo := []*models.UserInfo{}
	for _, userInfoMap := range userInfoMaps {
		uInfo = append(uInfo, userInfoMap)
	}

	matchIngRoomInfo := NewMatchIngRoomInfo(roomNo, uInfo)

	if len(userInfoMaps) == 1 {
		//单排
		if !mu.IsCanAddMatchIng(roomNo, 1) {
			mu.MatchIngRoom1User = append(mu.MatchIngRoom1User, matchIngRoomInfo)
		}
	} else if len(userInfoMaps) == 2 {
		//双排
		if !mu.IsCanAddMatchIng(roomNo, 2) {
			mu.MatchIngRoom2User = append(mu.MatchIngRoom2User, matchIngRoomInfo)
		}
	} else {
		global.GVA_LOG.Infof("AddMatchIngRoomInfo 匹配房间人数不对")
	}

	global.GVA_LOG.Infof("AddMatchIngRoomInfo 单排房间多少个{%v} 双排房间多少个{%v}", len(mu.MatchIngRoom1User), len(mu.MatchIngRoom2User))
}

func (trMgr *roomManager) CancelMatchIngUser(roomNo string) {
	trMgr.CommonRoomManager.MatchLock.Lock()
	defer trMgr.CommonRoomManager.MatchLock.Unlock()

	space := trMgr.Rooms[roomNo]
	trMgr.MatchIngRoom.CancelMatchRoom(space)
}

func (mu *MatchIngRoom) CancelMatchRoom(roomSpace *RoomSpace) {
	mu.Sync.Lock()
	mu.Sync.Unlock()
	//获取房间 信息
	//房间号 和 房间人
	roomNo := roomSpace.RoomInfo.RoomNo
	userInfoMaps := roomSpace.ComRoomSpace.UserInfos

	if len(userInfoMaps) == 1 {
		//单排
		newMatchIngRoom := []*MatchIngRoomInfo{}

		for k, _ := range mu.MatchIngRoom1User {
			matchIngRoom := mu.MatchIngRoom1User[k]
			if matchIngRoom.RoomNo == roomNo {
				continue
			}
			newMatchIngRoom = append(newMatchIngRoom, matchIngRoom)
		}
		mu.MatchIngRoom1User = newMatchIngRoom
	} else if len(userInfoMaps) == 2 {
		//双排
		newMatchIngRoom := []*MatchIngRoomInfo{}

		for k, _ := range mu.MatchIngRoom2User {
			matchIngRoom := mu.MatchIngRoom2User[k]
			if matchIngRoom.RoomNo == roomNo {
				continue
			}
			newMatchIngRoom = append(newMatchIngRoom, matchIngRoom)
		}
		mu.MatchIngRoom2User = newMatchIngRoom
	}

	global.GVA_LOG.Infof("CancelMatchRoom 单排房间多少个{%v} 双排房间多少个{%v}", len(mu.MatchIngRoom1User), len(mu.MatchIngRoom2User))
}

func (trMgr *roomManager) DelRoomByTime() {
	trMgr.CommonRoomManager.Sync.Lock()
	defer trMgr.CommonRoomManager.Sync.Unlock()

	for no, val := range trMgr.Rooms {
		//大于五分钟
		if time.Now().Unix()-val.ComRoomSpace.CurrentOpTime > RoomAlive {
			global.GVA_LOG.Infof("房间管理器清理过期房间:{%v},roomInfo:%v", val.RoomInfo.RoomNo, val.RoomInfo)
			val.RoomInfo.IsOpen = table.RoomStatusAbnormal
			//更新数据库 房间状态
			err := table.SaveMemeRoom(val.RoomInfo)
			if err != nil {
				global.GVA_LOG.Error("DelRoomByTime LeaveRoomNotStartGame", zap.Error(err))
			}
			trMgr.DelRoom(no)
		}
	}
}

func (trMgr *roomManager) AddRoomSpace(roomNo string, roomSpaceInfo *RoomSpace) {
	trMgr.CommonRoomManager.Sync.Lock()
	defer trMgr.CommonRoomManager.Sync.Unlock()

	_, b := trMgr.Rooms[roomNo]
	if !b {
		trMgr.Rooms[roomNo] = roomSpaceInfo
	}
}

func (trMgr *roomManager) GetRoomSpace(roomNo string) (*RoomSpace, error) {
	trMgr.CommonRoomManager.Sync.RLock()
	defer trMgr.CommonRoomManager.Sync.RUnlock()

	roomSpaceInfo, b := trMgr.Rooms[roomNo]
	if !b {
		return nil, errors.New("没有房间信息")
	}
	return roomSpaceInfo, nil
}

func (trMgr *roomManager) GetCurrRoomSpace() (*RoomSpace, error) {
	trMgr.CommonRoomManager.Sync.RLock()
	defer trMgr.CommonRoomManager.Sync.RUnlock()

	if len(trMgr.Rooms) <= 0 {
		return nil, errors.New("没有房间信息")
	}

	//todo 后续通过城市
	roomSpaceInfo := &RoomSpace{}
	for _, v := range trMgr.Rooms {
		roomSpaceInfo = v
		break
	}

	return roomSpaceInfo, nil
}

// SendMsgToRoomSpace 房间管理器 消息发送到房间
func (trMgr *roomManager) SendMsgToRoomSpace(roomNo string, message []byte) error {
	trMgr.CommonRoomManager.Sync.RLock()
	defer trMgr.CommonRoomManager.Sync.RUnlock()

	roomSpaceInfo, b := trMgr.Rooms[roomNo]
	if !b {
		global.GVA_LOG.Infof("SendMsgToRoomSpace 没有房间信息 %v", roomNo)
		return errors.New("没有房间信息")
	}

	if roomSpaceInfo.ComRoomSpace.ReceiveMsg != nil {
		roomSpaceInfo.ComRoomSpace.ReceiveMsg <- message
	}

	return nil
}

//func (trMgr *roomManager) DelRoomSpace(roomId string) {
//	trMgr.CommonRoomManager.Sync.Lock()
//	defer trMgr.CommonRoomManager.Sync.Unlock()
//	delete(trMgr.Rooms, roomId)
//}

func (trMgr *roomManager) DelRoom(roomNo string) {
	trMgr.CommonRoomManager.Sync.Lock()
	defer trMgr.CommonRoomManager.Sync.Unlock()
	global.GVA_LOG.Infof(" DelRoom,currTime:%v RoomNo:%v", helper.TimeIntToStr(time.Now().Unix()), roomNo)
	delete(trMgr.Rooms, roomNo)
}

func (trMgr *roomManager) SendCloseRoomMsg(closeInfo *models.CloseRoom) {
	trMgr.CommonRoomManager.Sync.Lock()
	defer trMgr.CommonRoomManager.Sync.Unlock()
	roomSpaceInfo, ok := trMgr.Rooms[closeInfo.RoomNo]
	if !ok {
		return
	}

	//从全局管理器中 删除结束房间
	roomSpaceInfo.ComRoomSpace.Close <- true
}

// MatchRoomArithmetic 匹配房间的算法
func (trMgr *roomManager) MatchRoomArithmetic() {
	trMgr.CommonRoomManager.MatchLock.Lock()
	defer trMgr.CommonRoomManager.MatchLock.Unlock()

	//寻找匹配中的用户

	//单排
	trMgr.DealMatchIng1User()

	//双排
	//trMgr.DealMatchIng2User()

	//双排 + 单排
	//trMgr.DealMatchIngUser()

	//补充机器人逻辑
	trMgr.DealMatchIngUserAddRobot()
}

func (trMgr *roomManager) InitDBData() {

	//初始化全局动物派对房间

	//1 先创建对局空间
	roomSpaceInfo := GetRoomSpace()
	//添加对局用户

	record, err := table.GetMemeRoomByIdDesc()
	if err != nil {
		global.GVA_LOG.Error("InitDBData GetMemeRoomByIdDesc err:", zap.Error(err))
	}
	period := "1"
	if record.ID >= 0 {
		period = helper.Itoa(helper.Atoi(record.Period) + 1)
	}

	animalPartyRoom := table.NewAnimalPartyRoom("1", "1", uuid.New().String(), config.AnimalPartyGlobal, "匹配房间", period,
		table.TavernRoomOpen, table.RoomTypeMatch, 0, 0, 0, 0)
	err = table.CreateMemeRoom(animalPartyRoom)
	if err != nil {
		global.GVA_LOG.Error("NewAnimalPartyRoom:{%v},roomInfo:%v", zap.Error(err), zap.Any("NewAnimalPartyRoom", animalPartyRoom.RoomNo))
		return
	}

	//给用户创建房间 并发送游戏开始的广播
	//slotRoom, err := table.SlotRoomByRoomNo(config.AnimalPartyGlobal)
	//if err != nil {
	//	global.GVA_LOG.Error("MebJoinRoom ", zap.Error(err))
	//}
	roomSpaceInfo.RoomInfo = animalPartyRoom

	//游戏开始
	roomSpaceInfo.ComRoomSpace.IsStartGame = true

	//游戏每 小轮状态 游戏开始
	roomSpaceInfo.ComRoomSpace.ChangeGameState(BetIng)

	roomSpaceInfo.ComRoomSpace.SetGameStartTime(helper.LocalTime().Unix()) //游戏开始时间

	//添加到全局房间管理器
	SlotRoomManager.AddRoomSpace(animalPartyRoom.Name, roomSpaceInfo)

	//颜色配置
	roomSpaceInfo.ColorConfigs = GetColorWheel()
	go roomSpaceInfo.Start()
}

func (trMgr *roomManager) SendMsg(message []byte) {
	//解析消息 找到对应的房间广播消息
	//todo
	//这里需要解析 message 信息，得到房间编号，再给房间里的用户发消息

	//roomSpaceInfo := trMgr.TavernRooms["1"]
	//for _, client := range roomSpaceInfo.UserClient {
	//	select {
	//	case client.Send <- message:
	//	default:
	//		//close(client.Send)
	//	}
	//}

	//或者给全部的房间发消息
	//for range trMgr.TavernRooms {
	//for _, client := range roomSpaceInfo.UserClient {
	//	select {
	//	case client.Send <- message:
	//	default:
	//		//close(client.Send)
	//	}
	//}
	//}

}

// IsCanStartPlay 是否可以开始游戏
func (trMgr *roomManager) IsCanStartPlay(roomNo string) bool {
	//房间所有的人是否就绪
	roomSpaceInfo, err := trMgr.GetRoomSpace(roomNo)
	if err != nil {
		global.GVA_LOG.Error("IsCanStartPlay GetRoomSpace", zap.Error(err))
		return false
	}

	b, err := roomSpaceInfo.ComRoomSpace.IsCanStartPlay(roomSpaceInfo.RoomInfo.UserNumLimit)
	if err != nil && b == false {
		global.GVA_LOG.Error("IsCanStartPlay IsCanStartPlay", zap.Error(err))
		return false
	}

	return true
}
