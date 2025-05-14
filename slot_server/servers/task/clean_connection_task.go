// Package task 定时任务
package task

import (
	"fmt"
	"runtime/debug"
	"slot_server/servers/src/logic"
	"time"
)

// Init 初始化
func Init() {
	//Timer(3*time.Second, 20*time.Second, cleanConnection, "", nil, nil)

	//埋点 往后推5分钟 ，然后每5分钟执行一次
	//Timer(60*time.Second*5, 60*time.Second*5, GetBetOnList, "", nil, nil)

	//Timer(3*time.Second, time.Minute*300, FishAnalysis, "", nil, nil)
	//加L币 往后推5分钟 ，然后每12个小时执行一次
	//Timer(60*time.Second*5, 24*time.Hour, IncrUserCoin, "", nil, nil)

	//一次性脚本 返回用户的积分 （钓鱼）
	//Timer(time.Second*10, 100*time.Hour, IncrUserCoinFish, "", nil, nil)

	Timer(3*time.Second, 20*time.Second, AnimalPartyGlobal, "", nil, nil)

}

func AnimalPartyGlobal(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" IncrUserCoin stop", r, string(debug.Stack()))
		}
	}()

	logic.AnimalPartyGlobal()

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

func IncrUserCoin(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" IncrUserCoin stop", r, string(debug.Stack()))
		}
	}()

	return
}

func FishAnalysis(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(" 埋点任务停止 stop", r, string(debug.Stack()))
		}
	}()

	return
}
