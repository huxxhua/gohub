// Package redis 工具包
package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gohub/pkg/logger"
	"sync"
	"time"
)

// RedisClient Redis 服务
type RedisClient struct {
	Clint   *redis.Client
	Content context.Context
}

// once 确保全局的 Redis 对象只实例一次
var once sync.Once

// Redis 全局 Redis, 使用 db 1
var Redis *RedisClient

// ConnectRedis 连接 redis 数据库,设置全局 Redis 对象
func ConnectRedis(address string, username string, password string, db int) {
	once.Do(func() {
		Redis = NewClient(address, username, password, db)
	})
}

// NewClient 创建一个新的 redis 连接
func NewClient(address string, username string, password string, db int) *RedisClient {

	// 初始化自定的 RedisClient 实例
	rds := &RedisClient{}
	// 使用默认的 context
	rds.Content = context.Background()

	// 使用 redis 库里的 NewClient 初始化操作
	redis.NewClient(&redis.Options{
		Addr:     address,
		Username: username,
		Password: password,
		DB:       db,
	})

	// 测试下连接
	err := rds.Ping()
	logger.LogIf(err)

	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds RedisClient) Ping() error {
	_, err := rds.Clint.Ping(rds.Content).Result()

	return err
}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Clint.Set(rds.Content, key, value, expiration).Err(); err != nil {
		logger.ErrorString("Redis", "Set", err.Error())
		return false
	}
	return true
}

// Get 获取 key 对应的 value
func (rds RedisClient) Get(key string) string {
	res, err := rds.Clint.Get(rds.Content, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Get", err.Error())
		}
		return ""
	}
	return res
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds RedisClient) Has(key string) bool {
	_, err := rds.Clint.Get(rds.Content, key).Result()
	if err != nil {
		if err != redis.Nil {
			logger.ErrorString("Redis", "Get", err.Error())
		}
		return false
	}
	return true
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds RedisClient) Del(key ...string) bool {
	if err := rds.Clint.Del(rds.Content, key...).Err(); err != nil {
		logger.ErrorString("Redis", "Del", err.Error())
		return false
	}
	return true
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds RedisClient) FlushDB() bool {
	if err := rds.Clint.FlushDB(rds.Content).Err(); err != nil {
		logger.ErrorString("Redis", "FlushDB", err.Error())
		return false
	}
	return true
}

// Increment 当参数只有 1 个时，为 key，其值增加 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型
func (rds RedisClient) Increment(keys ...interface{}) bool {
	switch len(keys) {
	case 1:
		key := keys[0].(string)
		if err := rds.Clint.Incr(rds.Content, key).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	case 2:
		key := keys[0].(string)
		val := keys[1].(int64)

		if err := rds.Clint.IncrBy(rds.Content, key, val).Err(); err != nil {
			logger.ErrorString("Redis", "Increment", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Increment", "参数过多")
		return false
	}
	return true
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型
func (rds RedisClient) Decrement(keys ...interface{}) bool {
	switch len(keys) {
	case 1:
		key := keys[0].(string)
		if err := rds.Clint.Decr(rds.Content, key).Err(); err != nil {
			logger.ErrorString("Redis", "Decrement", err.Error())
			return false
		}
	case 2:
		key := keys[0].(string)
		val := keys[1].(int64)
		if err := rds.Clint.DecrBy(rds.Content, key, val).Err(); err != nil {
			logger.ErrorString("Redis", "Decrement", err.Error())
			return false
		}
	default:
		logger.ErrorString("Redis", "Decrement", "参数过多")
		return false
	}
	return true
}
