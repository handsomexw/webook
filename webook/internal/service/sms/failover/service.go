package failover

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"log"
	"sync/atomic"
)

type FailoverSMSService struct {
	//可能有多个服务商
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(services []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: services,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//缺点：大部分会在svcs[0]成功，负载不均衡；每次轮询，很耗时
	for _, service := range f.svcs {
		err := service.Send(ctx, tplId, args, numbers...)
		if err == nil {
			return err
		}
		log.Println(err)
	}
	return errors.New("发送失败，所以服务商都试过了")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//atomic原子操作
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i%length]
		err := svc.Send(ctx, tplId, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			log.Println(err)
		}
	}
	return errors.New("发送失败，全部服务商都试过")
}
