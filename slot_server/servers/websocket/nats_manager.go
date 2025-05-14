package websocket

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
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
	NatsMebJs nats.JetStreamContext
	Sync      sync.Mutex
}

var NastManager = natsManager{
	NatsConn:  nil,
	NatsMebJs: nil,
	Sync:      sync.Mutex{},
}

// 拉模式消费者服务
func (n *natsManager) consumer() {
	time.Sleep(2 * time.Second)
	js := n.NatsMebJs

	//cons, err := js.CreateConsumer(context.Background(), config.AnimalParty, jetstream.ConsumerConfig{
	//	Durable:       "AnimalPartyTopic_Consumer",
	//	FilterSubject: config.AnimalPartyTopic,
	//	AckPolicy:     jetstream.AckExplicitPolicy,
	//	MaxDeliver:    3,
	//	AckWait:       30 * time.Second,
	//})

	cons, err := js.PullSubscribe(config.AnimalPartyTopic, config.AnimalParty, nats.MaxRequestExpires(10*time.Second), nats.AckWait(30*time.Second), nats.MaxDeliver(3))
	if err != nil {
		global.GVA_LOG.Error("memeBattleServiceSubConsumer 创建消费者失败: ", zap.Error(err))
		return
	}

	if err != nil {
		global.GVA_LOG.Error("consumer 创建消费者失败: ", zap.Error(err))
		return
	}
	for {
		msgs, err := cons.Fetch(100, nats.MaxWait(5*time.Second))

		if errors.Is(err, context.DeadlineExceeded) {
			continue
		} else if err != nil {
			//log.Printf("拉取失败: %v", err)
			global.GVA_LOG.Infof("consumer 拉取失败:%v ", zap.Error(err))
			continue
		}
		//global.GVA_LOG.Infof("Processing consumer: %v", len(msgs.Messages()))

		// 处理消息
		for _, msg := range msgs {
			global.GVA_LOG.Infof("Processing message: %s", string(msg.Data))
			// 发送处理结果给生产者

			response := []byte("Processed: " + string(msg.Data))
			global.GVA_LOG.Infof("Sending response: %s", string(response))

			// 发送处理结果给生产者
			req := &pbs.NetMessage{}
			err := proto.Unmarshal(msg.Data, req)
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
				global.GVA_LOG.Infof("Acknowledged: %v", string(msg.Data))
			}

		}
	}
}

func (n *natsManager) Producer(msgData []byte) {
	_, err := n.NatsMebJs.Publish(config.AnimalPartyTopicResp, msgData)

	//marshal, _ := proto.Marshal(msg)
	//if _, err := js.Publish(AnimalPartyTopic, marshal); err != nil {
	//	global.GVA_LOG.Error("NewPublishMessages 发布消息失败: %v", zap.Error(err))
	//} else {
	//	global.GVA_LOG.Infof("NewPublishMessages 发布消息: %s", string(marshal))
	//}

	if err != nil {
		global.GVA_LOG.Error("Failed to send response: %v", zap.Error(err))
	}
}

func memeBattleEnsureStream(js nats.JetStreamContext) error {
	streamConfig := &nats.StreamConfig{
		Name:      config.AnimalParty,
		Subjects:  []string{config.AnimalPartyTopic, config.AnimalPartyTopicResp},
		Retention: nats.WorkQueuePolicy,
		//MaxBytes:  1 * 1024 * 1024 * 1024, // 1GB
		Storage:  nats.FileStorage, // 存储类型（文件存储）
		MaxAge:   2 * time.Hour,
		Replicas: 1,
	}

	// 检查流是否已存在
	if info, err := js.StreamInfo(config.AnimalParty); err == nil {
		global.GVA_LOG.Infof("流 %s 已存在，跳过创建\n", info.Config.Name)
		return nil
	}

	// 创建新流
	_, err := js.AddStream(streamConfig)
	return err

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
	js, err := n.NatsConn.JetStream()

	//js, err := jetstream.New(n.NatsConn)

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
