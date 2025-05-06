package cache

import (
	"context"
	"fmt"
	"gateway/global"

	"github.com/redis/go-redis/v9"
	"time"
)

const (
	ServiceConfigPrefix    = "service:conf"
	ServiceConfigCacheTime = time.Minute
)

func GetServiceConfigKey() (key string) {
	key = fmt.Sprintf("%s", ServiceConfigPrefix)
	return
}

func SetGetServiceConfigKeyExpPre(val int) (err error) {
	redisClient := global.GVA_REDIS
	key := GetServiceConfigKey()
	err = redisClient.Set(context.Background(), key, val, ServiceConfigCacheTime).Err()
	if err != nil {
		return
	}
	return
}

func GetGetServiceConfigKeyExpPre() (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GetServiceConfigKey()

	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetGetServiceConfigKeyExpPre err:%v", err)
		return "", err
	}
	global.GVA_LOG.Infof("GetGetServiceConfigKeyExpPre key %v result {%v} err %v", key, result, err)
	return result, nil
}
