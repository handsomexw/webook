package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type jwtHadler struct {
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64  `json:"userid"`
	UserAgent string `json:"useragent"`
}

func (u jwtHadler) setJwtToken(ctx *gin.Context, uid int64) error {
	myclaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		UserId:    uid,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myclaims)

	//token := jwt.New(jwt.SigningMethodHS512)
	tokenStr, err := token.SignedString([]byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"))
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
