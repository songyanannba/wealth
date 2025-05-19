package websocket

import (
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/helper"
	"slot_server/lib/models/table"
)

type GameTurnState int

// 游戏开始前 -> 问题阶段 -> 随牌阶段  -> 出牌阶段 -> 点赞阶段 -> 计算阶段
// 随牌阶段：1，端侧倒计时15秒 （服务端倒计时25秒：待定） ；2，触发条件：全部加载完的时候 ，结束条件：是倒计时结束的时候
// 出牌阶段：1，端侧倒计时15秒（真实是15+15） 服务端倒计时25秒（真实是15+25） ；2，触发条件：随牌结束的时候 ，结束条件：全部用户出完牌 或者是倒计时结束的时候
// 点赞阶段：1，端侧倒计时10秒 服务端倒计时20秒 ；2，触发条件：出牌结束 ，结束条件：倒计时结束 或者全部用户点赞
// 计算阶段：触发条件：最后一轮全部点赞完成的时候

// 游戏状态
const (
	Game GameTurnState = iota

	/**
	游戏开始阶段
	EnGameStartExec 触发条件：当房主点击开始的时候
	*/

	EnGameStartExec //游戏开始 还没执行发送广播

	EnGameStartIng //游戏开始 已经发送完广播 在开始中 还有用户没有加载

	//BetIng 押注阶段
	BetIng

	//EnWheelAnimalPartyCalculateExec 押注结束 给端侧发送动物排序
	EnWheelAnimalPartyCalculateExec

	//EnAnimalPartyCalculateExecIng 开始执行计算
	EnAnimalPartyCalculateExecIng
	EnAnimalPartyCalculateExec

	StartNextPeriod

	//EnLoadExec 加载
	EnLoadExec // 加载阶段

	RemakeCardIng //随牌阶段

	// EnLikePageExec 点赞阶段
	EnLikePageExec
	EnLikeCardIng //点赞进行中

	//EnNextTurnExec 计算或者结束阶段
	EnNextTurnExec  //进入下一轮阶段 进入下轮轮后是随牌阶段
	EnCalculateExec //计算并结束结束阶段

	GameOver //本局结束
)

type ClientGameTurnState int

const (
	//0=游戏未开始
	//1=游戏开始但是没有加载完成
	//2=用户随牌阶段
	//3=用户出牌阶段
	//4=用户点赞阶段
	//5=点赞界面 等待结算或者进入下一轮

	CliGame ClientGameTurnState = iota
	CliStartAndLongIng
	CliRemakeCard
	CliOutCard
	CliLikePage
	CliNextTurnOrCalculate
)

// CurrGameTurnStateAndDownTime 获取当前的游戏阶段和当前阶段的倒计时
func (trs *RoomSpace) CurrGameTurnStateAndDownTime() (int, int64) {
	var (
		turnState         int
		turnTime          int64
		currTime          = helper.LocalTime().Unix()
		countdownTime     = trs.ComRoomSpace.GetCountdownTime()
		likeCountdownTime = trs.ComRoomSpace.GetLikeCountdownTime()
	)
	//这个要转换成返回 端侧的状态
	state := trs.ComRoomSpace.GetGameState()

	//游戏开始 但是没有加载完成
	if state == EnGameStartExec || state == EnGameStartIng || state == EnLoadExec {
		//加载 状态
		turnState = int(CliStartAndLongIng)
		return turnState, turnTime
	}

	//随牌阶段 + 出牌阶段
	//EnLikePageExec 说明都出过牌 进入到点赞页面的执行状态 但是定时器还没执行
	if state == RemakeCardIng || state == EnLikePageExec {
		//根据倒计时判断当前状态
		return RemakeCardIngAndOutCartIng(currTime, countdownTime)
	}

	//点赞阶段
	//EnNextTurnExec  EnCalculateExec 说明都已经点赞 但是下一个状态的方法还没执行
	if state == EnLikeCardIng {
		//进入点赞阶段
		return EnLikePageExecAndEnLikeCardIng(currTime, likeCountdownTime)
	}

	//在点赞页面等待 等待进入下一轮 或者结算的消息
	if state == EnNextTurnExec || state == EnCalculateExec {
		turnState = int(CliNextTurnOrCalculate)
		return turnState, turnTime
	}

	return turnState, turnTime
}

func (trs *RoomSpace) ServerSimplifyGetStateAndTime() (int, bool) {
	var (
		turnState int
		//当前时间
		currTime = helper.LocalTime().Unix()
		//进入每一轮的时间
		countdownTime = trs.ComRoomSpace.GetCountdownTime()
		//进入点在叶脉呢
		likeCountdownTime = trs.ComRoomSpace.GetLikeCountdownTime()
	)
	//这个要转换成返回 端侧的状态
	state := trs.ComRoomSpace.GetGameState()

	//随牌阶段 + 出牌阶段
	if (state == RemakeCardIng || state == EnLikePageExec) && countdownTime != 0 {
		//根据倒计时判断当前状态
		gapTime := currTime - countdownTime
		turnState = int(CliOutCard)
		if gapTime < CommTimeOutDouble+CommTimeDelay {
			//不到托管时间
			return turnState, false
		} else {
			//触发托管时间 出牌的托管
			return turnState, true
		}
	}

	//点赞阶段
	if state == EnLikeCardIng && likeCountdownTime != 0 {
		//进入点赞阶段 等待消息 还没执行方法
		likeGapTime := currTime - likeCountdownTime
		turnState = int(CliLikePage)
		if likeGapTime < CommTimeOut {
			//不到托管时间
			return turnState, false
		} else {
			//触发托管时间 点在的托管
			return turnState, true
		}
	}
	return turnState, false

}

func RemakeCardIngAndOutCartIng(currTime, countdownTime int64) (int, int64) {
	var (
		turnState int
		turnTime  int64
	)
	//每轮的开始时间
	gapTime := currTime - countdownTime
	global.GVA_LOG.Infof("RemakeCardIngAndOutCartIng gapTime:%v", gapTime)

	// 和当前时间比较
	//【0 -- 15】秒内就是随牌阶段
	if gapTime >= 0 && gapTime < CommTimeOut {
		turnState = int(CliRemakeCard)
		turnTime = gapTime
		return turnState, turnTime
	}

	turnState = int(CliOutCard)

	//【15 -- 30】内就是出牌阶段
	if gapTime >= CommTimeOut && gapTime < CommTimeOutDouble {
		turnState = int(CliOutCard)
		turnTime = gapTime
	}

	//【30 -- 40】内是服务端强制执行阶段 ，端侧不做操作
	//if gapTime >= 30 && gapTime < 40 {
	//	turnState = int(CliOutCard)
	//	//没有倒计时 客户端就停在 出牌阶段就可以
	//}

	return turnState, turnTime
}

func EnLikePageExecAndEnLikeCardIng(currTime, likeCountdownTime int64) (int, int64) {
	var (
		turnState int
		turnTime  int64
	)

	likeGapTime := currTime - likeCountdownTime
	global.GVA_LOG.Infof("CurrGameTurnStateAndDownTime EnLikePageExecAndEnLikeCardIng like gapTime:%v", likeGapTime)

	turnState = int(CliLikePage)

	//进入点赞阶段 等待消息 还没执行方法
	if likeCountdownTime == 0 {
		//
	} else {
		// 和当前时间比较
		//【0 -- 15】秒内就是点赞倒计时
		if likeGapTime >= 0 && likeGapTime < CommTimeOut {
			turnTime = likeGapTime
		}
		//【15 -- 25】内是服务端强制执行阶段 ，端侧不做操作
		if likeGapTime >= CommTimeOut && likeGapTime < (CommTimeOut+CommTimeDelay) {
			//turnTime = likeGapTime
		}

	}
	return turnState, turnTime
}

func (trs *RoomSpace) RegisterTurnStateFunc(key GameTurnState, GameStateFunc func(trs *RoomSpace)) {
	trs.GameStateFuncMapMutex.Lock()
	defer trs.GameStateFuncMapMutex.Unlock()
	trs.GameStateMap[key] = GameStateFunc
	return
}

func (trs *RoomSpace) GetTurnStateFuncHandlers(key GameTurnState) (value GameStateFunc, ok bool) {
	trs.GameStateFuncMapMutex.RLock()
	defer trs.GameStateFuncMapMutex.RUnlock()
	value, ok = trs.GameStateMap[key]
	return
}

func (trs *RoomSpace) ExecProcessTurnStateFunc(key GameTurnState) {
	global.GVA_LOG.Infof("ExecProcessTurnStateFunc game State 游戏状态%v", key)
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("ExecProcessTurnStateFunc 处理数据 stop", zap.Any("", r))
		}
	}()

	// 采用 map 注册的方式
	if value, ok := trs.GetTurnStateFuncHandlers(key); ok {
		value(trs)
		trs.ExecAutoNextTurnState(key)
	} else {
		//global.GVA_LOG.Error("ExecProcessTurnStateFunc 处理数据 路由不存在", zap.Any("key", key))
		return
	}
	//global.GVA_LOG.Infof("RoomNo{%v},处理 comMsg.MsgId %v 返回数据data:%v ", trs.RoomInfo.RoomNo)
	return
}

func (trs *RoomSpace) InItTurnStateFunc() {
	//房主开始游戏的执行逻辑
	//trs.RegisterTurnStateFunc(EnGameStartExec, SendGameStartBroadcast)
	////记载完成的执行逻辑
	//trs.RegisterTurnStateFunc(EnLoadExec, SendEnLoadBroadcast)
	////都出过牌的时候 发送进入点赞页面的广播
	//trs.RegisterTurnStateFunc(EnLikePageExec, EntryLikePage)
	////点赞结束 进入下一轮
	//trs.RegisterTurnStateFunc(EnNextTurnExec, NextTurnExecFunc)
	//点赞结束 计算
	//trs.RegisterTurnStateFunc(EnCalculateExec, CalculateExecFunc)

	//发送动物排序
	trs.RegisterTurnStateFunc(EnWheelAnimalPartyCalculateExec, WheelAnimalSortCalculateExec)

	//计算
	trs.RegisterTurnStateFunc(EnAnimalPartyCalculateExec, WheelAnimalPartyCalculateExec)
}

func (trs *RoomSpace) ExecAutoNextTurnState(key GameTurnState) {
	//if key == EnGameStartExec {
	//	//游戏的开始状态 只能进入游戏开始结束状态
	//	//DisGameStart 功能是发送广播
	//	trs.ComRoomSpace.GameStateTransition(EnGameStartExec, EnGameStartIng)
	//} else if key == EnLoadExec {
	//	//
	//	trs.ComRoomSpace.GameStateTransition(EnLoadExec, RemakeCardIng)
	//} else if key == EnLikePageExec {
	//	//
	//	trs.ComRoomSpace.GameStateTransition(EnLikePageExec, EnLikeCardIng)
	//} else if key == EnNextTurnExec {
	//	//
	//	trs.ComRoomSpace.GameStateTransition(EnNextTurnExec, RemakeCardIng)
	//} else if key == EnCalculateExec {
	//	//最后一轮进入结算状态
	//trs.ComRoomSpace.GameStateTransition(EnCalculateExec, GameOver)
	//}

	if key == EnWheelAnimalPartyCalculateExec {
		trs.ComRoomSpace.GameStateTransition(EnWheelAnimalPartyCalculateExec, EnAnimalPartyCalculateExec)
	}

	if key == EnAnimalPartyCalculateExec {
		trs.CloseRoom(trs.RoomInfo.Name, table.RoomStatusStop)
		trs.RoomInfo.IsOpen = table.RoomStatusIng
		trs.ComRoomSpace.GameStateTransition(EnAnimalPartyCalculateExec, StartNextPeriod)
	}
}

// ChangeGameState 改变游戏状态
func (rs *ComRoomSpace) ChangeGameState(state GameTurnState) {
	rs.SetGameState(state)
}

func (rs *ComRoomSpace) GameStateTransition(from, to GameTurnState) bool {
	if rs.GetGameState() != from {
		global.GVA_LOG.Infof("GameStateTransition 失败 from:{%v} to:{%v} ", from, to)
		return false
	} else {
		global.GVA_LOG.Infof("GameStateTransition 成功 from:{%v} to:{%v}", from, to)
		rs.ChangeGameState(to)
		return true
	}
}

func (rs *ComRoomSpace) GetGameState() GameTurnState {
	return rs.TurnMateInfo.GameTurnStatus
}

func (rs *ComRoomSpace) SetGameState(status GameTurnState) {
	rs.TurnMateInfo.GameTurnStatus = status
}

type GameState int

//
//const (
//	GameStateDef GameState = iota
//	//AllUserLoad 优化
//	AllUserLoad     //全部用户加载完成的时候
//	AllUserLoadOver //广播 全部用户加载完成的时候
//)

// ExecFuncByGameState 根据当前状态 判断要执行的方法
//func (trs *RoomSpace) ExecFuncByGameState() {
//
//	switch gState {
//	case DisLoad:
//	case EnLikePage: //进入点赞
//		//如果是全部用户已经出牌 就是进入点赞页面
//		trs.EntryLikePage()
//	case EnLikeCard: //进入下一轮/游戏结束
//		//点赞阶段 全部用户完成点赞的时候 ｜ 本轮是否都已点赞 如果是，进入下一轮或者游戏结束
//		trs.NextTurnOrCalculateAndEnd()
//	default:
//		global.GVA_LOG.Infof("ExecFuncByGameState gState 游戏状态 default")
//	}
//
//}
