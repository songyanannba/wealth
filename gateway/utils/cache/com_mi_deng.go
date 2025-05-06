package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	ComMiDengPrefix          = "com:mi:deng:"
	ComMiDengExpPreCacheTime = time.Second * 5
)

func GetComMiDengPrefixKey(uid, fun string) (key string) {
	key = fmt.Sprintf("%s%s%s", ComMiDengPrefix, uid, fun)
	return
}

// SetComMiDengFuncExpPre 幂等
func SetComMiDengFuncExpPre(uid string, fun, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetComMiDengPrefixKey(uid, fun)
	err = redisClient.Set(context.Background(), key, val, ComMiDengExpPreCacheTime).Err()
	if err != nil {
		return
	}
	return
}

// GetComMiDengFuncExpPre 根据uid和方法做的 通用 幂等
func GetComMiDengFuncExpPre(uid string, funcName string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GetComMiDengPrefixKey(uid, funcName)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetComMiDengFuncExpPre uid %v  err %v", err, uid)
		return "", err
	}
	global.GVA_LOG.Infof("GetComMiDengFuncExpPre uid %v key %v result {%v} err %v", uid, key, result, err)
	return result, nil
}

func DelComMiDengFuncExpPre(uid string, bet string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetComMiDengPrefixKey(uid, bet)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return
	}
	return
}
