package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	OreSellCoinPrefix          = "ore:sell:coin"
	OreSellCoinExpPreCacheTime = time.Second * 5
)

func GetOreSellCoinPrefixKey(uid, fun string) (key string) {
	key = fmt.Sprintf("%s%s%s", OreSellCoinPrefix, uid, fun)
	return
}

// SetOreSellCoinFuncExpPre 购买矿币
func SetOreSellCoinFuncExpPre(uid string, fun, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetOreSellCoinPrefixKey(uid, fun)
	err = redisClient.Set(context.Background(), key, val, OreSellCoinExpPreCacheTime).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

// GetOreSellCoinFuncExpPre 购买矿币
func GetOreSellCoinFuncExpPre(uid string, funcName string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GetOreSellCoinPrefixKey(uid, funcName)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetOreSellCoinFuncExpPre uid %v  err %v", err, uid)
		return "", err
	}
	global.GVA_LOG.Infof("GetOreSellCoinFuncExpPre uid %v key %v result {%v} err %v", uid, key, result, err)
	return result, nil
}

func DelOreSellCoinFuncExpPre(uid string, bet string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetOreSellCoinPrefixKey(uid, bet)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return
	}
	return
}
