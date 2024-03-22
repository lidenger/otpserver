package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/result"
)

func AddServer(ctx *gin.Context) {
	var p *param.ServerParam
	p = validParam(ctx, p)
	if p == nil {
		return
	}
	err := service.ServerSvcIns.Add(ctx, p)
	result.R(ctx, err, "")
}

func GetServer(ctx *gin.Context) {
	sign := ctx.Param("sign")
	model, err := service.ServerSvcIns.GetBySign(ctx, sign)
	result.R(ctx, err, model)
}
