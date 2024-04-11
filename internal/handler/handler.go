package handler

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/pkg/result"
	"strconv"
	"strings"
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

func validPagingParam(ctx *gin.Context) *param.PagingParam {
	p := &param.PagingParam{}
	p.PageNo = getIntParamByQuery(ctx, "pageNo")
	if p.PageNo == -1 {
		return nil
	}
	p.PageSize = getIntParamByQuery(ctx, "pageSize")
	if p.PageSize == -1 {
		return nil
	}
	if p.PageNo < 1 {
		p.PageNo = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
	searchTxt, exists := ctx.GetQuery("searchTxt")
	if !exists {
		p.SearchTxt = ""
	} else {
		p.SearchTxt = strings.TrimSpace(searchTxt)
	}
	return p
}

// 在query中获取int参数, 0为不存在，-1为参数不合法
func getIntParamByQuery(ctx *gin.Context, key string) int {
	valStr, exists := ctx.GetQuery(key)
	if !exists {
		return 0
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		result.ParamErr(ctx, key+"参数不合法:"+err.Error())
		return -1
	}
	return val
}
