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
	model, err := service.ServerSvcIns.GetBySign(ctx, sign, true)
	result.R(ctx, err, model)
}

func PagingServer(ctx *gin.Context) {
	pagingParam := validPagingParam(ctx)
	if pagingParam == nil {
		return
	}
	p := &param.ServerPagingParam{}
	p.PageNo = pagingParam.PageNo
	p.PageSize = pagingParam.PageSize
	p.SearchTxt = pagingParam.SearchTxt
	isEnable := getIntParamByQuery(ctx, "isEnable")
	if isEnable == -1 {
		return
	}
	data, total, err := service.ServerSvcIns.Paging(ctx, p)
	result.R(ctx, err, result.MakePagingResult(data, total))
}

func SetServerEnable(ctx *gin.Context) {
	var p *param.ServerParam
	p = validParam(ctx, p)
	if p == nil {
		return
	}
	err := service.ServerSvcIns.SetEnable(ctx, p)
	result.R(ctx, err, "")
}

func EditBase(ctx *gin.Context) {
	var p *param.ServerParam
	p = validParam(ctx, p)
	if p == nil {
		return
	}
	err := service.ServerSvcIns.EditBase(ctx, p)
	result.R(ctx, err, "")
}
