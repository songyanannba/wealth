package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	FishUserFuncPrefix          = "fish:user:func"
	FishUserFuncExpPreCacheTime = time.Second * 2
)

func GetUserFuncPreKey(uid, fun string) (key string) {
	key = fmt.Sprintf("%s%s%s", FishUserFuncPrefix, uid, fun)
	return
}

// SetFishUserFuncExpPre 设置方法的调用间隔时间
func SetFishUserFuncExpPre(uid string, fun, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetUserFuncPreKey(uid, fun)
	err = redisClient.Set(context.Background(), key, val, FishUserFuncExpPreCacheTime).Err()
	if err != nil {
		return
	}
	return
}

// GetFishUserFuncExpPre 获取方法的调用间隔时间
func GetFishUserFuncExpPre(uid string, funcName string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GetUserFuncPreKey(uid, funcName)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetFishUserFuncExpPre uid %v  err %v", err, uid)
		return "", err
	}
	global.GVA_LOG.Infof("GetFishUserFuncExpPre uid %v key %v result {%v} err %v", uid, key, result, err)
	return result, nil
}

func DelFishUserFuncExpPre(uid string, bet string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetUserFuncPreKey(uid, bet)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return
	}
	return
}
