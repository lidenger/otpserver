package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/result"
)

func validServer(ctx *gin.Context, sign string) bool {
	if len(sign) == 0 {
		result.ParamErr(ctx, "服务标识不能为空")
		return false
	}
	serverM, err := service.ServerSvcIns.GetBySign(ctx, sign, true)
	if err != nil {
		result.R(ctx, err, "")
		return false
	}
	if serverM == nil {
		result.ParamErr(ctx, "服务标识不正确")
		return false
	}
	return true
}

func GetServerIpList(ctx *gin.Context) {
	sign := ctx.Param("sign")
	if !validServer(ctx, sign) {
		return
	}
	data, err := service.ServerIpListSvcIns.GetBySign(ctx, sign)
	result.R(ctx, err, data)
}

func validSignAndIpParam(ctx *gin.Context) *param.ServerIpListParam {
	var p *param.ServerIpListParam
	p = validParam(ctx, p)
	if p == nil {
		return nil
	}
	if !validServer(ctx, p.Sign) {
		return nil
	} else {
		return p
	}
}

func RemoveServerIpList(ctx *gin.Context) {
	p := validSignAndIpParam(ctx)
	if validSignAndIpParam(ctx) == nil {
		return
	}
	err := service.ServerIpListSvcIns.Remove(ctx, p)
	result.R(ctx, err, "")
}

func AddServerIpList(ctx *gin.Context) {
	p := validSignAndIpParam(ctx)
	if validSignAndIpParam(ctx) == nil {
		return
	}
	err := service.ServerIpListSvcIns.Add(ctx, p)
	result.R(ctx, err, "")
}
