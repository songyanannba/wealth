// Package websocket 处理
package websocket

import (
	"github.com/gorilla/websocket"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime    = 10
	heartbeatExpirationTimeDev = 60 * 10
)

// 用户登录
type login struct {
	AppID  uint32
	UserID string
	Client *Client
}

// GetKey 获取 key
func (l *login) GetKey() (key string) {
	//key = GetUserKey(l.AppID, l.UserID)

	return
}

// Client 用户连接
type Client struct {
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

// NewClient 初始化
func NewClient(addr string, socket *websocket.Conn, firstTime uint64, protocType int) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 1000000),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		ProtocType:    protocType,
	}
	return
}

// GetKey 获取 key
func (c *Client) GetKey() (key string) {
	//key = GetUserKey(c.AppID, c.UserID)
	return
}
