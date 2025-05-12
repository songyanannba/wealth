package initialize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"slot_server/lib/global"
)

func Redis() {
	redisCfg := global.GVA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:         redisCfg.Addr,
		Password:     redisCfg.Password, // no password set
		DB:           redisCfg.DB,       // use default DB
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.GVA_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.GVA_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.GVA_REDIS = client
	}
}

//func Redis() {
//	redisCfg := global.GVA_CONFIG.Redis
//
//	client := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs:    []string{redisCfg.Addr},
//		Password: redisCfg.Password, // no password set
//	})
//
//	pong, err := client.Ping(context.Background()).Result()
//	if err != nil {
//		panic(fmt.Sprintf("redis connect ping failed, err:%s", err.Error()))
//	} else {
//		global.GVA_LOG.Info("redis connect ping response:", zap.String("pong", pong))
//		global.GVA_REDIS = client
//	}
//}
