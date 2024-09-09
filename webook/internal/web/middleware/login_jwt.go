package middleware

import (
	"basic-go/webook/internal/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type LoginJtMiddlewareBuilder struct {
	path []string
	cmd  redis.Cmdable
}

func NewLoginJwtMiddlewareBuilder() *LoginJtMiddlewareBuilder {
	return &LoginJtMiddlewareBuilder{}
}

func (l *LoginJtMiddlewareBuilder) IgnorePath(path ...string) *LoginJtMiddlewareBuilder {
	for _, ph := range path {
		l.path = append(l.path, ph)
	}
	return l
}

func (l *LoginJtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.path {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//
		tokenHeader := web.ExtractToken(ctx)
		myclaims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenHeader, myclaims, func(token *jwt.Token) (interface{}, error) {
			return []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || myclaims.UserId == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if myclaims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题,这些错误我觉得应该由前端重定向到登录界面？当然后端也可以
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		//now := time.Now()
		//if myclaims.ExpiresAt.Sub(now) < time.Second*50 {
		//	myclaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	//t := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)
		//	newToken, _ := token.SignedString([]byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"))
		//	ctx.Header("Token", newToken)
		//}
		//还有降级策略，如果redis崩了，那就直接登录，不判断
		cnt, err := l.cmd.Exists(ctx, fmt.Sprintf("userd:ssid:%d", myclaims.UserId)).Result()
		if err != nil || cnt > 0 {
			//系统错误，或者用户已经退出了
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Set("claims", myclaims)
	}
}
