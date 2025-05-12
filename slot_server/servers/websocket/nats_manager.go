package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"log"
	"slot_server/lib/config"
	"slot_server/lib/global"
	"slot_server/protoc/pbs"
	"sync"
	"time"
)

type natsManager struct {
	NatsConn  *nats.Conn
	NatsMebJs jetstream.JetStream
	//PlayersSub map[string]*nats.Subscription
	Sync sync.Mutex
}

var NastManager = natsManager{
	NatsConn:  nil,
	NatsMebJs: nil,
	//PlayersSub: make(map[string]*nats.Subscription),
	Sync: sync.Mutex{},
}

// 拉模式消费者服务
func (n *natsManager) consumer() {
	time.Sleep(2 * time.Second)
	js := n.NatsMebJs

	cons, err := js.CreateConsumer(context.Background(), config.MemeBattle, jetstream.ConsumerConfig{
		Durable:       "MemeBattleTopic_Consumer",
		FilterSubject: config.MemeBattleTopic,
		AckPolicy:     jetstream.AckExplicitPolicy,
		MaxDeliver:    3,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		global.GVA_LOG.Error("consumer 创建消费者失败: ", zap.Error(err))
		return
	}
	for {
		msgs, err := cons.Fetch(100, jetstream.FetchMaxWait(5*time.Second))

		if errors.Is(err, context.DeadlineExceeded) {
			continue
		} else if err != nil {
			log.Printf("拉取失败: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		//global.GVA_LOG.Infof("Processing consumer: %v", len(msgs.Messages()))

		// 处理消息
		for msg := range msgs.Messages() {
			global.GVA_LOG.Infof("Processing message: %s", string(msg.Data()))
			// 发送处理结果给生产者

			response := []byte("Processed: " + string(msg.Data()))
			global.GVA_LOG.Infof("Sending response: %s", string(response))

			// 发送处理结果给生产者
			req := &pbs.NetMessage{}
			err := proto.Unmarshal(msg.Data(), req)
			if err != nil {
				global.GVA_LOG.Error("Error unmarshalling message: %v", zap.Error(err))
				continue
			}
			if req.MsgId <= 0 {
				global.GVA_LOG.Infof("consumer Skipping message because msgId is zero,%v", req)
				continue
			}

			ProcessData(req)
			// 确认消息已处理
			if err := msg.Ack(); err != nil {
				global.GVA_LOG.Infof("Failed to acknowledge message: %v", err)
			} else {
				global.GVA_LOG.Infof("Acknowledged: %v", string(msg.Data()))
			}

		}
	}
}

func (n *natsManager) Producer(msgData []byte) {
	_, err := n.NatsMebJs.Publish(context.Background(), config.MemeBattleTopicResp, msgData)
	if err != nil {
		global.GVA_LOG.Infof("Failed to send response: %v", err)
	}
}

func memeBattleEnsureStream(js jetstream.JetStream) error {
	_, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:      config.MemeBattle,
		Subjects:  []string{config.MemeBattleTopic, config.MemeBattleTopicResp},
		Retention: jetstream.WorkQueuePolicy,
		Storage:   jetstream.FileStorage, // 存储类型（文件存储）
		MaxAge:    2 * time.Hour,         // 消息保留时间
		Replicas:  1,
	})

	if err != nil && !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
		return fmt.Errorf("创建流失败: %w", err)
	}
	return nil
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
		global.GVA_LOG.Error("Start nats conn err = ", zap.Error(err))
		return
	}

	n.NatsConn = connect

	n.CreatNatsPubMTJs()

	// 启动生产者和消费者
	go n.consumer()

	global.GVA_LOG.Info("Start StartConsumer")

}

func (n *natsManager) CreatNatsPubMTJs() {
	//js, err := n.NatsConn.JetStream()

	js, err := jetstream.New(n.NatsConn)

	if err != nil {
		n.Close()
		global.GVA_LOG.Error("nats js err = ", zap.Error(err))
		return
	}
	//meme 的nats
	n.NatsMebJs = js

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
