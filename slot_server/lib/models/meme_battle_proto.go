package models

import (
	"time"
)

// MatchSuccResp 快速匹配成功协议
type MatchSuccResp struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	RoomCom
	RoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

// DealCardsMsg 发牌
type DealCardsMsg struct {
	ProtoNum  string  `json:"proto_numb"`
	Timestamp int64   `json:"timestamp"`
	UserId    string  `json:"user_id"`         //被质疑者id
	RoomNo    string  `json:"room_no"`         //房间编号
	Turn      int     `json:"turn"`            //第几小轮
	Cards     []*Card `json:"cards,omitempty"` //被质疑者的牌
}

// ReadyMsg 准备就绪广播
type ReadyMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"user_id"` //谁就绪了
	RoomNo    string `json:"room_no"` //房间编号
}

// LeaveRoomMsg 离开房间广播
type LeaveRoomMsg struct {
	ProtoNum     string `json:"proto_numb"`
	Timestamp    int64  `json:"timestamp"`
	UserId       string `json:"user_id"`        //谁就绪了
	RoomNo       string `json:"room_no"`        //房间编号
	IsOwnerLeave bool   `json:"is_owner_leave"` //
	NewOwner     string `json:"new_owner"`
}

type UserStateMsg struct {
	ProtoNum   string    `json:"proto_numb"`
	Timestamp  int64     `json:"timestamp"`
	UserId     string    `json:"user_id"`               //谁就绪了
	RoomNo     string    `json:"room_no"`               //房间编号
	IsContinue bool      `json:"is_continue"`           //是否5连抽卡
	RoomDetail *RoomItem `json:"room_detail,omitempty"` //用户所在房间
}

// KickRoomMsg 踢人广播
type KickRoomMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"user_id"` //谁就绪了
	RoomNo    string `json:"room_no"` //房间编号
}

type InviteFriendMsg struct {
	ProtoNum  string       `json:"proto_numb"`
	Timestamp int64        `json:"timestamp"`
	UserId    string       `json:"user_id"`
	RoomNo    string       `json:"room_no"`              //房间编号
	OwnerInfo MemeRoomUser `json:"owner_info,omitempty"` //房主信息
}

// StartPlayMsg 开始游戏广播
type StartPlayMsg struct {
	ProtoNum         string         `json:"proto_numb"`
	Timestamp        int64          `json:"timestamp"`
	RoomNo           string         `json:"room_no"`                  //房间编号
	MemeRoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

// EntryLikePageMsg 进入点赞页面
type EntryLikePageMsg struct {
	ProtoNum  string  `json:"proto_numb"`
	Timestamp int64   `json:"timestamp"`
	RoomNo    string  `json:"room_no"` //房间编号
	OutCards  []*Card `json:"out_cards,omitempty"`
}

type CalculateRankMsg struct {
	ProtoNum       string            `json:"proto_numb"`
	Timestamp      int64             `json:"timestamp"`
	RoomNo         string            `json:"room_no"` //房间编号
	LikeDetailList []*UserLikeDetail `json:"like_detail_list,omitempty"`
}

type LoadMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	RoomCom
	RoomUserList   []MemeRoomUser  `json:"room_user_list,omitempty"`   //
	OtherUserCards []UserCartState `json:"other_user_cards,omitempty"` //其他用户的情况
	OutCards       []*Card         `json:"out_cards,omitempty"`        //出牌
	LikeCards      []*LikeCard     `json:"like_cards,omitempty"`       //被点赞的牌
}

// JoinRoomMsg 加入房间
type JoinRoomMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	RoomCom
	RoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

type CreateRoomMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	RoomCom
	RoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

func NewCreateRoomMsg(protoNum string, roomCom RoomCom, roomUserList []MemeRoomUser) CreateRoomMsg {
	return CreateRoomMsg{
		ProtoNum:     protoNum,
		Timestamp:    time.Now().Unix(),
		RoomCom:      roomCom,
		RoomUserList: roomUserList,
	}
}

type IssueMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	Issue     *Issue `json:"issue"`             //
	UserId    string `json:"user_id,omitempty"` //下一个 出牌用户
}

// OperateCardsMsg 操作牌
type OperateCardsMsg struct {
	ProtoNum   string  `json:"proto_numb"`
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"user_id"`        //出牌用户
	OutCardNum int     `json:"out_card_num"`   //出牌数量
	CardNum    int     `json:"card_num"`       //还有几张牌
	Pitch      float32 `json:"pitch"`          //端侧用
	Yaw        float32 `json:"yaw"`            //端侧用
	Looking    bool    `json:"looking"`        //端侧用 看牌传
	EmojiId    string  `json:"emoji_id"`       //端侧用 表情传
	Card       []*Card `json:"card,omitempty"` //
}

type LikeCardsMsg struct {
	ProtoNum   string  `json:"proto_numb"`
	Timestamp  int64   `json:"timestamp"`
	LikeUserId string  `json:"like_user_id"`   //被点赞用户
	UserId     string  `json:"user_id"`        //点赞用户
	Card       []*Card `json:"card,omitempty"` //
}

type GameOverMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	RoomCom
	RoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

type NextOutCardUserMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
	UserId    string `json:"user_id"` //出牌用户
}

type HandListCard struct {
	CardId int    `json:"card_id"`
	Name   string `json:"name"`   // 名字
	Suffix string `json:"suffix"` //后缀类型
	IsOwn  bool   `json:"is_own"`
	Level  int    `json:"level"`
}

type UnpackCardMsg struct {
	ProtoNum     string          `json:"proto_numb"`
	Timestamp    int64           `json:"timestamp"`
	HandListCard []*HandListCard `json:"hand_list_card,omitempty"`
}

type CardVersionListMsg struct {
	ProtoNum        string             `json:"proto_numb"`
	Timestamp       int64              `json:"timestamp"`
	CardVersionList []*CardVersionList `json:"card_version_list,omitempty"`
}

type CardVersionList struct {
	Version int `json:"version"`
}
