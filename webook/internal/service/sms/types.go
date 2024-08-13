package sms

import "context"

// 适配多个供应商
// 手机号，appid，签名，模板，参数
type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
type Name struct {
}

// 更适配版本
type ServiceV1 interface {
	Send(ctx context.Context, tplId string, args []NameValue, numbers ...string) error
}

type NameValue struct {
	Val  string
	Name string
}
