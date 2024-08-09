package middleware

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJtMiddlewareBuilder struct {
	path []string
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

		tokenHeader := ctx.GetHeader("Token")
		sesg := strings.Split(tokenHeader, " ")
		if len(sesg) != 1 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenHeader = sesg[0]
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

		now := time.Now()
		if myclaims.ExpiresAt.Sub(now) < time.Second*50 {
			myclaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			//t := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)
			newToken, _ := token.SignedString([]byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"))
			ctx.Header("Token", newToken)
		}

		ctx.Set("claims", myclaims)
	}
}
