package auth

import (
	"basic-go/webook/internal/service/sms"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key string
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string
}

func NewSMSService(svc sms.Service, key string) sms.Service {
	return &SMSService{
		svc: svc,
		key: key,
	}
}

func (s *SMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//tplid 模板id， 在这里代表业务方的token,token中解读出tplid，

	var tc Claims

	token, err := jwt.ParseWithClaims(tplId, &tc, func(t *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !token.Valid {
		return errors.New("invalid token")
	}

	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}
