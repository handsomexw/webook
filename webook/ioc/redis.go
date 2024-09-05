package ioc

import (
	"basic-go/webook/config"
	"basic-go/webook/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Redisaddr.Addr,
		Password: "",
		DB:       0,
	})
	return redisClient
}

func InitLimiter(cmd redis.Cmdable) ratelimit.Limiter {
	return ratelimit.NewRedisSlidingWindowLimiter(cmd, time.Second, 10)
}
