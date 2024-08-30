package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/ratelimit"
	"context"
	"fmt"
)

type RateLimitService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRateLimitService(svc sms.Service, l ratelimit.Limiter) sms.Service {
	return &RateLimitService{
		svc:     svc,
		limiter: l,
	}
}

func (s *RateLimitService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//
	limiter, err := s.limiter.Limit(ctx, "sms/login")
	if err != nil {
		return fmt.Errorf("limiter.Limit: %w", err)
	}
	if limiter {
		return fmt.Errorf("limiter.Limit: %w", err)
	}
	err = s.svc.Send(ctx, tplId, args, numbers...)
	return err
}
