package websocket

import (
	"encoding/json"
	"fmt"
	"gateway/protoc/pbs"
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

var WsClientService = wsClientService{
	context: make(chan []byte, 1024),
}

func (ws *wsClientService) Send(msgType int32, userId, serviceId string, message proto.Message) {

	//messageMarshal, _ := proto.Marshal(message)
	req := pbs.Test1Req{UserId: "1"}
	reqMarshal, _ := proto.Marshal(&req)

	netMessage := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "1",
			Token:    "xxx",
			Platform: "a",
		},
		AckHead: &pbs.AckHead{
			Uid:     "1",
			Code:    0,
			Message: "",
		},
		ServiceId: "1",
		MsgId:     1,
		Content:   reqMarshal,
	}
	msgMarshal, _ := proto.Marshal(&netMessage)

	//msgMarshal, _ := proto.Marshal(msg)
	ws.context <- msgMarshal
}

func (ws *wsClientService) Start() {
	header := http.Header{}

	//如果支持protoc协议 创建连接的时候请求头必须加协议类型
	header.Set("protoc_type", "1")

	u := url.URL{
		Scheme: "ws",
		Host:   "127.0.0.1:8099",
		Path:   "gate_way",
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

	ws.Test123()
	//ws.TestGame()

	time.Sleep(1000 * time.Second * 10000)
}

func (ws *wsClientService) Read() {
	fmt.Println("Read for")

	defer ws.conn.Close()

	for {
		var err error
		mType, msg, err := ws.conn.ReadMessage()
		//if mType == websocket.BinaryMessage {
		if mType == websocket.TextMessage {
			netMsg := &pbs.NetMessage{}
			err = proto.Unmarshal(msg, netMsg)
			fmt.Println("wsClientService read = ", netMsg)
			if err == nil {
				test1Ack := &pbs.Test1Ack{}
				err = proto.Unmarshal(msg, test1Ack)
				fmt.Println("wsClientService read = ", netMsg.MsgId, test1Ack)
				//CliHandler.DaYin(netMsg)
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
			//err := ws.conn.WriteMessage(websocket.BinaryMessage, context)
			err := ws.conn.WriteMessage(websocket.TextMessage, context)
			if err != nil {
				fmt.Println("ws write err", err)
			} else {
				mm := &pbs.NetMessage{}
				err = proto.Unmarshal(context, mm)
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

func (ws *wsClientService) Test123() {

	req := pbs.Test1Req{UserId: "1"}
	reqMarshal, _ := proto.Marshal(&req)

	netMessage := pbs.NetMessage{
		ReqHead: &pbs.ReqHead{
			Uid:      "1",
			Token:    "xxx",
			Platform: "a",
		},
		AckHead: &pbs.AckHead{
			Uid:     "1",
			Code:    0,
			Message: "",
		},
		ServiceId: "1",
		MsgId:     1,
		Content:   reqMarshal,
	}
	msgMarshal, _ := proto.Marshal(&netMessage)
	ws.context <- msgMarshal
	/*for {
		time.Sleep(30 * time.Second)
	}*/
}
