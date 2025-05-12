package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"slot_server/lib/global"
	"time"
)

// 根据rpc协议ID控制接口是否可以访问 缓存
// GameServiceApisAuth

const (
	GameServiceApisAuthPrefix    = "game:service:apis:auth"
	GameServiceApisAuthCacheTime = time.Second * 60
)

func GameServiceApisAuthKey(msgId string) (key string) {
	key = fmt.Sprintf("%s%s", GameServiceApisAuthPrefix, msgId)
	return
}

func SetGameServiceApisAuth(msgId string, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GameServiceApisAuthKey(msgId)
	err = redisClient.Set(context.Background(), key, val, GameServiceApisAuthCacheTime).Err()
	if err != nil {
		global.GVA_LOG.Infof("SetGameServiceApisAuth err %v", err)
		return
	}
	return
}

func GetGameServiceApisAuth(msgId string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GameServiceApisAuthKey(msgId)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetGameServiceApisAuth err %v", err)
		return "", err
	}
	return result, nil
}

func DelGameServiceApisAuth(msgId string) (err error) {
	redisClient := global.GVA_REDIS
	key := GameServiceApisAuthKey(msgId)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("DelGameServiceApisAuth msgId %v  err %v", err, msgId)
		return err
	}
	return nil
}
