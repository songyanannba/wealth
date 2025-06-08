package service

import (
	"client/conf"
	"client/protoc/pbs"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type wsClientService struct {
	sync.Mutex
	context chan []byte
	conn    *websocket.Conn
}

func (ws *wsClientService) Send(msgType int32, userId, serviceId string, message proto.Message) {

	messageMarshal, _ := proto.Marshal(message)
	msg := &pbs.NetMessage{
		ServiceId: serviceId,
		Content:   messageMarshal,
	}

	msgMarshal, _ := proto.Marshal(msg)
	ws.context <- msgMarshal
}

var WsClientService = wsClientService{
	context: make(chan []byte, 1024),
}

func (ws *wsClientService) Start() {

	//tk := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk0NjExNjMsInN1YiI6IiIsInVzZXJfaWQiOiJiZDNmZmQ1Mi1kMjc3LTQ1NTAtODZjNy1hN2I1MDIxZWJmNTAifQ.qXMRxNi48pI9CK9_t2aSzB9vFW5QY6TC_UFZ-fxxFuA"

	header := http.Header{}
	header.Add("userId", "syn")
	header.Add("auth", "syn")
	//header.Add("gw-token", tk)

	u := url.URL{
		Scheme: "ws",
		Host:   conf.HOST,
		Path:   conf.PATH,
	}
	s := u.String()

	fmt.Println("url str == ", s)

	conn, _, err := websocket.DefaultDialer.Dial(s, header)
	defer conn.Close()

	ws.conn = conn
	if err != nil {
		fmt.Println("ws dail 服务拨号失败 = ", err)
		return
	}
	//go ws.Read()
	go ws.Write()
	go ws.Read()

	//ws.WSLogin()
	//ws.Test123()
	//time.Sleep(2 * time.Second)

	//获取在线列表
	//ws.TestOnLineUser()

	//获取当局信息
	//ws.TestCurrAPInfo()

	//ws.TestGame()

	//押注TestBetReq
	ws.TestBetReq()

	time.Sleep(1000 * time.Second * 10000)
}

func (ws *wsClientService) Read() {
	fmt.Println("Read for")

	defer ws.conn.Close()

	for {
		var err error
		mType, msg, err := ws.conn.ReadMessage()
		if mType == websocket.TextMessage {
			netMsg := &pbs.NetMessage{}
			err = proto.Unmarshal(msg, netMsg)
			if err == nil {
				//fmt.Println("wsClientService read = ", netMsg.Type)
				//CliHandler.DaYin(netMsg)
				value, ok := CommonService.GetHandlers(netMsg)
				if ok {
					value(netMsg)
				}
			} else {
				fmt.Println("wsClientService read err = ", err)
			}
		}

	}

}

func (ws *wsClientService) Write() {

	for {
		fmt.Println("write for")
		select {
		case context := <-ws.context:
			err := ws.conn.WriteMessage(websocket.TextMessage, context)
			if err != nil {
				fmt.Println("ws write err", err)
			} else {
				mm := &pbs.NetMessage{}
				proto.Unmarshal(context, mm)
				fmt.Println("ws write succ", mm)
			}
		}
	}

}

type RoomConfigReq struct {
	Seq  string     `json:"seq"`
	Cmd  string     `json:"cmd"`
	Data RoomConfig `json:"data"`
}
type RoomConfig struct {
	Uid   string `json:"uid"`
	Token string `json:"token"`
}

func (ws *wsClientService) TestGame() {

	request := &RoomConfigReq{
		Seq: "1",
		Cmd: "getRoomConfig",
		Data: RoomConfig{
			Uid:   "1",
			Token: "1",
		},
	}
	//marshal, _ := proto.Marshal(request)
	marshal, _ := json.Marshal(request)

	//req := &pb.NetMessage{
	//	ServiceId: "syn-service",
	//	UId:       "syn--",
	//	Content:   marshal,
	//	Type:      1,
	//}
	//reqM, _ := proto.Marshal(req)

	ws.context <- marshal
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}

func (ws *wsClientService) WSLogin() {

	tk := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk0NjExNjMsInN1YiI6IiIsInVzZXJfaWQiOiJiZDNmZmQ1Mi1kMjc3LTQ1NTAtODZjNy1hN2I1MDIxZWJmNTAifQ.qXMRxNi48pI9CK9_t2aSzB9vFW5QY6TC_UFZ-fxxFuA"
	req1 := &pbs.Login{
		AppId: 10,
		Token: tk,
	}
	req1M, _ := proto.Marshal(req1)
	reqq := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: "slot_server",
		MsgId:     int32(pbs.ProtocNum_LoginReq),
		Content:   req1M,
	}

	reqM, _ := proto.Marshal(reqq)

	ws.context <- reqM
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}

func (ws *wsClientService) Test123() {

	//request := &pb.GameMessage{
	//	To:   "123",
	//	Do:   "234",
	//	Todo: "4哈哈哈哈6",
	//}
	//marshal, _ := proto.Marshal(request)

	//req := &pb.NetMessage{
	//	ServiceId: "syn-service",
	//	UId:       "syn--",
	//	Content:   marshal,
	//	Type:      1,
	//}

	req1 := &pbs.Login{
		AppId: 10,
	}
	req1M, _ := proto.Marshal(req1)
	reqq := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: "slot_server",
		MsgId:     int32(pbs.ProtocNum_LoginReq),
		Content:   req1M,
	}

	reqM, _ := proto.Marshal(reqq)

	ws.context <- reqM
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}

func (ws *wsClientService) TestCurrAPInfo() {
	time.Sleep(3 * time.Second)
	req1 := &pbs.CurrAPInfoReq{}
	req1M, _ := proto.Marshal(req1)
	reqq := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: "slot_server",
		MsgId:     int32(pbs.ProtocNum_CurrAPInfoReq),
		Content:   req1M,
	}

	reqM, _ := proto.Marshal(reqq)

	ws.context <- reqM
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}

func (ws *wsClientService) TestBetReq() {
	time.Sleep(2 * time.Second)
	req1 := &pbs.UserBetReq{
		Bet:       1,
		GameId:    1,
		BetZoneId: 2,
	}
	req1M, _ := proto.Marshal(req1)
	reqq := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: "slot_server",
		MsgId:     int32(pbs.ProtocNum_betReq),
		Content:   req1M,
	}

	reqM, _ := proto.Marshal(reqq)

	ws.context <- reqM
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}

func (ws *wsClientService) TestOnLineUser() {
	time.Sleep(2 * time.Second)
	req1 := &pbs.OnLineUserListReq{}
	req1M, _ := proto.Marshal(req1)
	reqq := &pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "",
			Token:    "",
			Platform: "",
		},
		AckHead: &pbs.AckHead{
			Uid:     "",
			Code:    0,
			Message: "",
		},
		ServiceId: "slot_server",
		MsgId:     int32(pbs.ProtocNum_OnLineUserListReq),
		Content:   req1M,
	}

	reqM, _ := proto.Marshal(reqq)

	ws.context <- reqM
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}
