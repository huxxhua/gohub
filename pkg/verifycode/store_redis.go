package verifycode

import (
	"gohub/pkg/app"
	"gohub/pkg/config"
	"gohub/pkg/redis"
	"time"
)

// RedisStore 实现 verifycode.Store interface
type RedisStore struct {
	RedisClient *redis.RedisClient
	KeyPrefix   string
}

// Set 实现 verifycode.Store interface 的 Set 方法
func (s *RedisStore) Set(id string, value string) bool {

	ExpireTime := time.Minute * time.Duration(config.GetInt64("verifycode.expire_time"))
	// 本地环境方便调试
	if app.IsLocal() {
		ExpireTime = time.Minute * time.Duration(config.GetInt64("verifycode.debug_expire_time"))
	}
	return s.RedisClient.Set(s.KeyPrefix+id, value, ExpireTime)
}

// Get 实现 verifycode.Store interface 的 Get 方法
func (s *RedisStore) Get(id string, clear bool) (value string) {
	value = s.RedisClient.Get(s.KeyPrefix + id)
	if clear {
		s.RedisClient.Del(id)
	}
	return value
}

// Verify 实现 verifycode.Store interface 的 Verify 方法
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	value := s.Get(id, clear)
	return value == answer
}
