package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	idx  int32
	cnt  int32
	//阈值
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)

	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		//考虑并发安全
		//这行
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		//这行
		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)

	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		//不知道什么错误
		//可以进行切换，超时可能是偶发性的，可以重试，但是发送错误就直接换
		//这行，三行有联动
		atomic.StoreInt32(&t.idx, (idx+1)%int32(len(t.svcs)))
		atomic.StoreInt32(&t.cnt, 0)
	}

	return err
}
