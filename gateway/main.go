package main

import (
	"gateway/core"
	"gateway/global"
	"gateway/routers"
	"gateway/servers/task"
	"gateway/servers/websocket"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitServer() {
	//初始化
	core.BaseInit()

	//nats消息中间件
	websocket.NastManager.Start()

	// 定时任务
	task.Init()

	// 服务注册
	task.ServerInit()
}

func CloseServer() {
	core.CloseDB()

	//nats
	//websocket.NastManager.Close()
	global.GVA_LOG.Infof("服务结束 CloseServer")
}

func main() {
	InitServer()
	defer CloseServer()

	//gin服务
	router := gin.Default()

	// Gin 框架示例（其他框架逻辑类似）
	router.Use(cors.Default())

	// 初始化路由
	routers.Init(router)
	routers.WebsocketInit()

	go websocket.StartWebSocket()

	//go grpcserver.Init()
	global.GVA_LOG.Infof("gate_way 服务启动成功...")
	httpPort := global.GVA_VP.GetString("app.httpPort")
	_ = http.ListenAndServe(":"+httpPort, router)
}
