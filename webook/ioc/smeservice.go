package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/memory"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	return memory.NewService()
	//svc := ratelimit.NewRateLimitService(memory.NewService(), ratelimit2.NewRedisSlidingWindowLimiter(cmd, time.Second, 10))
	//return svc
}
