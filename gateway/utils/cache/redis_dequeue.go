package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

// 定义队列名称
const Queue_UserSellOre = "UserSellOre"

// Enqueue 入队（生产者）
func Enqueue(queueName string, task string) error {
	// 使用 LPUSH 将任务插入队列
	rdb := global.GVA_REDIS
	err := rdb.LPush(context.Background(), queueName, task).Err()
	if err != nil {
		global.GVA_LOG.Infof("queueName %v err:%v", Queue_UserSellOre, err)
		return err
	}
	return nil
}

// Dequeue 出队（消费者，非阻塞）
func Dequeue(queueName string) (string, error) {
	rdb := global.GVA_REDIS
	// 使用 RPOP 获取队列中的任务
	task, err := rdb.RPop(context.Background(), queueName).Result()
	if err == redis.Nil {
		global.GVA_LOG.Infof("队列为空 %v err:%v", Queue_UserSellOre, err)
		return "", fmt.Errorf("队列为空")
	} else if err != nil {
		global.GVA_LOG.Infof("出队失败 %v err:%v", Queue_UserSellOre, err)
		return "", fmt.Errorf("出队失败: %w", err)
	}
	return task, nil
}

// DequeueBlocking 阻塞式出队（消费者）
func DequeueBlocking(queueName string, timeout time.Duration) (string, error) {
	rdb := global.GVA_REDIS
	// 使用 BRPop 阻塞式获取队列中的任务
	result, err := rdb.BRPop(context.Background(), timeout, queueName).Result()
	if err == redis.Nil {
		global.GVA_LOG.Infof("队列为空 %v err:%v", Queue_UserSellOre, err)
		return "", fmt.Errorf("队列为空")
	} else if err != nil {
		global.GVA_LOG.Infof("阻塞式出队失败 %v err:%v", Queue_UserSellOre, err)
		return "", fmt.Errorf("阻塞式出队失败: %w", err)
	}
	// BRPop 返回一个切片，格式为 [queueName, task]
	if len(result) > 1 {
		return result[1], nil
	}
	return "", fmt.Errorf("任务格式错误")
}

func ExecDequeue(queueName string) {
	for {
		result, _ := global.GVA_REDIS.LLen(context.Background(), queueName).Result()
		global.GVA_LOG.Infof("上架任务 还有多少未执行:%v", result)

		task, err := DequeueBlocking(queueName, 10*time.Second) // 阻塞式出队
		if err != nil {
			global.GVA_LOG.Infof("ExecDequeue 获取任务失败 %v err:%v", Queue_UserSellOre, err)
		} else {
			//reqData := pbs.UserSellOreReq{}
			//err := proto.Unmarshal([]byte(task), &reqData)
			//reqDataMarshal, _ := proto.Marshal(&reqData)
			//// 调用 gRPC 方法
			//msgReq := pbs.NetMessage{
			//	ReqHead: &pbs.ReqHead{
			//		Uid:      0,
			//		Token:    "",
			//		Platform: "",
			//	},
			//	AckHead:   &pbs.AckHead{},
			//	ServiceId: "",
			//	MsgId:     int32(pbs.ProtocNum_PNUserSellOreReq),
			//	Content:   reqDataMarshal,
			//}
			//response, err := grpcclient.GetClient().CallMethod(&msgReq)
			//if err != nil {
			//	global.GVA_LOG.Error("ExecDequeue could not call method:", zap.Error(err))
			//	continue
			//}
			//if response.AckHead.Code != pbs.Code_OK {
			//	global.GVA_LOG.Infof("ExecDequeue response:%v", response)
			//	continue
			//}
			global.GVA_LOG.Infof("ExecDequeue 任务已处理 %v task:%v", Queue_UserSellOre, task)
		}
	}
}
