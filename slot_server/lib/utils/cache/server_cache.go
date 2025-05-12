// Package cache 缓存
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"slot_server/lib/global"
	"slot_server/lib/models"
	"strconv"
)

const (
	serversHashKey       = "fish:hash:servers" // 全部的服务器
	serversHashCacheTime = 2 * 60 * 60         // key过期时间
	serversHashTimeout   = 3 * 60              // 超时时间
)

func getServersHashKey() (key string) {
	key = fmt.Sprintf("%s", serversHashKey)

	return
}

// SetServerInfo 设置服务器信息
func SetServerInfo(server *models.Server, currentTime uint64) (err error) {
	key := getServersHashKey()
	value := fmt.Sprintf("%d", currentTime)
	redisClient := global.GVA_REDIS
	number, err := redisClient.Do(context.Background(), "hSet", key, server.String(), value).Int()
	if err != nil {
		global.GVA_LOG.Error("SetServerInfo", zap.Any("number", number), zap.Any(key, err))
		return
	}
	redisClient.Do(context.Background(), "Expire", key, serversHashCacheTime)
	return
}

// DelServerInfo 下线服务器信息
func DelServerInfo(server *models.Server) (err error) {
	key := getServersHashKey()
	redisClient := global.GVA_REDIS
	number, err := redisClient.Do(context.Background(), "hDel", key, server.String()).Int()
	if err != nil {
		global.GVA_LOG.Error("DelServerInfo", zap.Any("number", number), zap.Any(key, err))
		return
	}
	if number != 1 {
		return
	}
	redisClient.Do(context.Background(), "Expire", key, serversHashCacheTime)
	return
}

// GetServerAll 获取所有服务器
func GetServerAll(currentTime uint64) (servers []*models.Server, err error) {
	servers = make([]*models.Server, 0)
	key := getServersHashKey()
	redisClient := global.GVA_REDIS
	val, err := redisClient.Do(context.Background(), "hGetAll", key).Result()
	valByte, _ := json.Marshal(val)
	global.GVA_LOG.Infof("GetServerAll ", zap.Any(key, string(valByte)))
	serverMap, err := redisClient.HGetAll(context.Background(), key).Result()
	if err != nil {
		global.GVA_LOG.Error("SetServerInfo", zap.Any(key, err))
		return
	}
	for key, value := range serverMap {
		valueUint64, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			global.GVA_LOG.Error("GetServerAll", zap.Any(key, err))
			return nil, err
		}

		// 超时
		if valueUint64+serversHashTimeout <= currentTime {
			continue
		}
		server, err := models.StringToServer(key)
		if err != nil {
			global.GVA_LOG.Error("GetServerAll", zap.Any(key, err))
			return nil, err
		}
		servers = append(servers, server)
	}
	return
}
