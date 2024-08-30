package tencent

import (
	mysms "basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/ratelimit"
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
	mysms.Name
}

func NewService(appId string, signName string, client *sms.Client, limiter ratelimit.Limiter) *Service {
	return &Service{
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		client:   client,
		limiter:  limiter,
		//os.Getenv(),
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//在原有的代码上修改，侵入式修改不行

	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("短信发送失败: %s, %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(args []string) []*string {
	return slice.Map[string, *string](args, func(idx int, src string) *string {
		return &src
	})
}
