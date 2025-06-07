// Package task 定时任务
package task

import (
	"fmt"
	"gateway/global"
	"gateway/servers/websocket"
	"go.uber.org/zap"
	"runtime/debug"
	"time"
)

// Init 初始化
func Init() {
	Timer(3*time.Second, 20*time.Second, cleanConnection, "", nil, nil)

	//埋点 往后推5分钟 ，然后每5分钟执行一次
	Timer(60*time.Second*5, 60*time.Second*5, GetBetOnList, "", nil, nil)

}

// cleanConnection 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("ClearTimeoutConnections stop", zap.Any("r", r), zap.Any("cleanConnection Stack", string(debug.Stack())))
		}
	}()

	global.GVA_LOG.Infof("定时任务，清理超时连接 %v", param)
	websocket.ClearTimeoutConnections()
	websocket.PrintCurrUser()
	return
}

func GetBetOnList(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" 埋点任务停止 stop", r, string(debug.Stack()))
		}
	}()

	return
}
