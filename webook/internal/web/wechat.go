package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	svc  wechat.Service
	usvc service.UserService
	JwtHadler
	stateKey []byte
	cfg      Config
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
	//RegisteredClaims jwt.Claims
}

type Config struct {
	Secure bool
}

func NewOAuth2WechatHandler(svc wechat.Service, usvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:      svc,
		usvc:     usvc,
		stateKey: []byte("oN1)tV1{xA6#xM2/nR5/hU1#fH2$bU0$"),
		cfg: Config{
			Secure: true,
		},
		JwtHadler: NewJwtHadler(),
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", h.AuthURL)
	g.Any("/callback", h.Callback)
	//g.GET("/")

}

func (h *OAuth2WechatHandler) AuthURL(ctx *gin.Context) {
	state := uuid.New()
	url, err := h.svc.AuthURL(ctx, state)
	if err != nil || url == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造失败",
			Data: nil,
		})
		return
	}
	tokenStr, err := h.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
	}

	ctx.JSON(http.StatusOK, Result{
		Msg:  tokenStr,
		Data: url,
	})

}

func (h *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, StateClaims{
		State: state,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 3)),
		},
	})

	tokenStr, err := token.SignedString(h.stateKey)
	if err != nil {
		return "", err
	}
	ctx.SetCookie("jwt-state", tokenStr, 180,
		"/oauth2/wechat/callback", "", h.cfg.Secure, true)
	return tokenStr, nil
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")
	err := h.VerifyState(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "登录失败",
			Data: err.Error(),
		})
		return
	}

	//wechatInfo, err := h.svc.VerifyCode(ctx, code, state)
	//if err != nil {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg:  err.Error(),
	//		Data: nil,
	//	})
	//}
	wechatInfo := domain.WechatInfo{
		OpenID:  code,
		UnionID: state,
	}
	u, err := h.usvc.FindOrCreateWithWechat(ctx, wechatInfo)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: nil,
		})
		return
	}
	//设置token

	err = h.setLoginToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "成功",
		Data: nil,
	})
}

func (h *OAuth2WechatHandler) VerifyState(ctx *gin.Context, state string) error {
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到Cookie, %w", err)
	}
	var sc StateClaims
	c, err := jwt.ParseWithClaims(ck, &sc, func(token *jwt.Token) (interface{}, error) {
		return h.stateKey, nil
	})
	if err != nil || !c.Valid {
		return fmt.Errorf("toke无效")
	}
	if sc.State != state {
		return fmt.Errorf("state不相等")
	}
	return nil
}
