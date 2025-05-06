// Package models 数据模型
package models

// Request 通用请求数据格式
type Request struct {
	Seq       string      `json:"seq"`                  // 消息的唯一ID
	Cmd       string      `json:"cmd"`                  // 请求命令字
	ServiceId string      `json:"service_id,omitempty"` // 请求命令字
	Data      interface{} `json:"data,omitempty"`       // 数据 json
	//MsgData []byte      // 数据 json
}

// Login 登录请求数据
type Login struct {
	ServiceToken string `json:"service_token"`    // 验证用户是否登录
	AppID        uint32 `json:"app_id,omitempty"` // 0:挖矿未传 1:未定义 2:钓鱼 3:酒馆故事 4:meme_battle
	Token        string `json:"token,omitempty"`
	Nickname     string `json:"nickname,omitempty"` //新版本 不需要传
	UserID       string `json:"user_id,omitempty"`  //新版本 不需要传
	LoginDouYin
}

type LoginExt struct {
	Platform int
}

type LoginDouYin struct {
	DyCode          string `json:"code"`
	DYAnonymousCode string `json:"anonymous_code"`
}

// HeartBeat 心跳请求数据
type HeartBeat struct {
	UserID string `json:"user_id,omitempty"`
}
