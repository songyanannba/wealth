package meme_serve

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"sync"
	"time"
)

type GameStatus int

var (
	GameStatusInit     GameStatus = 1  //正在初始化中游戏还没开始
	GameBeforeStart    GameStatus = 2  //游戏即将开始 游戏开始前1到3分钟会收到信息
	GameBeforePrepare  GameStatus = 3  //游戏开始 还在准备阶段
	GameCalculateStart GameStatus = 4  //游戏开始 结算阶段 开始
	GameLayerPNoPass   GameStatus = 5  //游戏开始 结算阶段 本层不过关
	GameLayerPass      GameStatus = 6  //游戏开始 结算阶段 本层过关
	GameAttack         GameStatus = 7  //游戏开始 攻打阶段
	GameGoCrazy        GameStatus = 8  //狂暴阶段
	GameLayerEnd       GameStatus = 9  //每层结束的时间标识
	GamePass           GameStatus = 15 //通关
)

type memeRoomManager struct {
	//是否初始化开始
	IsInitStart bool
}

type RoomUser struct {
	Sync     *sync.RWMutex
	UsersMap map[string]*models.UserInfo
}

// MemeRoomManager 房间管理器
var MemeRoomManager = memeRoomManager{}

func (trMgr *memeRoomManager) Start() {
	trMgr.ManagerInitDBData()

	// 创建一个计时器
	serviceTimer := time.NewTicker(time.Second * 10)
	defer serviceTimer.Stop() //定时器不用了需要关闭

	serviceRoomTimer := time.NewTicker(time.Second * 200) //todo
	defer serviceRoomTimer.Stop()                         //定时器不用了需要关闭

	//匹配房间的定时器
	matchRoomTimer := time.NewTicker(time.Second * 5)
	defer matchRoomTimer.Stop()

	//这个定时器是分析房间游戏是否开始
	defer func() {
		if err := recover(); err != nil {
			global.GVA_LOG.Error("mtRoomManager :", zap.Any("recover ", err))
		}
		global.GVA_LOG.Info("mtRoomManager end")
	}()

	for {
		select {
		case <-serviceTimer.C:
			//打印一些调试信息
			global.GVA_LOG.Infof("房间管理器:计时器触发,房间用户数量:%d", 1)

		case <-serviceRoomTimer.C:

		case <-matchRoomTimer.C:

		}
	}
}

func (trMgr *memeRoomManager) AddManager(userID, nickname string) (err error) {
	if !trMgr.IsInitStart {
		global.GVA_LOG.Infof("AddManager 最新一期还没初始化 %v", userID)
		return
	}

	//添加房间用户
	//uidStr := userID
	////保留用户的基本信息
	//realUser := models.UserInfo{
	//	UserID:   uidStr,
	//	Nickname: nickname,
	//	UserProperty: models.UserProperty{
	//		Turn: 1,
	//		//HeardTime: time.Now().Unix(),
	//	},
	//	UserExt: models.UserExt{},
	//}
	return nil
}

func (trMgr *memeRoomManager) ManagerInitDBData() {

	trMgr.IsInitStart = true
}

// GetAllRoomUser 获取房间的用户
func (trMgr *memeRoomManager) GetAllRoomUser() []*models.UserInfo {

	return nil
}

func (trMgr *memeRoomManager) DelNotHeartRoomUser() {

}

func (trMgr *memeRoomManager) DelRoom() {
	//trMgr.Rooms = nil
}
