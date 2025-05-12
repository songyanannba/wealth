// Package task 定时任务
package task

import (
	"go.uber.org/zap"

	"runtime/debug"

	"time"
)

// ServerInit 服务初始化
func ServerInit() {
	Timer(2*time.Second, 60*time.Second, server, "", serverDefer, "")
}

// server 服务注册
func server(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("服务下线 stop", zap.Any("", r), zap.Any("", string(debug.Stack())))
		}
	}()

	return
}

// serverDefer 服务下线
func serverDefer(param interface{}) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			global.GVA_LOG.Error("服务下线 stop", zap.Any("", r), zap.Any("", string(debug.Stack())))
		}
	}()

	return
}
