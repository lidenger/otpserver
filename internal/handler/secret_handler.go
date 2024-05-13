package handler

import (
	"fmt"
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
	err := service.SecretSvcIns.Add(ctx, p.Account, p.IsEnable)
	result.R(ctx, err, "")
}

// GetAccountSecret 获取密钥信息
func GetAccountSecret(ctx *gin.Context) {
	account := ctx.Param("account")
	model, err := service.SecretSvcIns.GetByAccount(ctx, account, true)
	result.R(ctx, err, model)
}

// GetQRCodeContent 获取添加密钥二维码内容
// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func GetQRCodeContent(ctx *gin.Context) {
	account := ctx.Param("account")
	model, err := service.SecretSvcIns.GetByAccount(ctx, account, true)
	if err != nil {
		result.R(ctx, err, "")
		return
	}
	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
	content := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", account, model.SecretSeed, "otpserver")
	result.R(ctx, nil, content)
}

func PagingAccountSecret(ctx *gin.Context) {
	pagingParam := validPagingParam(ctx)
	if pagingParam == nil {
		return
	}
	p := &param.SecretPagingParam{}
	p.PageNo = pagingParam.PageNo
	p.PageSize = pagingParam.PageSize
	p.SearchTxt = pagingParam.SearchTxt
	isEnable := getIntParamByQuery(ctx, "isEnable")
	if isEnable == -1 {
		return
	}
	data, total, err := service.SecretSvcIns.Paging(ctx, p)
	result.R(ctx, err, result.MakePagingResult(data, total))
}

func SetSecretEnable(ctx *gin.Context) {
	var p *param.SecretParam
	p = validParam(ctx, p)
	if p == nil {
		return
	}
	err := service.SecretSvcIns.SetEnable(ctx, p.Account, p.IsEnable)
	result.R(ctx, err, "")
}
