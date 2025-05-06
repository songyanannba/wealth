package websocket

import (
	"fmt"
	"gateway/protoc/pbs"
)

type cliHandler struct {
}

var CliHandler = &cliHandler{}

func (ch *cliHandler) Start() {

}

func (ch *cliHandler) DaYin(msg *pbs.NetMessage) {

	fmt.Println("dayin msg == ", msg)

}
