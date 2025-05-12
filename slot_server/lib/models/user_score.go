package models

type GetRoomReq struct {
	Uid      string `json:"uid"`
	Bet      int    `json:"bet"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"` //平台的token
}

type GetRoomResp struct {
	RoomId   int            `json:"room_id"`   //房间ID
	RoomNo   string         `json:"room_no"`   //房间编号
	IsOpen   int            `json:"is_open"`   //0 未知  1 匹配中 2 满员游戏开始（进行中） 3 关闭 结束
	UserInfo []UserInfoResp `json:"user_info"` //缓冲区用户
}

type UserInfoResp struct {
	Uid          string `json:"uid"`
	Nickname     string `json:"nickname"`
	RoomNo       string `json:"room_no"`        //房间编号
	RoomConfigNo string `json:"room_config_no"` //房间配置编号
	Turn         int    `json:"turn"`           //在第几轮
	IsKilled     int    `json:"is_killed"`      // 1 被杀
	IsLeave      int    `json:"is_leave"`       // 1 离开 （见好就收） //todo
}

type GetUserScoreResp struct {
	Code int              `json:"code"`
	Msg  string           `json:"msg"`
	Time string           `json:"time"`
	Data GetUserScoreData `json:"data"`
	Test int              `json:"test"`
}

type GetUserScoreData struct {
	Score            float64 `json:"score"`
	RegisterTime     int     `json:"register_time"`
	PrevLoginTime    int     `json:"prev_login_time"`
	IsAuthentication int     `json:"is_authentication"` //0=没认证,1=认证
}

type UpdateUserScore struct {
	Uid    string `json:"uid"`     //
	RoomNo string `json:"room_no"` //房间编号
	Bet    string `json:"bet"`     //押注
	Settle string `json:"settle"`  //吃鸡 主动退出
	//游戏过程:
	//|| (1-4:大逃杀) 1=投注,2=主动退出,3=吃鸡 4=主动增加积分
	//||（11-20：钓鱼）
	//|| (21-30：挖矿）21:金矿销毁补偿 （矿石兑换金币）；22:金矿购买 减自己的金币（转增）; 23:金矿购买 加售卖人的金币（转增）; 24.挖金矿矿工购买（购买矿工） 25:金矿游戏补偿积分 26:重置年龄
	//|| (31-40：骗子酒店）
	Process string `json:"process"`
}

type UpdateUserScoreResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Time string      `json:"time"`
	Data interface{} `json:"data"`
	Test int         `json:"test"`
}

type AddUserScoreResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time string `json:"time"`
}

type AddUserScore struct {
	Uid     string `json:"uid"`     //
	RoomNo  string `json:"room_no"` //房间编号
	Num     string `json:"Num"`     //增加积分数量
	Process string `json:"process"` //游戏过程: 1=投注,2=主动退出,3=吃鸡 //4 主动增加积分
}

type JsonData struct {
	Uid       string `json:"uid"`
	RoomNo    string `json:"room_No"`
	Bet       string `json:"bet"`
	Settle    string `json:"settle"`
	Process   string `json:"process"`
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
}

type CloseRoom struct {
	SwingRodNo string `json:"swing_rod_no"` //房间编号
	IsStop     bool   `json:"is_stop"`
	RoomNo     string `json:"room_no"` //房间编号
}
