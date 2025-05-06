package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	RemakeAgePrefix    = "remake:age:"
	RemakeAgeCacheTime = time.Second * 5
)

func RemakeAgePrefixKey(uid, val string) (key string) {
	key = fmt.Sprintf("%s%s%s", RemakeAgePrefix, uid, val)
	return
}

// SetRemakeAge 幂等
func SetRemakeAge(uid string, val string) (err error) {
	redisClient := global.GVA_REDIS
	key := RemakeAgePrefixKey(uid, val)
	err = redisClient.Set(context.Background(), key, val, RemakeAgeCacheTime).Err()
	if err != nil {
		return
	}
	return
}

// GetRemakeAge 幂等
func GetRemakeAge(uid string, val string) (value string, err error) {
	redisClient := global.GVA_REDIS
	key := RemakeAgePrefixKey(uid, val)
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		global.GVA_LOG.Infof("GetRemakeAge uid %v  err %v", err, uid)
		return "", err
	}
	global.GVA_LOG.Infof("GetRemakeAge uid %v key %v result {%v} err %v", uid, key, result, err)
	return result, nil
}

func DelRemakeAge(uid string, bet string) (err error) {
	redisClient := global.GVA_REDIS
	key := RemakeAgePrefixKey(uid, bet)
	err = redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return
	}
	return
}
