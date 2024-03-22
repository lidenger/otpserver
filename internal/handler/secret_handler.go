package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/result"
)

// AddAccountSecret 新增账号密钥
func AddAccountSecret(ctx *gin.Context) {
	var p *param.SecretParam
	p = validParam(ctx, p)
	if p == nil {
		return
	}
	err := service.SecretSvcIns.Add(ctx, p.Account)
	result.R(ctx, err, "")
}

// GetAccountSecret 获取密钥信息
func GetAccountSecret(ctx *gin.Context) {
	account := ctx.Param("account")
	model, err := service.SecretSvcIns.GetByAccount(ctx, account)
	result.R(ctx, err, model)
}
