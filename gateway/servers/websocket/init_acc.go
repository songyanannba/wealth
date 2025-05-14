// Package websocket 处理
package websocket

import (
	"fmt"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	defaultAppID = 101 // 默认平台ID
)

var (
	clientManager = NewClientManager()                                    // 管理者
	appIDs        = []uint32{defaultAppID, 102, 103, 104, 1, 5, 6, 7, 10} // 全部的平台
	serverIp      string
	serverPort    string
)

// GetAppIDs 所有平台
func GetAppIDs() []uint32 {
	return appIDs
}

// GetServer 获取服务器
func GetServer() (server *models.Server) {
	server = models.NewServer(serverIp, serverPort)
	return
}

// IsLocal 判断是否为本机
func IsLocal(server *models.Server) (isLocal bool) {
	if server.Ip == serverIp && server.Port == serverPort {
		isLocal = true
	}
	return
}

// InAppIDs in app
func InAppIDs(appID uint32) (inAppID bool) {
	for _, value := range appIDs {
		if value == appID {
			inAppID = true
			return
		}
	}
	return
}

// GetDefaultAppID 获取默认 appID
func GetDefaultAppID() (appID uint32) {
	appID = defaultAppID
	return
}

// StartWebSocket 启动程序
func StartWebSocket() {
	serverIp = helper.GetServerIp()

	webSocketPort := global.GVA_VP.GetString("app.webSocketPort")
	socketPort := global.GVA_VP.GetString("app.webSocketPort")
	serverPort = socketPort

	http.HandleFunc("/gate_way", wsPage)

	// 添加处理程序
	go clientManager.start()
	//go TavernRoomManager.Start() //房间服务

	//关闭nats
	go NastManager.slotServiceSubConsumer()

	global.GVA_LOG.Infof("WebSocket 启动程序成功 %v:%v ", serverIp, serverPort)
	_ = http.ListenAndServe(":"+webSocketPort, nil)
}

func wsPage(w http.ResponseWriter, req *http.Request) {

	gwToken1 := req.Header.Get("gw_token")
	fmt.Println(gwToken1)

	// 升级协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		global.GVA_LOG.Infof("升级协议 ua:%v referer:%v", r.Header["User-Agent"], r.Header["Referer"])
		//xCustomHeader := r.Header["X-Custom-Header"]
		//environmentVal := global.GVA_VP.GetString("app.environment")
		//if environmentVal == "pro" {
		//	if len(xCustomHeader) <= 0 || xCustomHeader[0] != "game" {
		//		global.GVA_LOG.Error("请求头没有负载均衡参数wsPage")
		//		return false
		//	}
		//}

		return true
	}}).Upgrade(w, req, nil)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	global.GVA_LOG.Infof("webSocket 建立连接:%v", conn.RemoteAddr().String())

	var protocType int
	//xProtocType := req.Header.Get("protoc_type")
	//if len(xProtocType) != 0 && xProtocType == "1" {
	//	protocType = 1
	//}
	protocType = 1
	gwToken := req.Header.Get("gw-token")
	currentTime := uint64(time.Now().Unix())
	client := NewClient(conn.RemoteAddr().String(), gwToken, conn, currentTime, protocType)

	go client.read()
	go client.write()

	// 用户连接事件
	clientManager.Register <- client
}
