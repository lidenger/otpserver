package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/result"
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
	if len(p.Account) == 0 || len(p.Password) == 0 {
		result.ParamErr(ctx, "账号密码不正确")
		return
	}
	log.Info("%+v", p)
	result.Success(ctx, "success")
}

func Logout(c *gin.Context) {
	result.Success(c, "success")
}
