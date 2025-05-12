// Package task 定时任务
package task

import (
	"fmt"

	"time"

	"go.uber.org/zap"
	"runtime/debug"
)

// Init 初始化
func Init() {
	//矿工任务优化 所有任务同步执行
	Timer(time.Minute*30, time.Hour*1, RatTasks, "", nil, nil)

	//加L币 往后推5分钟 ，然后每12个小时执行一次
	Timer(60*time.Second, 24*time.Hour, IncrUserCoin, "", nil, nil)

}

// RatTasks 矿工任务
func RatTasks(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" 挖矿买点 stop", r, string(debug.Stack()))
		}
	}()

	global.GVA_LOG.Infof("矿工任务:开始 RatTasks newTime:%v", helper.TimeIntToStr(time.Now().Unix()))

	global.GVA_LOG.Infof("矿工任务:结束 RatTasks %v", helper.TimeIntToStr(time.Now().Unix()))
	return
}

// RatGetOrePerHour 挖矿逻辑
func RatGetOrePerHour(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" 挖矿买点 stop", r, string(debug.Stack()))
		}
	}()

	global.GVA_LOG.Infof("RatGetOrePerHour")

	return
}

func IncrUserCoin(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" IncrUserCoin stop", r, string(debug.Stack()))
		}
	}()

	//logic.IncrUserCoin()
	return
}

// cleanConnection 清理超时连接
func cleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("ClearTimeoutConnections stop", zap.Any("r", r), zap.Any("cleanConnection Stack", string(debug.Stack())))
		}
	}()

	//global.GVA_LOG.Infof("定时任务，清理超时连接 %v", param)
	//websocket.ClearTimeoutConnections()
	return
}
