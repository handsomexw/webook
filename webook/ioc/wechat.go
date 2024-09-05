package ioc

import "basic-go/webook/internal/service/oauth2/wechat"

func InitOAuth2WechatService() wechat.Service {
	appID := "123456"
	appSecret := "666"
	return wechat.NewWetchatService(appID, appSecret)
}
