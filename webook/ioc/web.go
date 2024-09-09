package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	ratelimit2 "basic-go/webook/pkg/ratelimit"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler,
	oauth2WechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	web.RegisterRoutes(server, hdl)
	oauth2WechatHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(limiter ratelimit2.Limiter) []gin.HandlerFunc {

	return []gin.HandlerFunc{
		corsHdl(),
		ignoreHdl(),
		func(context *gin.Context) {
			fmt.Println("这是第一个路由")
		},
		ratelimit.NewBuilder(limiter).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Token", "x-refresh-token"},
		//第一种方式
		//AllowOrigins: []string{"http://localhost:8081"},
		//第二种方式

		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") ||
				strings.HasPrefix(origin, "https://live.webook.com") {
				return true
			}
			fmt.Println("origin:", origin)
			return false
		},
		MaxAge: 12 * time.Hour,
	})
}

func ignoreHdl() gin.HandlerFunc {
	return middleware.NewLoginJwtMiddlewareBuilder().IgnorePath("/user/login").
		IgnorePath("/user/signup").IgnorePath("/user/login/jwt").
		IgnorePath("/user/login_sms/code/send").
		IgnorePath("/user/login_sms/code/verify").
		IgnorePath("/oauth2/wechat/authurl").
		IgnorePath("/oauth2/wechat/callback").
		Build()
}
