// Package models 数据模型
package models

import (
	"slot_server/lib/global"
	"time"
)

const (
	heartbeatTimeout = 1 * 60 // 用户心跳超时时间
)

type UserOnline struct {
	AccIp         string `json:"accIp"`         // acc Ip
	AccPort       string `json:"accPort"`       // acc 端口
	AppID         uint32 `json:"appID"`         // appID
	UserID        string `json:"userID"`        // 用户ID
	Nickname      string `json:"nickname"`      // 用户ID
	ClientIp      string `json:"clientIp"`      // 客户端Ip
	ClientPort    string `json:"clientPort"`    // 客户端端口
	LoginTime     uint64 `json:"loginTime"`     // 用户上次登录时间
	HeartbeatTime uint64 `json:"heartbeatTime"` // 用户上次心跳时间
	LogOutTime    uint64 `json:"logOutTime"`    // 用户退出登录的时间
	Qua           string `json:"qua"`           // qua
	DeviceInfo    string `json:"deviceInfo"`    // 设备信息
	IsLogoff      bool   `json:"isLogoff"`      // 是否下线
}

// UserLogin 用户登录
func UserLogin(accIp, accPort string, appID uint32, userID string, addr string,
	loginTime uint64, nickname string) (userOnline *UserOnline) {
	userOnline = &UserOnline{
		AccIp:         accIp,
		AccPort:       accPort,
		AppID:         appID,
		UserID:        userID,
		Nickname:      nickname,
		ClientIp:      addr,
		LoginTime:     loginTime,
		HeartbeatTime: loginTime,
		IsLogoff:      false,
	}
	return
}

// Heartbeat 用户心跳
func (u *UserOnline) Heartbeat(currentTime uint64) {
	u.HeartbeatTime = currentTime
	u.IsLogoff = false
	return
}

// LogOut 用户退出登录
func (u *UserOnline) LogOut() {
	currentTime := uint64(time.Now().Unix())
	u.LogOutTime = currentTime
	u.IsLogoff = true
	return
}

// IsOnline 用户是否在线
func (u *UserOnline) IsOnline() (online bool) {
	if u.IsLogoff {
		return
	}
	currentTime := uint64(time.Now().Unix())
	if u.HeartbeatTime < (currentTime - heartbeatTimeout) {
		global.GVA_LOG.Infof("用户是否在线 心跳超时 u.AppID:%v, u.UserID:%v, u.HeartbeatTime:%v", u.AppID, u.UserID, u.HeartbeatTime)
		return
	}
	if u.IsLogoff {
		global.GVA_LOG.Infof("用户是否在线 用户已经下线 AppID:%v UserID:%v , Nickname:%v", u.AppID, u.UserID, u.Nickname)
		return
	}
	return true
}

// UserIsLocal 用户是否在本台机器上
func (u *UserOnline) UserIsLocal(localIp, localPort string) (result bool) {
	if u.AccIp == localIp && u.AccPort == localPort {
		result = true
		return
	}
	return
}

type UserWeb struct {
	UserID    string `json:"userID"`    // 用户ID
	Nickname  string `json:"nickname"`  // 用户ID
	LoginTime int64  `json:"loginTime"` // 用户上次登录时间
	Token     string `json:"token"`
}
