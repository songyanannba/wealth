package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	SelectRoomPrefix    = "fish:select:room:" // 用户在线状态
	SelectRoomCacheTime = time.Second * 10
)

// GetSelectRoomKey 设置分布式锁
func GetSelectRoomKey(roomKey string) (key string) {
	key = fmt.Sprintf("%s%s", SelectRoomPrefix, roomKey)
	return
}

func SetProhibitSelectRoom(roomKey string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetSelectRoomKey(roomKey)
	err = redisClient.SetNX(context.Background(), key, roomKey, SelectRoomCacheTime).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func GetProhibitSelectRoom(roomKey string) (val string, err error) {
	key := GetSelectRoomKey(roomKey)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}
