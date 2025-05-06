package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gateway/global"
	"gateway/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	gateWayUserWebPrefix    = "gate_way:user:web:" // 用户在线状态
	gateWayUserWebCacheTime = 48 * 60 * 60
)

func gateWayUserWebKey(userKey string) (key string) {
	key = fmt.Sprintf("%s%s", gateWayUserWebPrefix, userKey)
	return
}

// GetGateWayUserWeb 获取用户在线信息
func GetGateWayUserWeb(userKey string) (userWeb *models.UserWeb, err error) {
	redisClient := global.GVA_REDIS
	key := gateWayUserWebKey(userKey)
	data, err := redisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			global.GVA_LOG.Error("GetUserOnlineInfo", zap.Any(userKey, err))
			return
		}
		global.GVA_LOG.Error("GetUserOnlineInfo", zap.Any(userKey, err))
		return
	}
	userWeb = &models.UserWeb{}
	err = json.Unmarshal(data, userWeb)
	if err != nil {
		global.GVA_LOG.Error("获取用户在线数据 web json Unmarshal ", zap.Any(userKey, err))
		return
	}
	global.GVA_LOG.Infof("获取用户在线数据 web:%v ", userWeb)
	return
}

// SetGateWayUserWeb 设置用户在线数据
func SetGateWayUserWeb(userKey string, userWeb *models.UserWeb) (err error) {
	redisClient := global.GVA_REDIS
	key := gateWayUserWebKey(userKey)

	valueByte, err := json.Marshal(userWeb)
	if err != nil {
		global.GVA_LOG.Error("设置用户在线数据  web json Marshal ", zap.Any(key, err))
		return
	}

	_, err = redisClient.Do(context.Background(), "setEx", key, gateWayUserWebCacheTime, string(valueByte)).Result()
	if err != nil {
		global.GVA_LOG.Error("设置用户在线数据 web ", zap.Any(key, err))
		return
	}
	return
}
