package models

import (
	"slot_server/lib/models/table"
)

// Handarm 手枪

type Handarm struct {
	Bullets  []Bullet `json:"bullets"`   //初始只有1个子弹
	CurrSeat int      `json:"curr_seat"` //初始位置0
}

// Bullet 子弹
type Bullet struct {
	Id     int  `json:"id"`
	IsFill bool `json:"is_fill"` //是否填充子弹
}

type Card struct {
	CardId  int     `json:"card_id"`
	Type    int     `json:"card_type"`  //
	Level   int     `json:"card_level"` //等级 1=流辉级 2=幻彩级 3=璀璨
	Name    string  `json:"name"`       // 名字
	Point   int     `json:"point"`      // 点数
	Express int     `json:"express"`    //
	Suffix  string  `json:"suffix"`     //后缀类型
	ImgUrl  string  `json:"img_url"`
	UserID  string  `json:"user_id,omitempty"`
	AddRate float64 `json:"add_rate,omitempty"` //加成
}

type LikeCard struct {
	CardId     int     `json:"card_id"`                //牌ID
	Level      int     `json:"card_level"`             //等级 1=流辉级 2=幻彩级 3=璀璨
	AddRate    float64 `json:"add_rate,omitempty"`     //加成
	LikeNum    int     `json:"like_num,omitempty"`     //点赞次数
	LikeUserId string  `json:"like_user_id,omitempty"` //被点赞用户
	UserID     string  `json:"user_id,omitempty"`      //点赞用户
}

type UserLikeDetail struct {
	UserID      string  `json:"user_id,omitempty"`
	Nickname    string  `json:"nickname"`
	HeadPhoto   string  `json:"head_photo"`     //头像
	OnGoLinkNum int     `json:"on_go_like_num"` //连续点赞次数
	Integral    float64 `json:"integral"`       //获取积分（得分）
	Experience  float64 `json:"experience"`     //经验值
	MCoin       float64 `json:"m_coin"`         //币
}

type Issue struct {
	IssueId int    `json:"issue_id"`
	Level   int    `json:"level"`
	Class   int    `json:"class"`
	Desc    string `json:"desc"` // 问题描述
}

type CreateRoomInterior struct {
	UserID   string                 `json:"user_id"`
	RoomInfo *table.AnimalPartyRoom `json:"room_info"`
}

// 房间列表

type RoomListReq struct {
	UserID    string `json:"user_id"`
	RoomType  int    `json:"room_type"`  // 房间类型 1:纸牌 2 骰子
	RoomLevel int    `json:"room_level"` //房间 等级 0:初级 1:中级 2:高级 10全部
	RoomId    int    `json:"room_id"`    //房间ID
}

type RoomListResp struct {
	IsHaveNextPage int        `json:"is_have_next_page" example:"1"` //0 没有下一页 ； 1 有下一页
	RoomList       []RoomItem `json:"room_list"`                     //
}

type RoomItem struct {
	RoomCom
	RoomUserList []MemeRoomUser `json:"room_user_list,omitempty"` //
}

type RoomCom struct {
	RoomId       int    `json:"room_Id"`                //房间编号
	RoomNo       string `json:"room_no"`                //房间编号
	UserId       string `json:"user_id"`                //
	Turn         int    `json:"turn"`                   //在第几轮
	RoomName     string `json:"room_name"`              //房间 名字
	Status       int8   `json:"status"`                 //房间状态: 1=开放中,2=已满员,3=已解散,4=进行中,5=已结束 6=异常房间 7=服务字段清理残存房间 8=清理匹配成功用户之前的房间
	UserNumLimit int    `json:"user_num_limit"`         //用户人数限制 2人场 3 人场 4人场
	RoomType     int    `json:"room_type"`              //房间 类型 1
	RoomLevel    int    `json:"room_level"`             //房间 等级 0:初级 1:中级 2:高级 10全部
	CurrIssue    *Issue `json:"curr_issue,omitempty"`   //当前轮 问题
	NextRoomNo   string `json:"next_room_no,omitempty"` //下一个房间编号
	TimeDown     int64  `json:"time_down"`              //游戏状态倒计时
	//0=游戏未开始
	//1=游戏开始但是没有加载完成
	//2=用户随牌阶段
	//3=用户出牌阶段
	//4=用户点赞阶段
	//5=点赞界面 等待结算或者进入下一轮
	GameStatus int `json:"game_status"` //在游戏中的状态
}

func NewRoomCom(roomNo, userId, roomName string, roomId, turn, userNumLimit, roomType, roomLevel int, status int8) RoomCom {
	return RoomCom{
		RoomId:       roomId,
		RoomNo:       roomNo,
		UserId:       userId,
		Turn:         turn,
		RoomName:     roomName,
		Status:       status,
		UserNumLimit: userNumLimit,
		RoomType:     roomType,
		RoomLevel:    roomLevel,
	}
}

// MemeRoomUser 房间用户结构
type MemeRoomUser struct {
	UserID       string        `json:"user_id"`
	Nickname     string        `json:"nickname"`
	Turn         int           `json:"turn"`                 //在第几轮
	IsRobot      int           `json:"is_robot"`             //是否机器人 0:真实用户 1:机器人
	IsLeave      int           `json:"is_leave"`             // 1 离开 （见好就收）
	IsOwner      bool          `json:"is_owner"`             //是否房主
	IsReady      int           `json:"is_ready"`             //是否就绪 0 没有就绪 1就绪
	Seat         int           `json:"seat"`                 //位子、顺序 从0开始
	UserLimitNum int           `json:"user_limit_num"`       //房间人数限制
	UserCards    UserCartState `json:"user_cards,omitempty"` //用户自己的牌
	WinPrice     float64       `json:"win_price"`            // 最终赢钱
	Bet          float64       `json:"bet"`                  // 押注
	//TimeDown     int64         `json:"time_down"`            //游戏状态倒计时
	//0=游戏未开始
	//1=游戏开始但是没有加载完成
	//2=用户随牌阶段
	//3=用户出牌阶段
	//4=用户点赞阶段
	//5=点赞界面 等待结算或者进入下一轮
	//GameStatus int `json:"game_status"` //在游戏中的状态
}

type UserCartState struct {
	UserID     string  `json:"user_id"`
	OutCardNum int     `json:"out_card_num"`        //出牌数量
	CardNum    int     `json:"card_num"`            //当前还有 几张牌
	Card       []*Card `json:"card_list,omitempty"` //手里的牌
}

type HandarmState struct {
	ShootNum  int `json:"shoot_num"`  //射击次数
	AllBullet int `json:"all_bullet"` //全部子弹数量
}

// JoinRoomReq 加入房间
type JoinRoomReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

//准备就绪

type ReadyReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

type ReadyResp struct {
}

// LeaveRoomReq 离开房间
type LeaveRoomReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

type LeaveRoomResp struct {
}

// KickRoomReq 踢人
type KickRoomReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

type KickRoomInter struct {
	UserID  string `json:"user_id"`
	OwnerId string `json:"owner_id"` //房主
	RoomNo  string `json:"room_no"`  //房间编号
}

type KickRoomResp struct {
}

type InviteFriendReq struct {
	UserID string `json:"user_id"` //被邀请人的用户ID
	RoomNo string `json:"room_no"` //房间编号
}

// InviteFriendInter  邀请
type InviteFriendInter struct {
	InviteUserID string `json:"invite_user_id"`
	OwnerId      string `json:"owner_id"` //房主
	RoomNo       string `json:"room_no"`  //房间编号
}

type FraudCardResp struct {
	Card *Card `json:"card"`
}

// 看牌、选牌、

// OperateCardReq 出牌
type OperateCardReq struct {
	UserID  string  `json:"user_id"`
	EmojiId string  `json:"emoji_id"`
	RoomNo  string  `json:"room_no"`  //房间编号
	OpeType int     `json:"ope_type"` //0:看牌 1:出牌 2:摇头 3:表情
	Pitch   float32 `json:"pitch"`
	Yaw     float32 `json:"yaw"`
	Looking bool    `json:"looking"` //看牌传
	Card    []*Card `json:"card_list"`
}

// OperateCardResp 返回剩余的牌
type OperateCardResp struct {
	Card []*Card `json:"card_list,omitempty"`
}

//下家质疑

type IsDoubtReq struct {
	UserID  string `json:"user_id"`
	RoomNo  string `json:"room_no"`  //房间编号
	IsDoubt int    `json:"is_doubt"` //0:质疑 1:不质疑
}

// 服务广播双方牌

type IsDoubtResp struct {
}

//结算

type GameOverReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

type GameOverResp struct {
	RoomNo       string         `json:"room_no"`        //房间编号
	RoomUserList []MemeRoomUser `json:"room_user_list"` //
}

type CharacterListReq struct {
	UserID string `json:"user_id"`
}

type CharacterListResp struct {
	CharacterList []CharacterList `json:"character_list"` //
}

type CharacterList struct {
	CharacterId int `json:"character_Id"` //角色ID
	//CharacterName string `json:"character_name"` //角色名字
}

type RoomConfigReq struct {
}

type RoomConfigResp struct {
	RoomConfig []RoomConfig `json:"room_config"` //
}

type RoomConfig struct {
	Bet            float64 `json:"bet"`
	AdmissionPrice float64 `json:"admission_price" `
	RoomLevel      int8    `json:"room_level"`
}

type ChooseCharacterReq struct {
	UserID      string `json:"user_id"`
	CharacterId int    `json:"character_Id"` //角色ID
}

type ChooseCharacterResp struct {
	CharacterId   int    `json:"character_Id"`   //角色ID
	CharacterName string `json:"character_name"` //角色名字
}

type UserStateReq struct {
	UserID string `json:"user_id"`
	RoomNo string `json:"room_no"` //房间编号
}

type UserStateResp struct {
	CharacterId int       `json:"character_Id"`          //角色ID
	RoomDetail  *RoomItem `json:"room_detail,omitempty"` //用户所在房间
}

type CardLikeReq struct {
	LikeUserID string  `json:"like_user_id"` //被点赞的用户ID
	UserID     string  `json:"user_id"`      //用户ID
	RoomNo     string  `json:"room_no"`      //房间编号
	Card       []*Card `json:"card_list"`    //被点赞的牌
}

type CardLikeResp struct {
}

type UserFriendResp struct {
	IsHaveNextPage bool          `json:"is_have_next_page"` //是否还有下一页
	UserFriend     []*UserFriend `json:"user_friend,omitempty"`
}

type UserFriend struct {
	FriendUserId string `json:"friend_user_id"`
	Nickname     string `json:"nickname"`
	FriendId     int    `json:"friend_id"`
}

type AuditUserResp struct {
	IsHaveNextPage bool         `json:"is_have_next_page"` //是否还有下一页
	AuditUser      []*AuditUser `json:"audit_user,omitempty"`
}

type AuditUser struct {
	ApplicationUser string `json:"application_user"`
	Nickname        string `json:"nickname"`
	AuditId         int    `json:"audit_id"`
}
