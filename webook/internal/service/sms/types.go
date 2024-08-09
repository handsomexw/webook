package sms

import "context"

// 适配多个供应商
// 手机号，appid，签名，模板，参数
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
