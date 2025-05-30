package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	FishUserDaySignInPrefix       = "fish:user:day_sign_in"
	FishUserDaySignInPreCacheTime = time.Hour * 24
)

func FishUserDaySignInPreKey(uid, fun string) (key string) {
	key = fmt.Sprintf("%s%s%s", FishUserDaySignInPrefix, uid, fun)
	return
}

// SetFishUserDaySignInExpPre 每日签到的标识
func SetFishUserDaySignInExpPre(uid string, fun, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := FishUserDaySignInPreKey(uid, fun)
	err = redisClient.Set(context.Background(), key, val, FishUserDaySignInPreCacheTime).Err()
	if err != nil {
		return
	}
	return
}

// GetFishUserDaySignInPre 获取方法的调用间隔时间
func GetFishUserDaySignInPre(uid string, funcName string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := FishUserDaySignInPreKey(uid, funcName)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetFishUserDaySignInPre uid %v  err %v", err, uid)
		return "", err
	}
	global.GVA_LOG.Infof("GetFishUserDaySignInPre uid %v key %v result {%v} err %v", uid, key, result, err)
	return result, nil
}
