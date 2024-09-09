package web

import "github.com/gin-gonic/gin"

func RegisterRoutes(server *gin.Engine, u *UserHandler) {
	//u := NewUserHandler()
	ug := server.Group("/user")
	ug.POST("/login", u.Login)

	ug.POST("/edit", u.Edit)

	ug.POST("/signup", u.SignUp)
	ug.POST("/login/jwt", u.LoginJwt)
	ug.GET("/login/jwt", u.LoginJwt)
	ug.POST("/login/profile", u.Profile)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms/code/verify", u.VerifyLoginSMSCode)
	ug.POST("/logout", u.Logout)
	ug.POST("/login/refresh", u.RefreshToken)
}
