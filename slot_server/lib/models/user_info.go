package models

import (
	"fmt"
	"github.com/gorilla/websocket"
	"slot_server/lib/helper"
)

type ClientInfo struct {
	Addr          string          // 客户端地址
	Socket        *websocket.Conn // 用户连接
	Send          chan []byte     // 待发送的数据
	AppID         uint32          // 登录的平台ID app/web/ios
	UserID        string          // 用户ID，用户登录以后才有
	FirstTime     uint64          // 首次连接事件
	HeartbeatTime uint64          // 用户上次心跳时间
	LoginTime     uint64          // 登录时间 登录以后才有
	Token         string
	Nickname      string
	ProtocType    int // 0:json协议 1:protoc协议
}

type UserInfo struct {
	UserID       string       `json:"user_id"`
	Nickname     string       `json:"nickname"`
	UserProperty UserProperty `json:"user_property"` //属性
	UserExt      UserExt      `json:"user_ext"`      //扩展
	ClientInfo
}

func (t *UserInfo) Copy() *UserInfo {
	c := *t
	return &c
}

func NewUserInfo(userID string, nickname string, prop UserProperty, userExt UserExt) UserInfo {
	return UserInfo{
		UserID:       userID,
		Nickname:     nickname,
		UserProperty: prop,
		UserExt:      userExt,
	}
}

// UserProperty 用户的属性
type UserProperty struct {
	Turn         int     `json:"turn"`           //在第几轮
	IsLeave      int     `json:"is_leave"`       // 1 离开 （见好就收）
	IsOwner      bool    `json:"is_owner"`       //是否房主
	IsReady      int     `json:"is_ready"`       //是否就绪 0 没有就绪 1就绪
	PriorityAct  bool    `json:"priority_act"`   //优先出牌
	IsMyTurn     bool    `json:"is_my_turn"`     //是否轮到自己出牌
	Seat         int     `json:"seat"`           //位子、顺序 从0开始
	UserLimitNum int     `json:"user_limit_num"` //房间人数限制
	WinPrice     float64 `json:"win_price"`      //最终赢钱
	Bet          float64 `json:"bet"`            //押注
	CharacterId  int     `json:"character_Id"`   //角色ID
	GameStatus   int     `json:"game_status"`    //0=游戏为开始阶段 1=问题阶段 2=随牌阶段 3=出牌阶段 4=点赞阶段 5=计算阶段 6=游戏结束
	//机器人
	IsRobot    int `json:"is_robot"`    //是否机器人 0:真实用户 1:机器人
	RobotClass int `json:"robot_class"` //1:机器人1 2:机器人2 3:机器人3
	RobotStrategy
}

func (u *UserInfo) GetString() string {
	str := "用户ID:" + u.UserID + "用户昵称" + u.Nickname + "\n"

	if u.UserIsRobot() {
		str += "是机器人" + "\n"
		str += fmt.Sprintf("机器人类型是:%v", u.UserProperty.RobotClass) + "\n"
	}
	return str
}

// RobotStrategy 机器人策略
type RobotStrategy struct {
	OutCardCountDown int64 `json:"out_card_count_down"` //出牌倒计时 (机器人用)
	ReMakeCardDown   int64 `json:"re_make_card_down"`   //重置牌倒计时 (机器人用)
	LikeCardDown     int64 `json:"like_card_down"`      //点赞倒计时 (机器人用)
}

func NewUserProperty(turn, CharacterId int, isOwner bool, bet float64) UserProperty {
	return UserProperty{
		Turn:         turn,
		IsLeave:      0,
		IsOwner:      isOwner,
		IsReady:      0,
		PriorityAct:  false,
		IsMyTurn:     false,
		Seat:         0,
		UserLimitNum: 0,
		WinPrice:     0,
		Bet:          bet,
		CharacterId:  CharacterId,
	}
}

// UserExt 用户额外需要
type UserExt struct {
	CountdownTimeNum int    `json:"count_down_time_num"`    //超时出牌次数
	SwingRodNo       string `json:"swing_rod_no,omitempty"` //挥杆编号
	RoomNo           string `json:"room_no,omitempty"`      //房间编号
	//HandarmConfig    HandarmConfig `json:"handarm_config,omitempty"` //骗子酒馆 用户的子弹
}

type HandarmConfig struct {
	Seat      int `json:"seat"`       //子弹的位置
	CurrSeat  int `json:"curr_seat"`  //当前位置
	AllBullet int `json:"all_bullet"` //全部子弹数量
}

//func NewHandarmConfig(seat, currSeat, allBullet int) HandarmConfig {
//	return HandarmConfig{
//		Seat:      seat,
//		CurrSeat:  currSeat,
//		AllBullet: allBullet,
//	}
//}

func (u *UserInfo) AddUserTurn() {
	u.UserProperty.Turn++
}

func (u *UserInfo) GetUserTurn() int {
	return u.UserProperty.Turn
}

func GetOutCardCountDownTimeInt(outCardCountDownTime int64) int64 {
	return helper.LocalTime().Unix() + outCardCountDownTime
}

// SetOutCardCountDown 出牌倒计时（包含随牌时间）
// 1 加载 或者 进入下一轮的时候出发
// 2 托管
func (u *UserInfo) SetOutCardCountDown(outCardCountDown int64) {
	u.UserProperty.OutCardCountDown = outCardCountDown
}

func (u *UserInfo) GetOutCardCountDown() int64 {
	return u.UserProperty.OutCardCountDown
}

func GetReMakeCardDownTimeInt(reMakeCardDownTime int64) int64 {
	return helper.LocalTime().Unix() + reMakeCardDownTime
}

func (u *UserInfo) SetReMakeCardDown(reMakeCardDown int64) {
	u.UserProperty.ReMakeCardDown = reMakeCardDown
}

func (u *UserInfo) GetReMakeCardDown() int64 {
	return u.UserProperty.ReMakeCardDown
}

func (u *UserInfo) SetGameStatus(gameStatus int) {
	u.UserProperty.GameStatus = gameStatus
}

func (u *UserInfo) GetGameStatus() int {
	return u.UserProperty.GameStatus
}

func (rs *UserInfo) AddCountdownTimeNum() {
	rs.UserExt.CountdownTimeNum += 1
}

func (rs *UserInfo) GetCdTimeNum() int {
	return rs.UserExt.CountdownTimeNum
}

func (u *UserInfo) SetUserIsMyTurn(isMyTurn bool) {
	u.UserProperty.IsMyTurn = isMyTurn
}

func (u *UserInfo) GetUserIsMyTurn() bool {
	return u.UserProperty.IsMyTurn
}

func (u *UserInfo) UserIsRobot() bool {
	if u.UserProperty.IsRobot == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserInfo) SetUserIsReady(isReady int) {
	u.UserProperty.IsReady = isReady
}

func (u *UserInfo) GetUserIsReady() int {
	return u.UserProperty.IsReady
}

func (u *UserInfo) SetUserIsOwner(isOwner bool) {
	u.UserProperty.IsOwner = isOwner
}

func (u *UserInfo) GetUserIsIsOwner() bool {
	return u.UserProperty.IsOwner
}

//func (u *UserInfo) SetUserIsKilled(isKilled int) {
//	u.UserProperty.IsKilled = isKilled
//}
//
//func (u *UserInfo) GetUserIsKilled() int {
//	return u.UserProperty.IsKilled
//}
//
//func (u *UserInfo) UserIsKilled() bool {
//	return u.UserProperty.IsKilled == 1
//}
