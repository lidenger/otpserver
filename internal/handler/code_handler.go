package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/result"
)

func ValidCode(ctx *gin.Context) {
	account, exists := ctx.GetQuery("account")
	if !exists {
		result.ParamErr(ctx, "缺失account参数")
		return
	}
	code, exists := ctx.GetQuery("code")
	if !exists {
		result.ParamErr(ctx, "缺失code参数")
		return
	}
	if len(code) != 6 {
		result.ParamErr(ctx, "code不是6位数")
		return
	}
	valid, err := service.ValidCode(ctx, account, code)
	result.R(ctx, err, valid)
}
