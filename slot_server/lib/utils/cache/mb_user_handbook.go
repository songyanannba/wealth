package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"slot_server/lib/global"
	"time"
)

const (
	MbUserHandbookKey          = "mb_user_handbook_item"
	MbUserHandbookKeyCacheTime = 24 * time.Hour
)

func MbUserHandbookKeyKeyKey(userId string, cardId int) (key string) {
	key = fmt.Sprintf("%s:%v:%v", MbUserHandbookKey, userId, cardId)
	return
}

func SetMbUserHandbook(userId string, cardId int, value string) (err error) {
	redisClient := global.GVA_REDIS
	key := MbUserHandbookKeyKeyKey(userId, cardId)
	err = redisClient.Set(context.Background(), key, value, MbUserHandbookKeyCacheTime).Err()
	if err != nil {
		return
	}
	return
}

func GetMbUserHandbook(userId string, cardId int) (val string, err error) {
	key := MbUserHandbookKeyKeyKey(userId, cardId)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}
