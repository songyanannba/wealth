package websocket

import (
	"context"
	"errors"
	"gateway/common"
	"gateway/global"
	"gateway/protoc/pbs"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"time"
)

const (
	AnimalParty = "animal_party1" // 流名称

	AnimalPartyTopic        = "animal.party.topic1"   // 流绑定的主题
	AnimalPartyConsumerName = "animal_party_consumer" //消费者

	AnimalPartyTopicResp           = "animal.party.topi1.resp"    // 流绑定的 主题
	AnimalPartyProducerSubjectResp = "animal_party_resp_consumer" // 消费者
)

// memeBattleEnsureStream 确保流存在
//func memeBattleEnsureStream(js jetstream.JetStream) error {
//	_, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
//		Name:      AnimalParty,
//		Subjects:  []string{AnimalPartyTopic, AnimalPartyTopicResp},
//		Retention: jetstream.WorkQueuePolicy,
//		Storage:   jetstream.FileStorage, // 存储类型（文件存储）
//		MaxAge:    2 * time.Hour,         // 消息保留时间
//		Replicas:  1,
//	})
//	if err != nil && !errors.Is(err, jetstream.ErrStreamNameAlreadyInUse) {
//		return fmt.Errorf("创建流失败: %w", err)
//	}
//	return nil
//}

func slotServerEnsureStream(js nats.JetStreamContext) error {
	streamConfig := &nats.StreamConfig{
		Name:      AnimalParty,
		Subjects:  []string{AnimalPartyTopic, AnimalPartyTopicResp},
		Retention: nats.LimitsPolicy,
		MaxBytes:  1 * 1024 * 1024 * 1024, // 1GB
		Storage:   nats.FileStorage,       // 存储类型（文件存储）
		MaxAge:    2 * time.Hour,
		Replicas:  1,
	}

	// 检查流是否已存在
	if info, err := js.StreamInfo(AnimalParty); err == nil {
		global.GVA_LOG.Infof("流 %s 已存在，跳过创建\n", info.Config.Name)
		return nil
	}

	// 创建新流
	_, err := js.AddStream(streamConfig)
	return err
}

func (n *natsManager) NewPublishMessages(msg *pbs.NetMessage, js nats.JetStreamContext) {
	marshal, _ := proto.Marshal(msg)
	if _, err := js.Publish(AnimalPartyTopic, marshal); err != nil {
		global.GVA_LOG.Error("NewPublishMessages 发布消息失败: %v", zap.Error(err))
	} else {
		global.GVA_LOG.Infof("NewPublishMessages 发布消息: %s", string(marshal))
	}
}

func (n *natsManager) slotServiceSubConsumer() {
	time.Sleep(2 * time.Second)
	js, err := n.GetNatsJs(AnimalParty)
	if err != nil {
		global.GVA_LOG.Infof("memeBattleServiceSubConsumer %v", zap.Error(err))
		return
	}

	cons, err := js.PullSubscribe(AnimalPartyTopicResp, "AnimalPartyTopicResp_Consumer", nats.MaxRequestExpires(10*time.Second))
	if err != nil {
		global.GVA_LOG.Error("memeBattleServiceSubConsumer 创建消费者失败: ", zap.Error(err))
		return
	}
	for {
		msgs, err := cons.Fetch(100, nats.MaxWait(5*time.Second))

		if errors.Is(err, context.DeadlineExceeded) {
			continue
		} else if err != nil {
			global.GVA_LOG.Infof("拉取失败: %v", zap.Error(err))
			//time.Sleep(5 * time.Second)
			continue
		}

		global.GVA_LOG.Infof("Processing slotServiceSubConsumer:")

		// 处理消息
		for _, msgData := range msgs {
			global.GVA_LOG.Infof("Processing message: %s", string(msgData.Data))

			// 发送处理结果给生产者
			req := &pbs.NetMessage{}
			err := proto.Unmarshal(msgData.Data, req)
			if err != nil {
				global.GVA_LOG.Error("Error unmarshalling message: %v", zap.Error(err))
				continue
			}

			global.GVA_LOG.Infof("Processing memeBattleServiceSubConsumer: req: %v ,Content:%v", req, string(req.Content))

			if req.MsgId <= 0 {
				global.GVA_LOG.Infof("memeBattleServiceSubConsumer Skipping message because msgId is zero,%v", req)
				continue
			}

			// 确认收到消息
			if err := msgData.Ack(); err != nil {
				global.GVA_LOG.Infof("Failed to acknowledge message: %v", err)
			} else {
				global.GVA_LOG.Infof("Acknowledged: %s", string(msgData.Data))
			}

			if req.MsgId == 399 {
				global.GVA_LOG.Infof("协议消息心跳返回 不通知客户端")
				continue
			}

			//有的协议是广播 现在根据返回头是否存在uid来判断
			if req.AckHead.Uid == "" {
				//jsonData
				//clientManager.sendAppIDAll([]byte(), appID, ignoreClient)

			} else {
				uidStr := req.AckHead.Uid
				clientInfo := GetUserClient(common.AppId10, uidStr)
				if clientInfo == nil {
					global.GVA_LOG.Infof(" memeBattleServiceSubConsumer 用户没有客户端,用户可能没登陆 UserID:%v ", req)
					continue
				}

				//直接返回客户端
				if req.AckHead.Code != pbs.Code_OK {
					message := common.GetErrorMessage(uint32(req.AckHead.Code), "")
					req.AckHead.Message = message
				}
				ackReMarshal, _ := proto.Marshal(req)
				clientInfo.SendMsg(ackReMarshal)

				//处理返回的消息 返回客户端
				// 采用 map 注册的方式
				//if value, ok := getNatsProtoResp(req.MsgId); ok {
				//	var (
				//		code uint32
				//		//respMsgId uint32
				//		cmd  string
				//		data interface{}
				//	)
				//
				//	if req.AckHead.Code != pbs.Code_OK {
				//		cmd = strconv.Itoa(int(req.MsgId))
				//		code = uint32(req.AckHead.Code)
				//	} else {
				//		_, code, data = value(req.MsgId, req.Content)
				//		cmd = strconv.Itoa(int(req.MsgId))
				//		//message := common.GetErrorMessage(code, "")
				//		//responseHead := models.NewResponseHead("", cmd, code, message, data)
				//		//headByte, err := json.Marshal(responseHead)
				//		//if err != nil {
				//		//	global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
				//		//	continue
				//		//}
				//		//clientInfo.SendMsg(headByte)
				//		//global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", clientInfo.Addr, clientInfo.AppID, clientInfo.UserID, req.MsgId, code)
				//	}
				//
				//	ackMsgId, ackCode, contentByte := value(client, netMessage.MsgId, netMessage.Content)
				//	ackMsg := common.GetErrorMessage(ackCode, "")
				//	//client.SendMsg(headByte)
				//	netMessageResp := &pbs.NetMessage{
				//		ReqHead: netMessage.ReqHead,
				//		AckHead: &pbs.AckHead{
				//			Uid:     netMessage.ReqHead.Uid,
				//			Code:    pbs.Code(ackCode),
				//			Message: ackMsg,
				//		},
				//		ServiceId: netMessage.ServiceId,
				//		MsgId:     ackMsgId,
				//		Content:   contentByte,
				//	}
				//	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
				//	global.GVA_LOG.Infof("gate_way_response send headByte:%v ", string(contentByte))
				//	client.SendMsg(netMessageRespMarshal)
				//
				//	message := common.GetErrorMessage(code, "")
				//	responseHead := models.NewResponseHead("", cmd, code, message, data)
				//	headByte, err := json.Marshal(responseHead)
				//	if err != nil {
				//		global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
				//		continue
				//	}
				//	clientInfo.SendMsg(headByte)
				//	global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", clientInfo.Addr, clientInfo.AppID, clientInfo.UserID, req.MsgId, code)
				//
				//} else {
				//	global.GVA_LOG.Error("处理数据 路由不存在", zap.Any("Addr", clientInfo.Addr), zap.Any("cmd", req.MsgId))
				//	continue
				//}
			}

		}
	}
}

// 拉模式消费者服务
//func (n *natsManager) memeBattleServiceSubConsumer() {
//	time.Sleep(2 * time.Second)
//	js, err := n.GetNatsJs(AnimalParty)
//	if err != nil {
//		global.GVA_LOG.Infof("memeBattleServiceSubConsumer %v", zap.Error(err))
//		return
//	}
//
//	cons, err := js.CreateConsumer(context.Background(), AnimalParty, jetstream.ConsumerConfig{
//		Durable:       "AnimalPartyTopicResp_Consumer",
//		FilterSubject: AnimalPartyTopicResp,
//		AckPolicy:     jetstream.AckExplicitPolicy,
//		MaxDeliver:    3,
//		AckWait:       30 * time.Second,
//	})
//	if err != nil {
//		global.GVA_LOG.Error("memeBattleServiceSubConsumer 创建消费者失败: ", zap.Error(err))
//		return
//	}
//	for {
//		msgs, err := cons.Fetch(100, jetstream.FetchMaxWait(5*time.Second))
//
//		if errors.Is(err, context.DeadlineExceeded) {
//			continue
//		} else if err != nil {
//			log.Printf("拉取失败: %v", err)
//			time.Sleep(5 * time.Second)
//			continue
//		}
//
//		global.GVA_LOG.Infof("Processing memeBattleServiceSubConsumer:")
//
//		// 处理消息
//		for msg := range msgs.Messages() {
//			global.GVA_LOG.Infof("Processing message: %s", string(msg.Data()))
//
//			// 发送处理结果给生产者
//			req := &pbs.NetMessage{}
//			err := proto.Unmarshal(msg.Data(), req)
//			if err != nil {
//				global.GVA_LOG.Error("Error unmarshalling message: %v", zap.Error(err))
//				continue
//			}
//
//			global.GVA_LOG.Infof("Processing memeBattleServiceSubConsumer: req: %v ,Content:%v", req, string(req.Content))
//
//			if req.MsgId <= 0 {
//				global.GVA_LOG.Infof("memeBattleServiceSubConsumer Skipping message because msgId is zero,%v", req)
//				continue
//			}
//
//			// 确认收到消息
//			if err := msg.Ack(); err != nil {
//				global.GVA_LOG.Infof("Failed to acknowledge message: %v", err)
//			} else {
//				global.GVA_LOG.Infof("Acknowledged: %s", string(msg.Data()))
//			}
//
//			if req.MsgId == 399 {
//				global.GVA_LOG.Infof("协议消息心跳返回 不通知客户端")
//				continue
//			}
//
//			//有的协议是广播 现在根据返回头是否存在uid来判断
//			if req.AckHead.Uid == "" {
//				//jsonData
//				//clientManager.sendAppIDAll([]byte(), appID, ignoreClient)
//
//			} else {
//				uidStr := req.AckHead.Uid
//				clientInfo := GetUserClient(common.AppId10, uidStr)
//				if clientInfo == nil {
//					global.GVA_LOG.Infof(" memeBattleServiceSubConsumer 用户没有客户端,用户可能没登陆 UserID:%v ", req)
//					continue
//				}
//
//				//直接返回客户端
//				if req.AckHead.Code != pbs.Code_OK {
//					message := common.GetErrorMessage(uint32(req.AckHead.Code), "")
//					req.AckHead.Message = message
//				}
//				ackReMarshal, _ := proto.Marshal(req)
//				clientInfo.SendMsg(ackReMarshal)
//
//				//处理返回的消息 返回客户端
//				// 采用 map 注册的方式
//				//if value, ok := getNatsProtoResp(req.MsgId); ok {
//				//	var (
//				//		code uint32
//				//		//respMsgId uint32
//				//		cmd  string
//				//		data interface{}
//				//	)
//				//
//				//	if req.AckHead.Code != pbs.Code_OK {
//				//		cmd = strconv.Itoa(int(req.MsgId))
//				//		code = uint32(req.AckHead.Code)
//				//	} else {
//				//		_, code, data = value(req.MsgId, req.Content)
//				//		cmd = strconv.Itoa(int(req.MsgId))
//				//		//message := common.GetErrorMessage(code, "")
//				//		//responseHead := models.NewResponseHead("", cmd, code, message, data)
//				//		//headByte, err := json.Marshal(responseHead)
//				//		//if err != nil {
//				//		//	global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
//				//		//	continue
//				//		//}
//				//		//clientInfo.SendMsg(headByte)
//				//		//global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", clientInfo.Addr, clientInfo.AppID, clientInfo.UserID, req.MsgId, code)
//				//	}
//				//
//				//	ackMsgId, ackCode, contentByte := value(client, netMessage.MsgId, netMessage.Content)
//				//	ackMsg := common.GetErrorMessage(ackCode, "")
//				//	//client.SendMsg(headByte)
//				//	netMessageResp := &pbs.NetMessage{
//				//		ReqHead: netMessage.ReqHead,
//				//		AckHead: &pbs.AckHead{
//				//			Uid:     netMessage.ReqHead.Uid,
//				//			Code:    pbs.Code(ackCode),
//				//			Message: ackMsg,
//				//		},
//				//		ServiceId: netMessage.ServiceId,
//				//		MsgId:     ackMsgId,
//				//		Content:   contentByte,
//				//	}
//				//	netMessageRespMarshal, _ := proto.Marshal(netMessageResp)
//				//	global.GVA_LOG.Infof("gate_way_response send headByte:%v ", string(contentByte))
//				//	client.SendMsg(netMessageRespMarshal)
//				//
//				//	message := common.GetErrorMessage(code, "")
//				//	responseHead := models.NewResponseHead("", cmd, code, message, data)
//				//	headByte, err := json.Marshal(responseHead)
//				//	if err != nil {
//				//		global.GVA_LOG.Infof("处理数据 json Marshal %v", err)
//				//		continue
//				//	}
//				//	clientInfo.SendMsg(headByte)
//				//	global.GVA_LOG.Infof("gate_way_response send %v %v %v cmd %vcode %v ", clientInfo.Addr, clientInfo.AppID, clientInfo.UserID, req.MsgId, code)
//				//
//				//} else {
//				//	global.GVA_LOG.Error("处理数据 路由不存在", zap.Any("Addr", clientInfo.Addr), zap.Any("cmd", req.MsgId))
//				//	continue
//				//}
//			}
//
//		}
//	}
//}
