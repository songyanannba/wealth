package websocket

import (
	"errors"
	"fmt"
	"gateway/global"
	"gateway/protoc/pbs"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
)

type natsManager struct {
	NatsConn     *nats.Conn
	NatsPubJsMap map[string]jetstream.JetStream
	//PlayersSub  map[string]*nats.Subscription
	Sync *sync.RWMutex
}

var NastManager = natsManager{
	NatsConn:     nil,
	NatsPubJsMap: make(map[string]jetstream.JetStream),
	//PlayersSub:  make(map[string]*nats.Subscription),
	Sync: new(sync.RWMutex),
}

func (n *natsManager) Start() {
	natsUrl := global.GVA_VP.GetString("app.natsUrl")
	connect, err := nats.Connect(
		natsUrl,
		//fmt.Sprintf("nats://127.0.0.1:%d", 4222),
		nats.Timeout(10*time.Second),
		nats.MaxReconnects(3),
		nats.ReconnectWait(5*time.Second),
	)
	if err != nil {
		global.GVA_LOG.Error("nats conn err = ", zap.Error(err))
		return
	}
	n.NatsConn = connect

	//js, err := connect.JetStream()
	//if err != nil {
	//	global.GVA_LOG.Error("nats js err = ", zap.Error(err))
	//	return
	//}
	//// 初始化 JetStream
	//err = memeBattleEnsureStream(js)
	//if err != nil {
	//	log.Fatalf("Failed to initialize JetStream: %v", err)
	//}

	n.CreatNatsPubMTJs(MemeBattle)

	//defer connect.Close()
}

func (n *natsManager) CreatNatsPubMTJs(key string) {
	//js, err := n.NatsConn.JetStream()

	js, err := jetstream.New(n.NatsConn)

	if err != nil {
		n.Close()
		global.GVA_LOG.Error("nats js err = ", zap.Error(err))
		return
	}
	//meme 的nats
	n.NatsPubJsMap[key] = js

	// 初始化 JetStream
	err = memeBattleEnsureStream(js)
	if err != nil {
		log.Fatalf("Failed to initialize JetStream: %v", err)
	}
}

func (n *natsManager) Close() {
	n.NatsConn.Close()
}

func (n *natsManager) servicePub(where string, msg *pbs.NetMessage) {
	marshal, _ := proto.Marshal(msg)
	err := n.NatsConn.Publish(where, marshal)
	if err != nil {
		fmt.Println("natsManager err : ", err)
	}
}

//func (n *natsManager) SendMemeJs(msg *pbs.NetMessage) {
//	err := n.publishMessages(msg)
//	if err != nil {
//		global.GVA_LOG.Error("SendMemeJs", zap.Error(err))
//	}
//}

func (n *natsManager) GetNatsJs(topic string) (jetstream.JetStream, error) {
	n.Sync.RLock()
	defer n.Sync.RUnlock()
	js, ok := n.NatsPubJsMap[topic]
	if !ok {
		global.GVA_LOG.Error("natsManager err : MemeBattle not exist")
		return nil, errors.New("not js")
	}
	return js, nil
}

func (n *natsManager) SendMemeJs(msg *pbs.NetMessage) {
	n.Sync.RLock()
	defer n.Sync.RUnlock()
	js, ok := n.NatsPubJsMap[MemeBattle]
	if !ok {
		global.GVA_LOG.Error("natsManager err : MemeBattle not exist")
		return
	}

	n.NewPublishMessages(msg, js)

}

func (n *natsManager) Send(where string, msg *pbs.NetMessage) {
	n.servicePub(where, msg)
}
