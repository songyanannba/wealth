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
	WinSeat          int
	WinAnimalConfig  *AnimalConfig    //赢钱的动物位置
	AnimalConfigs    []*AnimalConfig  //当前排序
	WinBetZoneConfig []*BetZoneConfig //赢钱区域
}

type ComRoomSpace struct {
	//是否被保护
	IsProtection bool
	//房主
	UserOwner *models.UserInfo
	//当前房间的用户
	//key是用户ID {是在获取房间/押注的时候填充的数据，所有的真实用户}
	//UserClient map[string]*Client

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

	//区域 ｜ 用户
	BetZoneUserInfoMap map[int]map[string]*models.UserInfo

	UserIsWin bool

	//每轮游戏环节状态
	GameTurnStatus GameTurnState

	//开始游戏广播后的时间
	GameStartTime int64

	//第几轮
	Turn int

	//全部加载 或者 进入下一轮的时间
	//也就是每小轮的开始时间
	CountdownTime int64

	//每小轮的点赞倒计时
	LikeCountdownTime int64

	//在每一轮的最后确定时间
	// key是轮数 用户在每一轮的信息 【选择房间的时候/强制开始的时候】
	TurnUserInfoMap map[int][]*models.UserInfo

	//每轮的骗子牌 key(int)是轮数 card 是每轮骗子牌
	FraudCard map[int]*models.Card

	//每轮被选的问题 key(int)是轮数 Issue 是每轮问题
	SelectIssue map[int]*models.Issue

	//key(int)是轮数 string 是用户ID  card 是每轮当前拥有的牌
	Cards map[int]map[string][]*models.Card

	//出牌集合  key(int)是轮数 string 是用户ID  TurnUserCard 是每轮出牌的集合
	OutCards map[int]map[string][]*TurnUserCard

	//key(int)是轮数 var是每轮被点赞的牌
	likeCards map[int][]*models.LikeCard

	likeUserInfo map[int]map[string][]*models.LikeCard

	//只保留一下上个用户最近一次出过的牌 （冗余OutCards里面的数据）
	LastUserOutCards map[string]*TurnUserCard

	//是否重随过牌
	IsReGiveCard map[int]bool

	//本轮已经抽过的牌
	ExtractCard map[int]map[string][]*table.MbCardConfig

	//本轮未被抽过的牌
	NotExtractCard map[int]map[string][]*table.MbCardConfig

	//用户自己解锁的卡
	UserOwnCards map[string][]*table.MbCardConfig

	//是否可以发牌
	//IsCanSendCard map[int]bool
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
			Cards:              make(map[int]map[string][]*models.Card),
			OutCards:           make(map[int]map[string][]*TurnUserCard),
			FraudCard:          make(map[int]*models.Card),
			likeCards:          make(map[int][]*models.LikeCard),
			likeUserInfo:       make(map[int]map[string][]*models.LikeCard),
			LastUserOutCards:   make(map[string]*TurnUserCard),
			IsReGiveCard:       make(map[int]bool),
			SelectIssue:        make(map[int]*models.Issue),
			ExtractCard:        make(map[int]map[string][]*table.MbCardConfig),
			NotExtractCard:     make(map[int]map[string][]*table.MbCardConfig),
			UserOwnCards:       make(map[string][]*table.MbCardConfig),
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

func (rs *ComRoomSpace) SetLikeUserInfo(userId string, likeCards *models.LikeCard) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()
	turn := rs.GetTurn()

	_, ok := rs.TurnMateInfo.likeUserInfo[turn]
	if !ok {
		rs.TurnMateInfo.likeUserInfo[turn] = make(map[string][]*models.LikeCard)
	}
	rs.TurnMateInfo.likeUserInfo[turn][userId] = append(rs.TurnMateInfo.likeUserInfo[turn][userId], likeCards)
}

func (rs *ComRoomSpace) GetAllLikeUserInfo() []*models.LikeCard {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()
	turn := rs.GetTurn()
	res := make([]*models.LikeCard, 0)

	turnLikeUserInfos, ok := rs.TurnMateInfo.likeUserInfo[turn]
	if !ok {
		return res
	}

	for _, turnLikeUserInfo := range turnLikeUserInfos {
		res = append(res, turnLikeUserInfo...)
	}
	return res
}

// GetLikeUserInfo 某个用户 某一轮 给别人点赞的纪录
func (rs *ComRoomSpace) GetLikeUserInfo(userId string) []*models.LikeCard {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()
	turn := rs.GetTurn()
	res := make([]*models.LikeCard, 0)

	turnLikeUserInfo, ok := rs.TurnMateInfo.likeUserInfo[turn]
	if !ok {
		return res
	}

	cards, okk := turnLikeUserInfo[userId]
	if !okk {
		return res
	}

	return cards
}

func (rs *ComRoomSpace) AddUserOwnCards(userId string, userCards []*table.MbCardConfig) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	_, ok := rs.TurnMateInfo.UserOwnCards[userId]
	if !ok {
		rs.TurnMateInfo.UserOwnCards[userId] = make([]*table.MbCardConfig, 0)
	}
	rs.TurnMateInfo.UserOwnCards[userId] = userCards
	return
}

func (rs *ComRoomSpace) GetUserOwnCards(userId string) []*table.MbCardConfig {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()
	res := make([]*table.MbCardConfig, 0)
	userCards, ok := rs.TurnMateInfo.UserOwnCards[userId]
	if !ok {
		return res
	}
	res = userCards
	return res
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

// AddFraudCard 添加骗子牌
func (rs *ComRoomSpace) AddFraudCard(card *models.Card) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()
	//已经存在骗子牌
	_, ok := rs.TurnMateInfo.FraudCard[rs.TurnMateInfo.Turn]
	if ok {
		return
	}
	rs.TurnMateInfo.FraudCard[rs.TurnMateInfo.Turn] = card
}

func (rs *ComRoomSpace) AddSelectIssue(issue *models.Issue) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()
	//已经存在问题
	_, ok := rs.TurnMateInfo.SelectIssue[rs.TurnMateInfo.Turn]
	if ok {
		return
	}
	rs.TurnMateInfo.SelectIssue[rs.TurnMateInfo.Turn] = issue
}

// GetFraudCard 本轮是否存在骗子牌
func (rs *ComRoomSpace) GetFraudCard() (*models.Card, error) {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	//已经存在骗子牌
	card, ok := rs.TurnMateInfo.FraudCard[rs.TurnMateInfo.Turn]
	if !ok {
		//没有数据
		return nil, errors.New("没有骗子牌")
	}
	return card, nil
}

func (rs *ComRoomSpace) GetSelectIssue() (*models.Issue, error) {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	//已经存在骗子牌
	issue, ok := rs.TurnMateInfo.SelectIssue[rs.TurnMateInfo.Turn]
	if !ok {
		//没有数据
		return nil, errors.New("没有问题")
	}
	return issue, nil
}

// IsDealCards 本轮是否发过牌
func (rs *ComRoomSpace) IsDealCards() bool {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	//已经存在骗子牌
	v, ok := rs.TurnMateInfo.Cards[rs.TurnMateInfo.Turn]
	if ok && len(v) > 0 {
		return true
	}
	return false
}

func (rs *ComRoomSpace) IsAutoDealCards() bool {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	v, ok := rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn]
	if ok && len(v) > 0 {
		return true
	}
	return false
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

//func (rs *ComRoomSpace) NatsSendAimUserMsg(msg *pbs.NetMessage, userId string) {
//	//给客户消息
//	userInfo, ok := rs.UserInfos[userId]
//	if ok {
//		global.GVA_LOG.Infof("NatsSendAimUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, msg)
//		if len(userInfo.UserID) == 0 || userInfo.UserProperty.IsLeave == 1 || userInfo.UserProperty.IsRobot == 1 {
//			return
//		}
//		//userIDInt, err := strconv.Atoi(userInfo.UserID)
//		//if err != nil {
//		//	global.GVA_LOG.Error("NatsSendAimUserMsg err", zap.Any("err", err))
//		//	return
//		//}
//		msg.AckHead.Uid = userId
//		netMessageRespMarshal, _ := proto.Marshal(msg)
//		global.GVA_LOG.Infof("NatsSendAimUserMsg UserID:{%v} 给客户端发消息:{%v}", userInfo.UserID, msg)
//		NastManager.Producer(netMessageRespMarshal)
//	}
//
//}

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

// SendAimUserMsg 发给目标用户
//func (rs *ComRoomSpace) SendAimUserMsg(msgBty []byte, userId string) {
//	//给客户消息
//	client, ok := rs.UserClient[userId]
//	if ok {
//		global.GVA_LOG.Infof("SendAimUserMsg UserID:{%v} 给客户端发消息:{%v}", client.UserID, string(msgBty))
//		//client.SendMsg(msgBty)
//	}
//}

// GetCurrCard 获取用户当前的牌
func (rs *ComRoomSpace) GetCurrCard(userId string) ([]*models.Card, error) {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	var cards []*models.Card

	cards, ok := rs.TurnMateInfo.Cards[rs.TurnMateInfo.Turn][userId]
	global.GVA_LOG.Infof("GetCurrCard 用户:%v 轮数:%v", userId, rs.TurnMateInfo.Turn)
	if !ok {
		return cards, errors.New("没有牌")
	}

	return cards, nil

}

// AddCurrCard 只保留最近一次的牌
func (rs *ComRoomSpace) AddCurrCard(userId string, cards []*models.Card) error {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	_, ok := rs.TurnMateInfo.Cards[rs.TurnMateInfo.Turn]
	if !ok {
		rs.TurnMateInfo.Cards[rs.TurnMateInfo.Turn] = make(map[string][]*models.Card)
	}
	rs.TurnMateInfo.Cards[rs.TurnMateInfo.Turn][userId] = cards

	return nil
}

// CurrTurnFirstNotExtractCard 当前轮第一次初始化没出过的牌
func (rs *ComRoomSpace) CurrTurnFirstNotExtractCard(userId string, roomBaseCard []*table.MbCardConfig) error {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	_, ok := rs.TurnMateInfo.NotExtractCard[rs.TurnMateInfo.Turn]
	if !ok {
		rs.TurnMateInfo.NotExtractCard[rs.TurnMateInfo.Turn] = make(map[string][]*table.MbCardConfig)
	}
	rs.TurnMateInfo.NotExtractCard[rs.TurnMateInfo.Turn][userId] = roomBaseCard

	return nil
}

// GetNotExtractCard 获取当前轮 没有被随的牌
func (rs *ComRoomSpace) GetNotExtractCard(userId string) ([]*table.MbCardConfig, error) {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	data, ok := rs.TurnMateInfo.NotExtractCard[rs.TurnMateInfo.Turn]
	if !ok {
		return nil, errors.New("not carts ")
	}
	configs, ok := data[userId]
	if !ok {
		return nil, errors.New("not carts ")
	}
	return configs, nil
}

// ReMakeExtractCard 重置 未抽过的牌
func (rs *ComRoomSpace) ReMakeExtractCard(userId string, newCards []*table.MbCardConfig) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	turn := rs.GetTurn()

	_, ok := rs.TurnMateInfo.NotExtractCard[turn]
	if !ok {
		rs.TurnMateInfo.NotExtractCard[turn] = make(map[string][]*table.MbCardConfig)
	}
	rs.TurnMateInfo.NotExtractCard[turn][userId] = newCards
}

// SaveExtractCard 纪录已经抽过的牌
func (rs *ComRoomSpace) SaveExtractCard(userId string, outCards []*table.MbCardConfig) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	_, ok := rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn]
	if !ok {
		rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn] = make(map[string][]*table.MbCardConfig)
	}
	rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn][userId] = append(rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn][userId], outCards...)
}

func (rs *ComRoomSpace) GetExtractCard(userId string) []*table.MbCardConfig {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	date, ok := rs.TurnMateInfo.ExtractCard[rs.TurnMateInfo.Turn]
	if !ok {
		return make([]*table.MbCardConfig, 0)
	}
	configs, ok := date[userId]
	if !ok {
		return make([]*table.MbCardConfig, 0)
	}
	return configs
}

// ReMakeCurrCard 重置当前手中的卡
func (rs *ComRoomSpace) ReMakeCurrCard(userId string, reCards, outCards []*models.Card) {
	rs.TurnMateInfo.TurnSync.Lock()
	defer rs.TurnMateInfo.TurnSync.Unlock()

	turn := rs.GetTurn()
	//手里的牌
	rs.TurnMateInfo.Cards[turn][userId] = reCards

	//出牌
	turnUserCard := TurnUserCard{
		Ord:   turn,
		Cards: outCards,
	}

	_, ok := rs.TurnMateInfo.OutCards[turn]
	if !ok {
		rs.TurnMateInfo.OutCards[turn] = make(map[string][]*TurnUserCard)
	}

	outTurnUserCard, ok := rs.TurnMateInfo.OutCards[turn][userId]
	if !ok {
		rs.TurnMateInfo.OutCards[turn][userId] = append(rs.TurnMateInfo.OutCards[turn][userId], &turnUserCard)
	} else {
		turnUserCard.Ord = len(outTurnUserCard)
		rs.TurnMateInfo.OutCards[turn][userId] = append(rs.TurnMateInfo.OutCards[turn][userId], &turnUserCard)
	}

	rs.TurnMateInfo.LastUserOutCards[userId] = &turnUserCard
}

// IsAllUserOutCart 本轮是否全部用户出过牌
func (rs *ComRoomSpace) IsAllUserOutCart() bool {
	OutCardsMap, ok := rs.TurnMateInfo.OutCards[rs.GetTurn()]
	if !ok {
		//没有当前轮信息
		global.GVA_LOG.Error("IsAllUserOutCart 本轮是否全部用户出过牌")
		return false
	}

	if len(OutCardsMap) != len(rs.UserInfos) {
		return false
	}
	return true
}

// GetUserOutEdCards 获取用户已经出国的牌
func (rs *ComRoomSpace) GetUserOutEdCards(userId string) []*models.Card {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	var resCards []*models.Card

	_, ok := rs.TurnMateInfo.OutCards[rs.GetTurn()]
	if !ok {
		return resCards
	}

	_, ok = rs.TurnMateInfo.OutCards[rs.GetTurn()][userId]
	if !ok {
		return resCards
	} else {
		userCards := rs.TurnMateInfo.OutCards[rs.GetTurn()][userId]
		for _, card := range userCards {
			resCards = append(resCards, card.Cards...)
		}
	}

	return resCards
}

// GetAllUserOutEdCards 所有用户 本轮出牌的集合
func (rs *ComRoomSpace) GetAllUserOutEdCards() []*models.Card {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	var resCards []*models.Card

	for _, uInfo := range rs.UserInfos {
		cards := rs.GetUserOutEdCards(uInfo.UserID)
		for _, card := range cards {
			card.UserID = uInfo.UserID
			resCards = append(resCards, card)
		}
	}
	return resCards
}

// GetUserOutEdCardExcludeUser 排除某个用户 本轮出牌的集合
func (rs *ComRoomSpace) GetUserOutEdCardExcludeUser(userId string) []*models.Card {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()

	var resCards []*models.Card

	for _, uInfo := range rs.UserInfos {
		if userId == uInfo.UserID {
			continue
		}
		cards := rs.GetUserOutEdCards(uInfo.UserID)
		for _, card := range cards {
			card.UserID = uInfo.UserID
			resCards = append(resCards, card)
		}
	}
	return resCards
}

// GetLastOutCard 获取上次出的牌
func (rs *ComRoomSpace) GetLastOutCard(userId string) (*TurnUserCard, error) {
	rs.TurnMateInfo.TurnSync.RLock()
	defer rs.TurnMateInfo.TurnSync.RUnlock()
	//userId是上次出牌的 uid

	//先从冗余数据里面查一下数据 有的话直接返回
	lastUserOutCard, ok := rs.TurnMateInfo.LastUserOutCards[userId]
	global.GVA_LOG.Infof("GetLastOutCard 获取上次出的牌 %v ,%v", lastUserOutCard, ok)
	if ok == false {
		return lastUserOutCard, nil
	}

	//出牌
	lastUserOutCards, ok := rs.TurnMateInfo.OutCards[rs.GetTurn()][userId]
	global.GVA_LOG.Infof("GetLastOutCard 获取上次出的牌 lastTurnUserCards %v ,%v", lastUserOutCards, ok)
	if !ok {
		return nil, errors.New("上一个用户没有出过牌,无效质疑")
	}

	lastUserOutCard = lastUserOutCards[len(lastUserOutCards)-1]

	return lastUserOutCard, nil

}

//func (rs *ComRoomSpace) SetUserHandarmConfig() {
//	rs.TurnMateInfo.TurnSync.Lock()
//	defer rs.TurnMateInfo.TurnSync.Unlock()
//
//	for k, _ := range rs.UserInfos {
//		userInfo := rs.UserInfos[k]
//		//开始游戏前 初始化一下用户的左轮枪信息
//		userInfo.UserExt.HandarmConfig = models.NewHandarmConfig(helper.RandInt(6)+1, 0, 6)
//	}
//}

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

// GetTurnFraudCard 获取当前轮的骗子牌
func (rs *ComRoomSpace) GetTurnFraudCard() *models.Card {
	card := rs.TurnMateInfo.FraudCard[rs.GetTurn()]
	return card
}

// GetTurnCards 用户手里的当前牌
func (rs *ComRoomSpace) GetTurnCards(userId string) []*models.Card {
	res := make([]*models.Card, 0)
	turnCardInfoMap, ok := rs.TurnMateInfo.Cards[rs.GetTurn()]
	if !ok {
		//没有当前轮信息
		global.GVA_LOG.Error("GetTurnCards 没有当前轮信息", zap.Any("GetTurn", userId))
		return res
	}
	cards, ook := turnCardInfoMap[userId]
	if !ook {
		//没有这个用户
		global.GVA_LOG.Error("GetTurnCards 没有这个用户", zap.Any("userId", userId))
		return res
	}

	res = cards
	return cards
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

// SetNextOutCardUser 下一轮要出牌的人 出牌的时候
//func (rs *ComRoomSpace) SetNextOutCardUser(userId string) (*models.UserInfo, error) {
//	rs.UserSync.RLock()
//	defer rs.UserSync.RUnlock()
//
//	_, ok := rs.UserInfos[userId]
//	if !ok {
//		return nil, errors.New("user not found")
//	}
//
//	nextUser, err := rs.FindNextUser(userId)
//	if err != nil {
//		return nil, err
//	}
//	nextSeat := nextUser.UserProperty.Seat
//
//	//当前用户的座位
//	//nextSeat := userInfo.UserProperty.Seat + 1
//	//if nextSeat > len(rs.UserInfos) {
//	//	nextSeat = 1
//	//}
//
//	//找到下一个出牌的人
//	resUser := &models.UserInfo{}
//	for _, uInfo := range rs.UserInfos {
//		if uInfo.UserProperty.Seat == nextSeat {
//			uInfo.SetUserIsMyTurn(true)
//			//服务端计算 40秒被出牌
//			uInfo.SetOutCardCountDown(models.GetOutCardCountDownTimeInt(OutCardCountDownTimeInt))
//			resUser = uInfo
//		} else {
//			uInfo.SetUserIsMyTurn(false)
//		}
//	}
//
//	//发送广播
//	//msgData := models.NextOutCardUserMsg{
//	//	ProtoNum:  models.TavernNextOutCardUser,
//	//	Timestamp: time.Now().Unix(),
//	//	UserId:    resUser.UserID,
//	//}
//	//responseHead := models.NewResponseHead("", models.TavernNextOutCardUser, common.OK, "", msgData)
//	//responseHeadByte, _ := json.Marshal(responseHead)
//
//	//给客户消息
//	//global.GVA_LOG.Infof("SetNextOutCardUser 下一个出牌的用户: %v", string(responseHeadByte))
//	//
//	//rs.SendAllUserMsg(responseHeadByte)
//
//	return resUser, nil
//}

// DoubtSetNextOutCardUser 质疑的时候 ：userId 失败用户ID
//func (rs *ComRoomSpace) DoubtSetNextOutCardUser(userId string) (*models.UserInfo, error) {
//	rs.UserSync.RLock()
//	defer rs.UserSync.RUnlock()
//
//	userInfo, ok := rs.UserInfos[userId]
//	if !ok {
//		return nil, errors.New("user not found")
//	}
//
//	//找到下一个出牌的人
//	resUser := &models.UserInfo{}
//
//	if !userInfo.UserIsKilled() {
//		//用户未死 下次就是用户出牌
//		for _, uInfo := range rs.UserInfos {
//			if userId == uInfo.UserID {
//				uInfo.SetUserIsMyTurn(true)
//				//服务端计算 40秒被出牌
//				uInfo.SetOutCardCountDown(models.GetOutCardCountDownTimeInt(OutCardCountDownTimeInt))
//				resUser = uInfo
//			} else {
//				uInfo.SetUserIsMyTurn(false)
//			}
//		}
//	} else {
//		//用户已死
//		userInfo.SetUserIsMyTurn(false)
//
//		//通过当前座位 找到下一个座位的用户
//		user, err := rs.FindNextUser(userId)
//		if err != nil {
//			return resUser, err
//		}
//		resUser = user
//		resUser.SetUserIsMyTurn(true)
//	}
//
//	return resUser, nil
//}

// FindNextUser 找到下一个活着的用户
//func (rs *ComRoomSpace) FindNextUser(userId string) (*models.UserInfo, error) {
//	//userId 失败用户ID
//	userInfo, ok := rs.UserInfos[userId]
//	if !ok {
//		return nil, errors.New("user not found")
//	}
//
//	//没有被杀的人数量
//	var notBeKillNum int
//	for k, _ := range rs.UserInfos {
//		if !rs.UserInfos[k].UserIsKilled() {
//			notBeKillNum++
//		}
//	}
//	if notBeKillNum == 1 {
//		return userInfo, errors.New("没有下一个用户 本局结束")
//	}
//
//	//找到下一个出牌的人
//	resUser := &models.UserInfo{}
//
//	//算法逻辑
//	//要么就是比当前座位大的最近一个座位 要么就是最小的座位
//	currSeat := userInfo.UserProperty.Seat
//
//	//通过当前座位 找到下一个座位的用户
//	seatInt := []int{}
//	notBeKillUser := []*models.UserInfo{}
//	for _, uInfo := range rs.UserInfos {
//		if !uInfo.UserIsKilled() {
//			notBeKillUser = append(notBeKillUser, uInfo)
//			seatInt = append(seatInt, uInfo.UserProperty.Seat)
//		}
//	}
//
//	//升序
//	sort.Ints(seatInt)
//	nextSeat := 0
//	for _, sInt := range seatInt {
//		if sInt > currSeat {
//			nextSeat = sInt
//			break
//		}
//	}
//	if nextSeat == 0 {
//		nextSeat = seatInt[0]
//	}
//
//	for _, uInfo := range notBeKillUser {
//		if uInfo.UserProperty.Seat == nextSeat {
//			resUser = uInfo
//		}
//	}
//
//	return resUser, nil
//}

//func (rs *ComRoomSpace) FindPreUser(userId string) (*models.UserInfo, error) {
//	userInfo, ok := rs.UserInfos[userId]
//	if !ok {
//		return nil, errors.New("user not found")
//	}
//
//	//没有被杀的人数量
//	var notBeKillNum int
//	for k, _ := range rs.UserInfos {
//		if !rs.UserInfos[k].UserIsKilled() {
//			notBeKillNum++
//		}
//	}
//	if notBeKillNum == 1 {
//		return userInfo, errors.New("没有下一个用户 本局结束")
//	}
//
//	//找到下一个出牌的人
//	resUser := &models.UserInfo{}
//
//	//算法逻辑
//	//要么就是比当前座位小的最近一个座位 要么就是最大的座位
//	currSeat := userInfo.UserProperty.Seat
//
//	//通过当前座位 找到下一个座位的用户
//	seatInt := []int{}
//	notBeKillUser := []*models.UserInfo{}
//	for _, uInfo := range rs.UserInfos {
//		if !uInfo.UserIsKilled() {
//			notBeKillUser = append(notBeKillUser, uInfo)
//			seatInt = append(seatInt, uInfo.UserProperty.Seat)
//		}
//	}
//
//	//降序排序
//	sort.Slice(seatInt, func(i, j int) bool {
//		return seatInt[i] > seatInt[j]
//	})
//
//	preSeat := 0
//	for _, sInt := range seatInt {
//		if currSeat > sInt {
//			preSeat = sInt
//			break
//		}
//	}
//	if preSeat == 0 {
//		preSeat = seatInt[0]
//	}
//
//	for _, uInfo := range notBeKillUser {
//		if uInfo.UserProperty.Seat == preSeat {
//			resUser = uInfo
//		}
//	}
//
//	return resUser, nil
//}

// GetPreOutCardUser 根据当前用户获取前一个用户
//func (rs *ComRoomSpace) GetPreOutCardUser(userId string) (*models.UserInfo, error) {
//	rs.UserSync.RLock()
//	defer rs.UserSync.RUnlock()
//
//	userInfo, ok := rs.UserInfos[userId]
//	if !ok {
//		return &models.UserInfo{}, errors.New("user not found")
//	}
//
//	//当前用户的座位
//	preSeat := userInfo.UserProperty.Seat - 1
//	if preSeat <= 0 {
//		preSeat = len(rs.UserInfos)
//	}
//
//	//找到下一个出牌的人
//	resUser := &models.UserInfo{}
//	for _, uInfo := range rs.UserInfos {
//		if uInfo.UserProperty.Seat == preSeat {
//			resUser = uInfo
//			break
//		}
//	}
//
//	return resUser, nil
//}

// isEndGame 是否结束游戏
//func (rs *ComRoomSpace) isEndGame() bool {
//	rs.TurnMateInfo.TurnSync.Lock()
//	defer rs.TurnMateInfo.TurnSync.Unlock()
//
//	//是否结束游戏
//	var isEndGame bool
//
//	//没有被杀的人数量
//	var notBeKillNum int
//	for k, _ := range rs.UserInfos {
//		if rs.UserInfos[k].GetUserIsKilled() != 1 {
//			notBeKillNum++
//		}
//	}
//
//	//就剩一个人 游戏结束
//	if notBeKillNum == 1 {
//		isEndGame = true
//	}
//
//	rs.SetStopGame(isEndGame)
//
//	return isEndGame
//}

// DoubtResult 质疑的结果 谁输 谁赢
//func (rs *ComRoomSpace) DoubtResult(doubtUserInfo, beDoubtUserInfo *models.UserInfo, turnUserCard *TurnUserCard, fraudCard *models.Card) (loserUid string, isHaveBullet bool) {
//	var isHaveThirdCardId bool
//	//如果
//	//出的是骗子牌 被质疑者输掉比赛
//	//规则 出的牌里面只有骗子牌ID和鬼牌ID 质疑者输，否则被质疑者输
//	for k, _ := range turnUserCard.Cards {
//		//4 鬼牌ID 写死 todo
//		if turnUserCard.Cards[k].CardId == fraudCard.CardId || turnUserCard.Cards[k].CardId == 4 {
//			continue
//		}
//		isHaveThirdCardId = true
//	}
//
//	//质疑者输
//	//var loserUid string
//	////是否有子弹
//	//var isHaveBullet bool
//
//	if !isHaveThirdCardId {
//		//没有除鬼牌 和骗子牌以外的其他牌
//		loserUid = doubtUserInfo.UserID
//
//		//看左轮枪是否有子弹
//		doubtUserInfo.AddUserHandarmCurrSeat()
//		if doubtUserInfo.UserSeatIsEquCurrSeat() {
//			//中枪
//			isHaveBullet = true
//			doubtUserInfo.SetUserIsKilled(1)
//			//err := dao.UpdateTavernRoomIsKill(loserUid, doubtUserInfo.UserExt.RoomNo)
//			//if err != nil {
//			//	global.GVA_LOG.Error("DoubtResult", zap.Any("err", err))
//			//}
//			//
//			//updateMap := dao.MakeUpdateData("is_killed", 1)
//			//updateMap["room_no"] = doubtUserInfo.UserExt.RoomNo
//			//dao.UpdateTavernUsersRoomRoomNo(doubtUserInfo.UserID, updateMap)
//		}
//
//	} else {
//		//被质疑者输
//		loserUid = beDoubtUserInfo.UserID
//
//		//看左轮枪是否有子弹
//		beDoubtUserInfo.AddUserHandarmCurrSeat()
//		if beDoubtUserInfo.UserSeatIsEquCurrSeat() {
//			//中枪
//			isHaveBullet = true
//			//beDoubtUserInfo.UserProperty.IsKilled = 1
//			//beDoubtUserInfo.SetUserIsKilled(1)
//			//err := dao.UpdateTavernRoomIsKill(loserUid, beDoubtUserInfo.UserExt.RoomNo)
//			//if err != nil {
//			//	global.GVA_LOG.Error("DoubtResult", zap.Any("err", err))
//			//}
//			//
//			//updateMap := dao.MakeUpdateData("is_killed", 1)
//			//updateMap["room_no"] = beDoubtUserInfo.UserExt.RoomNo
//			//dao.UpdateTavernUsersRoomRoomNo(beDoubtUserInfo.UserID, updateMap)
//		}
//	}
//	return loserUid, isHaveBullet
//}
