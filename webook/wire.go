//go:build wireinject

package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	//中间件如何处理
	//最后返回的类型要行程调用链
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewUserDao, cache.NewCodeCache,
		cache.NewUserCache, repository.NewUserRepository,
		repository.NewCodeRepository, service.NewCodeService,
		service.NewUserService, ioc.InitSMSService, web.NewUserHandler,
		ioc.InitWebServer, ioc.InitMiddlewares, ioc.InitLimiter,
		ioc.InitOAuth2WechatService, web.NewOAuth2WechatHandler,
	)
	return new(gin.Engine)

}
