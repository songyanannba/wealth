package websocket

import (
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"slot_server/lib/models/table"
	"slot_server/protoc/pbs"
	"sync"
	"time"
)

type PeriodSpace struct {

	//key是用户ID 用户信息 {用户押注/获取房间的时候添加 }
	UserInfos map[string]*models.UserInfo

	//接收消息处理
	ReceiveMsg chan []byte

	//广播类型的消息 消息中需要有房间号
	Broadcast chan []byte

	//关闭房间的chan
	Close chan bool

	//上次操作房间的时间 ｜ 存活时间
	CurrentOpTime int64

	//是否开始游戏
	IsStartGame bool

	//房间游戏是否结束 结束需要回收资源
	IsStopGame bool
	//房间结束的时间戳

	StopGameTime int64

	//每小局详情
	TurnMateInfo TurnMateInfo
	Sync         *sync.Mutex
	UserSync     *sync.RWMutex

	//是否匹配成功后，被清理的房间
	IsMatchClear bool

	//是否正在发牌
	IdDealCarding bool
}

// AllAnimalWheelSort 最外圈动物怕排序
type AllAnimalWheelSort struct {
	WinSeat             int
	WinAnimalConfig     *AnimalConfig       //赢钱的动物位置
	AnimalConfigs       []*AnimalConfig     //当前动物排序
	WinBetZoneConfig    []*BetZoneConfig    //赢钱区域
	WinBigOrSmallConfig *BigOrSmallConfig   //赢钱的大小位置
	BigOrSmallConfigs   []*BigOrSmallConfig //当前大小排序
}

type ComRoomSpace struct {
	//是否被保护
	IsProtection bool
	//房主
	UserOwner *models.UserInfo
	//当前房间的用户
	//key是用户ID {是在获取房间/押注的时候填充的数据，所有的真实用户}
	//UserClient map[string]*Client

	//机器人押注次数 //todo
	APRobotActionCount int

	//key是用户ID 用户信息 {用户押注/获取房间的时候添加 }
	UserInfos map[string]*models.UserInfo

	//全部的情况  如果是 LUCKY 有多个
	CurrAnimalWheelSort []*AllAnimalWheelSort

	//接收消息处理
	ReceiveMsg chan []byte

	//广播类型的消息 消息中需要有房间号
	Broadcast chan []byte

	//关闭房间的chan
	Close chan bool

	//上次操作房间的时间 ｜ 存活时间
	CurrentOpTime int64

	//是否开始游戏
	IsStartGame bool

	//房间游戏是否结束 结束需要回收资源
	IsStopGame bool
	//房间结束的时间戳

	StopGameTime int64

	//每小局详情
	TurnMateInfo TurnMateInfo
	Sync         *sync.Mutex
	UserSync     *sync.RWMutex

	//是否匹配成功后，被清理的房间
	IsMatchClear bool

	//是否正在发牌
	IdDealCarding bool
}

type TurnMateInfo struct {
	TurnSync *sync.RWMutex

	//第几轮
	Turn int

	//区域 ｜ 用户
	BetZoneUserInfoMap map[int]map[string]*models.UserInfo

	UserIsWin bool

	//每轮游戏环节状态
	GameTurnStatus GameTurnState

	//开始游戏广播后的时间
	GameStartTime int64

	//全部加载 或者 进入下一轮的时间
	//也就是每小轮的开始时间
	CountdownTime int64

	//每小轮的点赞倒计时
	LikeCountdownTime int64

	//在每一轮的最后确定时间
	// key是轮数 用户在每一轮的信息 【选择房间的时候/强制开始的时候】
	TurnUserInfoMap map[int][]*models.UserInfo

	//只保留一下上个用户最近一次出过的牌 （冗余OutCards里面的数据）
	LastUserOutCards map[string]*TurnUserCard

	//是否重随过牌
	IsReGiveCard map[int]bool
}

type TurnUserCard struct {
	Ord   int //次序
	Cards []*models.Card
}

type LikeCardInfo struct {
	//点赞用户
	likeUser *models.UserInfo
	//点在的卡
	Cards []*models.LikeCard
}

// DoubtInfo 被质疑的集合
type DoubtInfo struct {
	//质疑者
	DoubtUser *models.UserInfo
	//被质疑者
	BeDoubtUser *models.UserInfo
	//质疑的牌
	Cards []*models.Card
	//胜利者
	WinUser *models.UserInfo
}

func (rs *ComRoomSpace) SetCanSendCard(isCan bool) {
	//rs.TurnMateInfo.IsCanSendCard[rs.GetCurrTurn()] = isCan
}

func (rs *ComRoomSpace) GetCurrTurn() int {
	return rs.TurnMateInfo.Turn
}

func (rs *ComRoomSpace) SetReceiveMsg(msgId string, data []byte) {
	msg := models.ComMsg{
		MsgId: msgId,
		Data:  data,
	}
	msgByte, _ := json.Marshal(msg)
	rs.ReceiveMsg <- msgByte
}

func (rs *ComRoomSpace) UpdateTurnMateInfo(turn int, countdownTime int64, matchUser []*models.UserInfo) {
	//rs.TurnMateInfo.CountdownTime = countdownTime
	rs.TurnMateInfo.TurnUserInfoMap[turn] = matchUser
}

// GetComRoomSpace ss
func GetComRoomSpace() *ComRoomSpace {
	return &ComRoomSpace{
		UserOwner:           &models.UserInfo{},
		IsProtection:        false,
		UserInfos:           make(map[string]*models.UserInfo),
		Broadcast:           make(chan []byte),
		ReceiveMsg:          make(chan []byte, 10000),
		Close:               make(chan bool),
		CurrentOpTime:       time.Now().Unix(), //内存中房间协程创建时间
		APRobotActionCount:  0,
		IsStopGame:          false,
		IsStartGame:         false,
		CurrAnimalWheelSort: make([]*AllAnimalWheelSort, 0),
		Sync:                new(sync.Mutex),
		UserSync:            new(sync.RWMutex),
		TurnMateInfo: TurnMateInfo{
			TurnSync:           new(sync.RWMutex),
			BetZoneUserInfoMap: make(map[int]map[string]*models.UserInfo),
			Turn:               0,
			UserIsWin:          false,
			TurnUserInfoMap:    make(map[int][]*models.UserInfo),
			LastUserOutCards:   make(map[string]*TurnUserCard),
			IsReGiveCard:       make(map[int]bool),
			GameTurnStatus:     0,
			CountdownTime:      0,
			GameStartTime:      0,
			LikeCountdownTime:  0,
		},
	}
}

// AddBetZoneUserInfoMap 给对应区域增加用户 或者用户押注
func (rs *ComRoomSpace) AddBetZoneUserInfoMap(batZoneId int, bat float32, userInfo *models.UserInfo) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	_, ok := rs.TurnMateInfo.BetZoneUserInfoMap[batZoneId]
	if !ok {
		rs.TurnMateInfo.BetZoneUserInfoMap[batZoneId] = make(map[string]*models.UserInfo)
	}

	uInfo, ok := rs.TurnMateInfo.BetZoneUserInfoMap[batZoneId][userInfo.UserID]
	if !ok {
		userInfo.UserProperty.Bet = float64(bat)
		rs.TurnMateInfo.BetZoneUserInfoMap[batZoneId][userInfo.UserID] = userInfo
	} else {
		uInfo.UserProperty.Bet = uInfo.UserProperty.Bet + float64(bat)
		rs.TurnMateInfo.BetZoneUserInfoMap[batZoneId][userInfo.UserID] = uInfo
	}

}

// GetBetZoneUserInfos 根据押注区域ID 获取赢钱和输钱用户
func (rs *ComRoomSpace) GetBetZoneUserInfos(batZoneId int) (winUserArr []*models.UserInfo, loseUserArr []*models.UserInfo) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	for batZoneKey, mapUInfo := range rs.TurnMateInfo.BetZoneUserInfoMap {
		if batZoneKey == batZoneId {
			for _, uInfo := range mapUInfo {
				winUserArr = append(winUserArr, uInfo)
			}
		} else {
			for _, uInfo := range mapUInfo {
				loseUserArr = append(winUserArr, uInfo)
			}
		}
	}
	return
}

// AddUserInfos 添加用户信息
func (rs *ComRoomSpace) AddUserInfos(userId string, userInfo *models.UserInfo) {
	rs.UserSync.Lock()
	defer rs.UserSync.Unlock()

	_, ok := rs.UserInfos[userId]
	if !ok {
		rs.UserInfos[userId] = userInfo
	}
}

func (rs *ComRoomSpace) UpdateUserInfoAndUserClient(userId string, userInfo *models.UserInfo) {
	rs.UserSync.Lock()
	defer rs.UserSync.Unlock()
	rs.UserInfos[userId] = userInfo
}

// DelUserInfoAndUserClient 删除用户信息
func (rs *ComRoomSpace) DelUserInfoAndUserClient(userId string) {
	rs.UserSync.Lock()
	defer rs.UserSync.Unlock()
	_, ok := rs.UserInfos[userId]
	if ok {
		delete(rs.UserInfos, userId)
	}
}

func (rs *ComRoomSpace) DelRoomAllUserAndUserClient(roomNo string) {
	rs.UserSync.Lock()
	defer rs.UserSync.Unlock()
	for k, _ := range rs.UserInfos {
		userId := rs.UserInfos[k].UserID
		_, ok := rs.UserInfos[userId]
		if ok {
			delete(rs.UserInfos, userId)
		}
		err := table.DelRoomUsers(roomNo, userId)
		if err != nil {
			global.GVA_LOG.Error("DelRoomAllUserAndUserClient", zap.Any("err", err))
		}
	}
}

func (rs *ComRoomSpace) GetUserInfo(userId string) (*models.UserInfo, error) {
	rs.UserSync.RLock()
	defer rs.UserSync.RUnlock()

	uInfo, ok := rs.UserInfos[userId]
	if !ok {
		return nil, errors.New("user not found")
	}
	return uInfo, nil
}

func (rs *ComRoomSpace) SetUserSeat(userId string, seat int) error {
	rs.UserSync.RLock()
	defer rs.UserSync.RUnlock()

	uInfo, ok := rs.UserInfos[userId]
	if !ok {
		return errors.New("user not found")
	}
	uInfo.UserProperty.Seat = seat
	return nil
}

//func (rs *ComRoomSpace) GetUserClient(userId string) (*Client, error) {
//	rs.UserSync.RLock()
//	defer rs.UserSync.RUnlock()
//
//	cliInfo, ok := rs.UserClient[userId]
//	if !ok {
//		return nil, errors.New("user client not found")
//	}
//	return cliInfo, nil
//}

func (rs *ComRoomSpace) GetAllUserInfo() ([]*models.UserInfo, error) {
	rs.UserSync.RLock()
	defer rs.UserSync.RUnlock()

	allUser := make([]*models.UserInfo, 0)
	for _, val := range rs.UserInfos {
		allUser = append(allUser, val)
	}

	return allUser, nil
}

func (rs *ComRoomSpace) UserInfoToRoomUser() ([]models.MemeRoomUser, error) {
	rs.UserSync.RLock()
	defer rs.UserSync.RUnlock()

	roomUserLists := []models.MemeRoomUser{}

	for _, userItem := range rs.UserInfos {
		tavernRoomUser := models.MemeRoomUser{
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
			UserCards: models.UserCartState{
				OutCardNum: 0,
				CardNum:    0,
			},
		}
		roomUserLists = append(roomUserLists, tavernRoomUser)
	}

	return roomUserLists, nil
}

//func (rs *ComRoomSpace) AddSpaceUserClient(uid string, client *Client) {
//	rs.UserSync.Lock()
//	defer rs.UserSync.Unlock()
//	_, ok := rs.UserClient[uid]
//	if !ok {
//		rs.UserClient[uid] = client
//	}
//}

// IsCanStartPlay 是否可以开始游戏
func (rs *ComRoomSpace) IsCanStartPlay(userNumLimit int) (bool, error) {
	//房间所有的人是否就绪
	if len(rs.UserInfos) != userNumLimit {
		global.GVA_LOG.Error("ComRoomSpace IsCanStartPlay 房间人数不够")
		return false, errors.New("房间人数不够")
	}

	for _, userInfos := range rs.UserInfos {
		if userInfos.GetUserIsReady() != int(models.Ready) {
			return false, errors.New("没有全部就绪")
		}
	}
	return true, nil
}

func (rs *ComRoomSpace) SendAllUserMsg(msgBty []byte) {
	//给客户消息
	for _, userInfo := range rs.UserInfos {
		global.GVA_LOG.Infof("SendAllUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, string(msgBty))
		//nats.NastManager.Producer(netMessageRespMarshal)
	}

}

func (rs *ComRoomSpace) NatsSendAllUserMsg(msg *pbs.NetMessage) {
	//给客户消息
	for _, userInfo := range rs.UserInfos {
		//离开或者机器人不发送消息
		msg.AckHead.Uid = userInfo.UserID
		netMessageRespMarshal, _ := proto.Marshal(msg)
		global.GVA_LOG.Infof("NatsSendAllUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, msg)
		NastManager.Producer(netMessageRespMarshal)
	}
}

func (rs *ComRoomSpace) NatsSendAimUserMsg(msg *pbs.NetMessage, userId string) {
	//给客户消息
	msg.AckHead.Uid = userId
	netMessageRespMarshal, _ := proto.Marshal(msg)
	global.GVA_LOG.Infof("NatsSendAimUserMsg UserID:{%v} 给客户端发消息:{%v}", userId, msg)
	NastManager.Producer(netMessageRespMarshal)
}

// MatsSendExcludeUserMsg 排除制定用户
func (rs *ComRoomSpace) MatsSendExcludeUserMsg(msg *pbs.NetMessage, userId string) {
	//给客户消息
	for _, userInfo := range rs.UserInfos {
		if userInfo.UserID == userId {
			continue
		}
		if len(userInfo.UserID) == 0 || userInfo.UserProperty.IsLeave == 1 {
			continue
		}
		//userIDInt, err := strconv.Atoi(userInfo.UserID)
		//if err != nil {
		//	global.GVA_LOG.Error("MatsSendExcludeUserMsg err", zap.Any("err", err))
		//	continue
		//}

		msg.AckHead.Uid = userInfo.UserID
		netMessageRespMarshal, _ := proto.Marshal(msg)
		global.GVA_LOG.Infof("MatsSendExcludeUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, msg)
		NastManager.Producer(netMessageRespMarshal)
	}
}

// SendExcludeUserMsg 排除制定用户
func (rs *ComRoomSpace) SendExcludeUserMsg(msg *pbs.NetMessage, userId string) {
	//给客户消息
	for _, userInfo := range rs.UserInfos {
		if userInfo.UserID == userId {
			continue
		}
		if len(userInfo.UserID) == 0 || userInfo.UserProperty.IsLeave == 1 || userInfo.UserProperty.IsRobot == 1 {
			continue
		}
		//userIDInt, err := strconv.Atoi(userInfo.UserID)
		//if err != nil {
		//	global.GVA_LOG.Error("NatsSendAllUserMsg err", zap.Any("err", err))
		//	continue
		//}
		msg.AckHead.Uid = userInfo.UserID
		netMessageRespMarshal, _ := proto.Marshal(msg)
		global.GVA_LOG.Infof("NatsSendAllUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, msg)
		NastManager.Producer(netMessageRespMarshal)
	}
}

func (rs *ComRoomSpace) GetTurn() int {
	return rs.TurnMateInfo.Turn
}

func (rs *ComRoomSpace) AddTurn() {
	//房间属性
	rs.TurnMateInfo.Turn++

	//用户（属性）身上也同步当前轮
	for k, _ := range rs.UserInfos {
		rs.UserInfos[k].AddUserTurn()
	}
}

func (rs *ComRoomSpace) SetStopGame(isEndGame bool) {
	rs.IsStopGame = isEndGame
	rs.StopGameTime = time.Now().Unix()
}

func (rs *ComRoomSpace) SetCountdownTime(countdownTime int64) {
	rs.TurnMateInfo.CountdownTime = countdownTime
}

func (rs *ComRoomSpace) GetCountdownTime() int64 {
	return rs.TurnMateInfo.CountdownTime
}

func (rs *ComRoomSpace) SetLikeCountdownTime(countdownTime int64) {
	rs.TurnMateInfo.LikeCountdownTime = countdownTime
}

func (rs *ComRoomSpace) GetLikeCountdownTime() int64 {
	return rs.TurnMateInfo.LikeCountdownTime
}

func (rs *ComRoomSpace) SetGameStartTime(gameStartTime int64) {
	rs.TurnMateInfo.GameStartTime = gameStartTime
}

func (rs *ComRoomSpace) GetGameStartTime() int64 {
	return rs.TurnMateInfo.GameStartTime
}
