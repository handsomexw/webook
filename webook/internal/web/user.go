package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserService
	emailRegexp *regexp.Regexp
	codeService service.CodeService
	jwtHadler
}

const (
	emailRegexPattern = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

func NewUserHandler(svc service.UserService, codeService service.CodeService) *UserHandler {
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
	if err := ctx.BindJSON(&req); err != nil {
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
	//ctx.String(200, "注册成功")
	//ctx.JSON(200, Result{
	//	Code: 0,
	//	Msg:  "注册成功",
	//})
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
	//更新更新不敏感信息可以随便更新
	//更新敏感信息如邮箱，手机号，需要用到验证码模块
	type Req struct {
		NewName  string `json:"newName"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Req
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
	}
	//校验规则
	switch len(req.NewName) {
	case 0:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "用户名不能为空",
		})
	case 1024:
		ctx.JSON(http.StatusOK, gin.H{
			"message": "用户名太长",
		})
	default:
	}
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "日期格式不对",
			Data: nil,
		})
	}
	uc := ctx.MustGet("claims").(UserClaims)
	err = u.svc.UpdateNoeSensitiveInfo(ctx, domain.User{
		Id:       uc.UserId,
		Name:     req.NewName,
		Birthday: birthday,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 6,
			Msg:  "系统错误",
			Data: nil,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "修改成功",
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

func (u *UserHandler) Profile(ctx *gin.Context) {
	//c, ok := ctx.Get("claims")
	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"message": "读取claims错误",
	//	})
	//	return
	//}
	//claims, ok := c.(*UserClaims)
	//if !ok {
	//	ctx.JSON(http.StatusOK, gin.H{
	//		"message": "读取claims错误",
	//	})
	//	return
	//}
	//println(claims.UserId)
	type ProfileReq struct {
		Email string `json:"email"`
	}
	sess := sessions.Default(ctx)
	id := sess.Get("userId").(int64)
	c, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "系统错误",
		})
	}
	ctx.JSON(http.StatusOK, ProfileReq{
		Email: c.Email,
	})
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Request struct {
		Phone string `json:"phone"`
	}
	const biz = "login"
	var req Request
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 2,
			Msg:  "验证系统错误",
			Data: nil,
		})
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 3,
			Msg:  "手机号输入错误",
			Data: nil,
		})
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
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "发送成功",
		Data: nil,
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
