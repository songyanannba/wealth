// Package websocket 处理
package websocket

const (
	defaultAppID = 101 // 默认平台ID
)

var (
	//clientManager = NewClientManager()                                   // 管理者
	appIDs     = []uint32{defaultAppID, 102, 103, 104, 2, 1, 3, 4, 5} // 全部的平台
	serverIp   string
	serverPort string
)
