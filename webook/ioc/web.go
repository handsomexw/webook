package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	web.RegisterRoutes(server, hdl)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		ignoreHdl(),
		func(context *gin.Context) {
			fmt.Println("这是第一个路由")
		},
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Token"},
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
		IgnorePath("/user/login_sms/code/verify").Build()
}
