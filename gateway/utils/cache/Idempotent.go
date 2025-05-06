package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

// Idempotent
const (
	IdempotentPrefix      = "fish:idempotent:" // 用户在线状态
	IdempotentCacheTime   = time.Second * 300
	IdempotentCacheTimeNx = time.Second * 2
)

// GetIdempotentKey 幂等行
func GetIdempotentKey(uid, mark string) (key string) {
	key = fmt.Sprintf("%s%s%s", IdempotentPrefix, uid, mark)
	return
}

func SetIdempotent(uid, mark, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetIdempotentKey(uid, mark)
	err = redisClient.Set(context.Background(), key, val, IdempotentCacheTime).Err()
	if err != nil {
		return
	}
	return
}

func GetIdempotent(uid, mark string) (val string, err error) {
	redisClient := global.GVA_REDIS
	key := GetIdempotentKey(uid, mark)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}

// SetIdempotentNx 限制接口访问太频繁
func SetIdempotentNx(uid, mark, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetIdempotentKey(uid, mark)
	err = redisClient.SetNX(context.Background(), key, val, IdempotentCacheTimeNx).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func releaseLock(uid, mark string) {
	redisClient := global.GVA_REDIS
	key := GetIdempotentKey(uid, mark)
	_, err := redisClient.Del(context.Background(), key).Result()
	if err != nil {
		fmt.Println("Error while releasing lock:", err)
	}
}

func SetIdempotentNxExp(uid, mark, val string, exp int) (err error) {
	redisClient := global.GVA_REDIS
	key := GetIdempotentKey(uid, mark)
	duration := time.Duration(exp)
	err = redisClient.SetNX(context.Background(), key, val, IdempotentCacheTimeNx*duration).Err()
	if err != nil {
		return
	}
	return
}
