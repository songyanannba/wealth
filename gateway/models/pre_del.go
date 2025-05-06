package models

const (
	F_Conn_Succ   = "1" //	 连接成功
	F_result      = "2" //
	F_select_room = "3" //

)

type ComAck struct {
	//Seq      string    `json:"seq"`        // 序列号
	//Cmd      string    `json:"cmd"`        // 序列号
	ProtoNum string    `json:"proto_numb"` //协议号
	UserID   string    `json:"user_id"`    //协议好
	Response *Response `json:"response"`   // 消息体
}

type ConnSucc struct {
	//ComAck
	ProtoNum string `json:"proto_numb"` //协议号
	UserID   string `json:"user_id"`    //协议好
	CurrTime int64  `json:"curr_prize"` // 时间
	Desc     string `json:"desc"`
}

type LogicRespAck struct {
	ServiceToken string `json:"token"` // 验证用户是否登录
	UserID       string `json:"user_id"`
	Nickname     string `json:"nickname,omitempty"` //新版本 不需要传
}

type GetRoomReq struct {
	Uid      string `json:"uid"`
	Bet      int    `json:"bet"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"` //平台的token
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

type UpdateUserScoreResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Time string      `json:"time"`
	Data interface{} `json:"data"`
	Test int         `json:"test"`
}

func NewUpdateUserScoreResp() UpdateUserScoreResp {
	return UpdateUserScoreResp{
		Code: 0,
		Msg:  "",
		Time: "",
		Data: nil,
		Test: 0,
	}
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

type GetCoinBalanceReq struct {
	//UserID string `json:"user_id"`
	//Token  string `json:"token"`
}

type GetCoinBalanceResp struct {
	CoinList []CoinList `json:"list"`
}

type CoinList struct {
	CType int     `json:"c_type"`    // 币种类型 1:国王积分 10:L币
	Num   float64 `json:"price_num"` // 余额
}

type GetUserInfo struct {
	Code int          `json:"code"`
	Test int          `json:"test"`
	Msg  string       `json:"msg"`
	Time string       `json:"time"`
	Data UserInfoData `json:"data"`
}

type UserInfoData struct {
	Id               int    `json:"id"`                //用户ID
	Username         string `json:"username"`          //名称
	Nickname         string `json:"nickname"`          //昵称
	IsAuthentication int    `json:"is_authentication"` //0=没认证,1=认证
}
