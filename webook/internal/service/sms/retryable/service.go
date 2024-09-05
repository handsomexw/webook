package retryable

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"fmt"
)

type Service struct {
	svc      sms.Service
	retryMax int
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	err := s.svc.Send(ctx, tplId, args, numbers...)
	cnt := 1
	for err != nil && cnt <= s.retryMax {
		err = s.svc.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return nil
		}
		cnt++
	}
	return errors.New(fmt.Sprintf("%d次重试，失败", cnt-1))
}
