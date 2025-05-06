package cache

import (
	"context"
	"fmt"
	"gateway/global"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	UpdateUserScorePrefix          = "gateWay:update:user:score@" // todo
	UpdateUserScorePrefixCacheTime = time.Second * 60 * 60 * 24   //
)

func GetUpdateUserScoreKey(roomKey, uid, process string) (key string) {
	//gateWay:update:user:score:roomNo:uid:process
	//前缀+ 房间号 + 用户ID + 进程号
	key = fmt.Sprintf("%s%s%s%s", UpdateUserScorePrefix, roomKey, uid, process)
	return
}

func SetUpdateUserScore(roomKey, uid, process string) (err error) {
	redisClient := global.GVA_REDIS
	key := GetUpdateUserScoreKey(roomKey, uid, process)

	err = redisClient.Set(context.Background(), key, roomKey, UpdateUserScorePrefixCacheTime).Err()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func GetUpdateUserScore(roomKey, uid, process string) (val string, err error) {
	key := GetUpdateUserScoreKey(roomKey, uid, process)
	redisClient := global.GVA_REDIS
	result, err := redisClient.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return result, nil
}
