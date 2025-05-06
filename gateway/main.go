package main

import (
	"gateway/core"
	"gateway/global"
	"gateway/routers"
	"gateway/servers/grpcclient"
	"gateway/servers/task"
	"gateway/servers/websocket"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func InitServer() {
	//初始化
	core.BaseInit()

	//nats消息中间件
	websocket.NastManager.Start()

	// 初始化 gRPC 客户端
	oreRpcUrl := global.GVA_VP.GetString("app.slotRpcUrl")
	grpcclient.InitMebClient(oreRpcUrl)

	// 定时任务
	task.Init()

	// 服务注册
	task.ServerInit()
}

func CloseServer() {
	core.CloseDB()
	//grpcclient.Close()   // 确保在程序结束时关闭连接

	websocket.NastManager.Close()
	global.GVA_LOG.Infof("服务结束 CloseServer")
}

func main() {
	InitServer()
	defer CloseServer()

	//gin服务
	router := gin.Default()

	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()

	time.Sleep(10 * time.Microsecond)
	go websocket.StartWebSocket()

	//go grpcserver.Init()
	global.GVA_LOG.Infof("gate_way 服务启动成功...")
	httpPort := global.GVA_VP.GetString("app.httpPort")
	_ = http.ListenAndServe(":"+httpPort, router)
}
