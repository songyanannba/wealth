package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"slot_server/lib/global"
	"time"
)

const (
	MbCardConfigKey          = "mmb_card_config_version"
	MbCardConfigKeyCacheTime = 10 * time.Minute
)

func MbCardConfigKeyKey(version string) (key string) {
	key = fmt.Sprintf("%s:%s", MbCardConfigKey, version)
	return
}

func SetMbCardConfig(version string, value string) (err error) {
	redisClient := global.GVA_REDIS
	key := MbCardConfigKeyKey(version)
	err = redisClient.Set(context.Background(), key, value, MbCardConfigKeyCacheTime).Err()
	if err != nil {
		return
	}
	return
}

func GetMbCardConfig(version string) (val string, err error) {
	key := MbCardConfigKeyKey(version)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}

func DelMbCardConfig(version string) (err error) {
	redisClient := global.GVA_REDIS
	key := MbCardConfigKeyKey(version)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("DelMbCardConfig err %v", err)
		return err
	}
	return nil
}
