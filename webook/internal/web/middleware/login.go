package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	path []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(path ...string) *LoginMiddlewareBuilder {
	for _, ph := range path {
		l.path = append(l.path, ph)
	}
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		for _, path := range l.path {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			//?
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := sess.Get("updateTime")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		nowTime := time.Now().UnixMilli()
		if updateTime == nil {
			sess.Set("updateTime", nowTime)
			sess.Save()
			return
		}
		updateTimeVal, _ := updateTime.(int64)
		if nowTime-updateTimeVal > 50*1000 {
			sess.Set("updateTime", nowTime)
			sess.Save()
		}
	}
}
