package async

//
//import (
//	"basic-go/webook/internal/service/sms"
//	"context"
//	"time"
//)
//
//type AsyncService struct {
//	svc sms.Service
//}
//
//func NewAsyncService(svc sms.Service) sms.Service {
//	return &AsyncService{
//		svc: svc,
//	}
//}
//
//func (as *AsyncService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
//
//}
//
//func (as *AsyncService) needAsync() bool {
//	//触发异步的方案
//	//1.1基于响应时间，平均响应时间
//	//何时退出异步
//	//1.进入异步N分钟后
//	//2.保留1%的流量，0-100随机数，小于10的同步
//}
//
//func (as *AsyncService) AsyncSend(ctx context.Context, tplId string, args []string, numbers ...string) error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//
//}
