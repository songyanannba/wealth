// Package cache 缓存
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
)

const (
	userOnlinePrefix    = "l:city:user:online:" // 用户在线状态
	userOnlineCacheTime = 24 * 60 * 60
)

func getUserOnlineKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", userOnlinePrefix, userKey)
	return
}

// GetUserOnlineInfo 获取用户在线信息
func GetUserOnlineInfo(userKey string) (userOnline *models.UserOnline, err error) {
	redisClient := global.GVA_REDIS
	key := getUserOnlineKey(userKey)
	data, err := redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			global.GVA_LOG.Error("GetUserOnlineInfo", zap.Any(userKey, err))
			return
		}
		global.GVA_LOG.Error("GetUserOnlineInfo", zap.Any(userKey, err))
		return
	}
	userOnline = &models.UserOnline{}
	err = json.Unmarshal(data, userOnline)
	if err != nil {
		global.GVA_LOG.Error("获取用户在线数据 json Unmarshal ", zap.Any(userKey, err))
		return
	}
	global.GVA_LOG.Infof("获取用户在线数据 %v ,LoginTime %v , HeartbeatTime %v , AccIp %v , IsLogoff %v ", userKey, userOnline.LoginTime, userOnline.HeartbeatTime,
		userOnline.AccIp, userOnline.IsLogoff)
	return
}

// SetUserOnlineInfo 设置用户在线数据
func SetUserOnlineInfo(userKey string, userOnline *models.UserOnline) (err error) {
	redisClient := global.GVA_REDIS
	key := getUserOnlineKey(userKey)
	valueByte, err := json.Marshal(userOnline)
	if err != nil {
		global.GVA_LOG.Error("设置用户在线数据 json Marshal ", zap.Any(key, err))
		return
	}
	_, err = redisClient.Do(context.Background(), "setEx", key, userOnlineCacheTime, string(valueByte)).Result()
	if err != nil {
		global.GVA_LOG.Error("设置用户在线数据 ", zap.Any(key, err))
		return
	}
	return
}
