package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/result"
)

// AddAccountSecret 新增账号密钥
func AddAccountSecret(ctx *gin.Context) {
	var p *param.SecretParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.R(ctx, otperr.ErrParamIllegal(errors.New("非法参数")), "")
		return
	}
	err := service.AccountSecretSvc.Add(ctx, p.Account)
	result.R(ctx, err, "")
}
