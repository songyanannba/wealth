package main

import (
	"log"
	"os"
	"os/signal"
	"slot_server/lib/core"
	"slot_server/lib/global"
	"slot_server/lib/utils/queue"
	"slot_server/router"
	"slot_server/servers/task"
	"slot_server/servers/websocket"
	"syscall"
)

func InitServer() {
	//初始化
	core.BaseInit()

	global.QueueDataKeyMap = queue.NewQueueDataKeyMap()

	//nats消息中间件
	websocket.NastManager.Start()

	//路由初始化
	router.NatsSlotInit()

	// 定时任务
	task.Init()
	// 服务注册
	task.ServerInit()
}

func CloseServer() {
	core.CloseDB()
	websocket.NastManager.Close()
}

func main() {
	InitServer()
	defer CloseServer()

	go websocket.SlotRoomManager.Start() //房间服务

	// 等待终止信号
	sigs := make(chan os.Signal, 1)
	//SIGINT	2	Term	用户发送INTR字符(Ctrl+C)触发
	//SIGTERM	15	Term	结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			//关闭 gRPC 服务器
			log.Println("Server stopped. Shutting down server...")
			global.GVA_LOG.Infof("Shutting down server...")
			return
		}
	}

}
