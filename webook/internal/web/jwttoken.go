package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid"
	"strings"

	"time"
)

type JwtHadler struct {
	//access_token
	atKey []byte
	//refresh_token
	rtKey []byte
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

func NewJwtHadler() JwtHadler {
	return JwtHadler{
		atKey: []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"),
		rtKey: []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU1$"),
	}
}

func (u JwtHadler) setJwtToken(ctx *gin.Context, uid int64, ssid string) error {
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
	tokenStr, err := token.SignedString(u.atKey)
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

// 同时刷新长短token，用redis判断是否有效
func (u JwtHadler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	myclaims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7 * 24)),
		},
		Ssid:   ssid,
		UserId: uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)

	//token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString(u.rtKey)
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

func ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Token")
	msgs := strings.Split(tokenHeader, " ")
	if len(msgs) != 1 {
		return " "
	}
	return msgs[0]
}

func (u JwtHadler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New()
	err := u.setJwtToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = u.setRefreshToken(ctx, uid, ssid)
	return err
}
