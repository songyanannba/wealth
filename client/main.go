package main

import (
	"client/service"
	"fmt"
	"os"
	"os/signal"
)

func main() {

	controlC := make(chan os.Signal, 1)
	signal.Notify(controlC)

	service.CommonService.Start()

	service.CliHandler.Start()

	service.WsClientService.Start()

	fmt.Println("启动成功。。。")
	for {
		select {
		case <-controlC:
			return
		}
	}

}
