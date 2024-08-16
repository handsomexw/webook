package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         *service.UserService
	emailRegexp *regexp.Regexp
	codeService *service.CodeService
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64  `json:"userid"`
	UserAgent string `json:"useragent"`
}

const (
	emailRegexPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

func NewUserHandler(svc *service.UserService, codeService *service.CodeService) *UserHandler {
	return &UserHandler{
		emailRegexp: regexp.MustCompile(emailRegexPattern, regexp.RegexOptions(regexp.Unicode)),
		svc:         svc,
		codeService: codeService,
	}
}

func (u *UserHandler) RegisterRoutes(engine *gin.Engine) {}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}

	var req SignUpReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return
	}

	if is, _ := u.emailRegexp.MatchString(req.Email); !is {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "邮箱格式错误",
		})
		return
	}

	err := u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrorUserDuplicate) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "邮箱冲突",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message":  "注册成功",
		"user":     req.Email,
		"confirm":  req.ConfirmPassword,
		"password": req.Password,
	})
	//fmt.Println(req)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrorInvalidUserOrPassword) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "用户名或密码不对",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
		return
	}
	//设置session
	//如何shezhijson呢
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		//Secure:   true,
		//HttpOnly: true,
		MaxAge: 5,
	})
	sess.Set("userId", user.Id)
	err = sess.Save()
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
	})

}

func (u *UserHandler) Edit(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"messgae": "这是edit",
	})
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	err := sess.Save()
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "退出登录成功",
	})
}

func (u *UserHandler) LoginJwt(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrorInvalidUserOrPassword) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "用户名或密码不对",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
		return
	}
	//设置token
	//claims := jwt.MapClaims{
	//	"username": "example_user",
	//	"admin":    true,
	//}
	//这是一种携带信息的方法
	//token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//	"userId": user.Id,
	//})
	err = u.setJwtToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
	})
}

func (u *UserHandler) setJwtToken(ctx *gin.Context, uid int64) error {
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

func (u *UserHandler) Profile(ctx *gin.Context) {
	c, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "读取claims错误",
		})
		return
	}
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "读取claims错误",
		})
		return
	}
	println(claims.UserId)
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Request struct {
		Phone string `json:"phone"`
	}
	const biz = "login"
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "验证码校验系统错误",
		})
		return
	}
	err := u.codeService.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码发送错误",
			Data: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "发送成功",
	})

}
func (u *UserHandler) VerifyLoginSMSCode(ctx *gin.Context) {
	type Request struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	const biz = "login"
	var req Request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "验证码校验系统错误",
		})
		fmt.Println("这是绑定阶段")
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "手机号格式错误",
		})
		return
	}
	ok, err := u.codeService.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 2,
			Msg:  "验证码校验系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "验证码错误",
		})
		return
	}
	//id 如何获取
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 2,
			Msg:  "创建用户错误",
			Data: err.Error(),
		})
	}
	err = u.setJwtToken(ctx, user.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 2,
			Msg:  "创建用户错误",
			Data: err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "验证码校验通过",
		Data: "",
	})
}
