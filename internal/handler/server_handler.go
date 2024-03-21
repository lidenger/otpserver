package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/result"
)

func AddServer(ctx *gin.Context) {
	var p *param.ServerParam
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.ParamErr(ctx, "非法参数")
		return
	}
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		result.ParamErr(ctx, err.Error())
		return
	}
	err = service.ServerSvcIns.Add(ctx, p)
	result.R(ctx, err, "")
}
