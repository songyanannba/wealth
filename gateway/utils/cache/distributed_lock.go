package cache

import (
	"context"
	"gateway/global"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"time"
)

//// TryLock 尝试获取分布式锁
//func TryLock(key string, timeout time.Duration) error {
//	val := fmt.Sprintf("uuid:%s", uuid.New().String())
//	redisClient := global.GVA_REDIS
//	cmd := redisClient.Set(context.Background(), key, val, timeout)
//	if err := cmd.Err(); err != nil {
//		return err
//	}
//	return nil
//}
//
//// ReleaseLock 释放分布式锁
//func ReleaseLock(key string) error {
//	val := fmt.Sprintf("uuid:%s", uuid.New().String())
//	redisClient := global.GVA_REDIS
//	cmd := redisClient.Del(context.Background(), key)
//	if err := cmd.Err(); err != nil {
//		return err
//	}
//	return nil
//}

// RedisLock 是分布式锁的核心结构
type RedisLock struct {
	Client  *redis.Client // Redis 客户端
	Key     string        // 锁的键名
	Value   string        // 锁的唯一标识
	Timeout time.Duration // 锁的过期时间
}

// NewRedisLock 创建一个 RedisLock 实例
func NewRedisLock(key string, timeout time.Duration) *RedisLock {
	return &RedisLock{
		Client:  global.GVA_REDIS,
		Key:     key,
		Value:   uuid.New().String(), // 生成唯一标识符
		Timeout: timeout,
	}
}

// Acquire 获取锁
func (l *RedisLock) Acquire(ctx context.Context) (bool, error) {
	// 使用 SETNX 和 EX 实现原子操作
	success, err := l.Client.SetNX(ctx, l.Key, l.Value, l.Timeout).Result()
	if err != nil {
		return false, err
	}
	return success, nil
}

// Release 释放锁
func (l *RedisLock) Release(ctx context.Context) (bool, error) {
	// Lua 脚本验证锁的拥有者并删除键
	script := `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
	`
	result, err := l.Client.Eval(ctx, script, []string{l.Key}, l.Value).Result()
	if err != nil {
		return false, err
	}
	return result.(int64) == 1, nil
}
