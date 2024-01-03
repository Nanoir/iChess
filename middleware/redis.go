// redis.go
package middleware

import (
	_ "context"
	"github.com/go-redis/redis/v8"
)

// 全局的 Redis 客户端
var redisClient *redis.Client

// 初始化 Redis 客户端
func InitRedisClient() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 服务器地址
		// 其他配置选项，比如密码等
	})
	return redisClient
}
