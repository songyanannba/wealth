// Package websocket 处理
package websocket

import (
	"gateway/global"
	"gateway/protoc/pbs"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"runtime/debug"
)

const (
	// 用户连接超时时间
	heartbeatExpirationTime    = 10
	heartbeatExpirationTimeDev = 60 * 60
)

// 用户登录
type login struct {
	AppID  uint32
	UserID string
	Client *Client
}

// GetKey 获取 key
func (l *login) GetKey() (key string) {
	key = GetUserKey(l.AppID, l.UserID)

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
	SentOut       chan *pbs.NetMessage
}

// NewClient 初始化
func NewClient(addr, gwToken string, socket *websocket.Conn, firstTime uint64, protocType int) (client *Client) {
	client = &Client{
		Addr:          addr,
		Socket:        socket,
		Send:          make(chan []byte, 1000000),
		FirstTime:     firstTime,
		HeartbeatTime: firstTime,
		ProtocType:    protocType,
		Token:         gwToken,
	}
	return
}

// GetKey 获取 key
func (c *Client) GetKey() (key string) {
	key = GetUserKey(c.AppID, c.UserID)
	return
}

// 读取客户端数据
func (c *Client) read() {
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("write stop %v", zap.Any("read recover", string(debug.Stack())), zap.Any("r", r))
		}
	}()
	defer func() {
		global.GVA_LOG.Error("读取客户端数据 关闭send", zap.Any("c:", c.UserID), zap.Any("c:", c.Nickname))
		global.GVA_LOG.Error("读取客户端数据 关闭send", zap.Any("c:", *c))
		close(c.Send)
	}()
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			global.GVA_LOG.Error("read 读取客户端数据 错误", zap.Any("c.Addr ", err), zap.Any("message", string(message)))
			return
		}

		// 处理程序
		global.GVA_LOG.Infof("读取客户端数据 处理:%v", string(message))

		//ProcessData(c, message)

		ProcessDataNew(c, message)
	}
}

// 向客户端写数据
func (c *Client) write() {
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("write stop", zap.Any("Stack:", string(debug.Stack())), zap.Any("read r:", r))
		}
	}()
	defer func() {
		clientManager.Unregister <- c
		_ = c.Socket.Close()
		global.GVA_LOG.Infof("Client发送数据 c.UserID %v,c.Addr %v defer %v", c.UserID, c.Addr, c)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				// 发送数据错误 关闭连接
				global.GVA_LOG.Infof("Client 发送数据 关闭连接 c.UserID %v,  c.Addr:%v ;%v %v", c.UserID, c.Addr, "ok", ok)
				return
			}
			//_ = c.Socket.WriteMessage(websocket.BinaryMessage, message)
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {
	if c == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Infof("SendMsg stop: %v %v ", r, string(debug.Stack()))
		}
	}()
	c.Send <- msg
}

// close 关闭客户端连接
func (c *Client) close() {
	close(c.Send)
}

// Login 用户登录
func (c *Client) Login(appID uint32, userID string, loginTime uint64, nickname string, token string) {
	c.AppID = appID
	c.UserID = userID
	c.Nickname = nickname
	c.LoginTime = loginTime
	c.Token = token
	// 登录成功=心跳一次
	c.Heartbeat(loginTime)
}

// Heartbeat 用户心跳
func (c *Client) Heartbeat(currentTime uint64) {
	c.HeartbeatTime = currentTime
	return
}

// IsHeartbeatTimeout 心跳超时
func (c *Client) IsHeartbeatTimeout(currentTime uint64) (timeout bool) {
	IsHeartbeatTimeout := c.HeartbeatTime+heartbeatExpirationTime <= currentTime
	if global.GVA_VP.GetString("app.environment") == "dev" {
		IsHeartbeatTimeout = c.HeartbeatTime+heartbeatExpirationTimeDev <= currentTime
	}
	if IsHeartbeatTimeout {
		timeout = true
	}
	return
}

// IsLogin 是否登录了
func (c *Client) IsLogin() (isLogin bool) {
	// 用户登录了
	if c.UserID != "" {
		isLogin = true
		return
	}
	return
}
