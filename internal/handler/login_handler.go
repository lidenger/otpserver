package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/result"
	"time"
)

type loginParam struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	var p loginParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.ParamErr(ctx, "非法参数")
		return
	}
	if p.Account != "admin" || len(p.Password) == 0 {
		result.ParamErr(ctx, "账号密码不正确")
		return
	}
	pwdDigest := crypt.HmacDigest(cmd.P.RootKey256, p.Password)
	sysConf, err := service.ConfSvcIns.GetByKey(ctx, enum.AdminPasswordKey)
	if err != nil {
		result.R(ctx, otperr.ErrServer(err), "")
		return
	}
	if pwdDigest != sysConf.Val {
		result.ParamErr(ctx, "账号密码不正确")
		return
	}
	conf := serverconf.GetSysConf()
	maxAge := time.Duration(conf.Server.AdminLoginValidHour) * time.Hour
	token, err := service.GenLoginToken(p.Account)
	if err != nil {
		result.R(ctx, otperr.ErrServer(err), "")
		return
	}
	ctx.SetCookie("otp_login_token", token, int(maxAge.Seconds()), "/", conf.Server.Domain, false, false)

	result.Success(ctx, "success")
}

func Logout(ctx *gin.Context) {
	conf := serverconf.GetSysConf()
	ctx.SetCookie("otp_login_token", "invalid", 1, "/", conf.Server.Domain, false, false)
	result.Success(ctx, "success")
}
