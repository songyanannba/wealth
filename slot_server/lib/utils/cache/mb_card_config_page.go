package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"slot_server/lib/global"
	"time"
)

const (
	MbCardConfigPageKey          = "mb_card_config_page"
	MbCardConfigPageKeyCacheTime = 10 * time.Minute
)

func MbCardConfigPageKeyKey(lastId, level int) (key string) {
	key = fmt.Sprintf("%s:%v:%v", MbCardConfigPageKey, lastId, level)
	return
}

func SetMbCardConfigPage(lastId, level int, value string) (err error) {
	redisClient := global.GVA_REDIS
	key := MbCardConfigPageKeyKey(lastId, level)
	err = redisClient.Set(context.Background(), key, value, MbCardConfigPageKeyCacheTime).Err()
	if err != nil {
		return
	}
	return
}

// GetMbCardConfigPage 配置分野
func GetMbCardConfigPage(lastId, level int) (val string, err error) {
	key := MbCardConfigPageKeyKey(lastId, level)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}

//func DelMbCardConfigPage(version string) (err error) {
//	redisClient := global.GVA_REDIS
//	key := MbCardConfigPageKeyKey(version)
//	err = redisClient.Del(context.Background(), key).Err()
//	if err != nil && err != redis.Nil {
//		global.GVA_LOG.Infof("DelMbCardConfig err %v", err)
//		return err
//	}
//	return nil
//}
