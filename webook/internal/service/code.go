package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var codeTplId string = "1122334"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func (c *CodeService) Send(ctx context.Context, biz string, phone string) error {
	//biz 用于区别场景
	//三个个步骤：
	//生成验证码、塞进Redis、发送
	code := c.generateCode()
	err := c.repo.Story(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	return c.smsSvc.Send(ctx, codeTplId, []string{}, phone)

	//if err != nil {
	//	//redis有验证码，但是发送失败，发送失败有两个原因，不确定
	//	//发送超时
	//	return err
	//}
}

func (c *CodeService) Verify(ctx context.Context, biz string, phone string, inputcode string) (bool, error) {
	//使用redis phone_code:biz:phone
	//要做原子操作
	//
	return c.repo.Verify(ctx, biz, phone, inputcode)

}

func (c *CodeService) generateCode() string {
	//生成 0-999999随机数
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}

// func (c *CodeService) VerifyV1(ctx context.Context, biz string, phone string, inputcode string) error {
//
// }
