package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/pkg/result"
)

func validParam[T any](ctx *gin.Context, p *T) *T {
	if err := ctx.ShouldBindJSON(&p); err != nil {
		result.ParamErr(ctx, "非法参数")
		return nil
	}
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		result.ParamErr(ctx, err.Error())
		return nil
	}
	return p
}
