package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	FrequencyVerifyPrefix          = "wand:user:" //
	FrequencyVerifyPrefixCacheTime = time.Second * 5
)

func GetFrequencyVerifyKey(uid string) (key string) {
	key = fmt.Sprintf("%s%s", FrequencyVerifyPrefix, uid)
	return
}

func SetFrequencyVerify(uid, timeStr string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetFrequencyVerifyKey(uid)

	err = redisClient.Set(context.Background(), key, timeStr, FrequencyVerifyPrefixCacheTime).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// GetFrequencyVerify 销毁矿石（官方兑换）｜ 售卖矿石 ｜ 积分转赠（购买积分的接口） 需要请求这个接口
func GetFrequencyVerify(uid string) (val string, err error) {
	key := GetFrequencyVerifyKey(uid)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}

func DelFrequencyVerify(uid string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetFrequencyVerifyKey(uid)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}
