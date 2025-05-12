package models

// 通用协议号
const (
	MsgId1           = iota //
	MsgIdCanSendCard        //酒馆故事 开始发牌
	MsgId3

	MsgIdLeaveRoom = 4 //离开房间
)

type ComMsg struct {
	MsgId string //消息ID
	Data  []byte //消息内容
}

type TurnSendCard struct {
	Turn int
}

// 顶号重新登陆
const (
	RepLogin = "-1"
)

type RepLoginMsg struct {
	ProtoNum  string `json:"proto_numb"`
	Timestamp int64  `json:"timestamp"`
}
