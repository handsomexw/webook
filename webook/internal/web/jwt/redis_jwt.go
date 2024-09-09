package jwt

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

var (
	AtKey = []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$")
	RtKey = []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU1$")
)

type RedisJwtHandler struct {
	cmd redis.Cmdable
}

type UserClaims struct {
	jwt.RegisteredClaims
	Ssid      string `json:"ssid"`
	UserId    int64  `json:"userid"`
	UserAgent string `json:"useragent"`
}

type RefreshClaims struct {
	jwt.RegisteredClaims
	Ssid   string `json:"ssid"`
	UserId int64  `json:"userid"`
}

func NewRedisJwtHandler(cmd redis.Cmdable) Handler {
	return &RedisJwtHandler{
		cmd: cmd,
	}
}

func (r *RedisJwtHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Token")
	msgs := strings.Split(tokenHeader, " ")
	if len(msgs) != 1 {
		return " "
	}
	return msgs[0]
}

func (r *RedisJwtHandler) SetJwtToken(ctx *gin.Context, uid int64, ssid string) error {
	myclaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)

	//token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString(AtKey)
	if err != nil {
		//ctx.JSON(http.StatusOK, gin.H{
		//	"message": "jwt系统错误",
		//	"error":   err.Error(),
		//})
		return err
	}
	ctx.Header("Token", tokenStr)
	return nil
}

func (r *RedisJwtHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("Token", "")
	ctx.Header("x-refresh-token", "")

	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
	}
	return r.cmd.Set(ctx, fmt.Sprintf("user:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}

func (r *RedisJwtHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New()
	err := r.SetJwtToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, uid, ssid)
	return err
}

func (r *RedisJwtHandler) CheckSession(ctx *gin.Context, ssid string) error {
	_, err := r.cmd.Exists(ctx, fmt.Sprintf("userd:ssid:%s", ssid)).Result()
	return err
}

func (r *RedisJwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	myclaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7 * 24)),
		},
		Ssid:   ssid,
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)

	//token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString(RtKey)
	if err != nil {
		//ctx.JSON(http.StatusOK, gin.H{
		//	"message": "jwt系统错误",
		//	"error":   err.Error(),
		//})
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}
