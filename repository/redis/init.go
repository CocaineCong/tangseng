package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/CocaineCong/tangseng/config"
	logs "github.com/CocaineCong/tangseng/pkg/logger"
)

// RedisClient Redis缓存客户端单例
var RedisClient *redis.Client
var RedisContext = context.Background()

// InitRedis 在中间件中初始化redis链接
func InitRedis() {
	rConfig := config.Conf.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", rConfig.RedisHost, rConfig.RedisPort),
		Username: rConfig.RedisUsername,
		Password: rConfig.RedisPassword,
		DB:       rConfig.RedisDbName,
	})
	_, err := client.Ping(RedisContext).Result()
	if err != nil {
		logs.LogrusObj.Errorln(err)
		panic(err)
	}
	RedisClient = client
}
