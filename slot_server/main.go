package main

import (
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"slot_server/lib/core"
	"slot_server/lib/global"
	"slot_server/lib/utils/queue"
	"slot_server/protoc/pbs"
	"slot_server/router"
	"slot_server/servers/grpcserver"
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
	router.NatsRouterInit()

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

	go websocket.MemeRoomManager.Start() //房间服务

	grpcServer := grpc.NewServer()
	pbs.RegisterMemeBattleServiceServer(grpcServer, &grpcserver.MemeBattleService{})
	router.InitRouters()

	addr := global.GVA_VP.GetString("app.mebRpcUrl")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		global.GVA_LOG.Fatalf("failed to listen: %v", err)
	}
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			global.GVA_LOG.Fatalf("failed to serve: %v", err)
		}
	}()

	global.GVA_LOG.Infof("meme_battle 服务启动成功...")
	// 等待终止信号
	sigs := make(chan os.Signal, 1)
	//SIGINT	2	Term	用户发送INTR字符(Ctrl+C)触发
	//SIGTERM	15	Term	结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigs:
			//关闭 gRPC 服务器
			grpcServer.GracefulStop()
			log.Println("Server stopped. Shutting down server...")
			global.GVA_LOG.Infof("Shutting down server...")
			return
		}
	}

}
