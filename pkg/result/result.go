package result

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/code"
	"github.com/lidenger/otpserver/pkg/otperr"
	"net/http"
)

func result(ctx *gin.Context, httpCode int, code code.CODE, message string, data any) {
	ctx.JSON(httpCode, gin.H{
		"code": code,
		"msg":  message,
		"data": data,
	})
}

func Success(ctx *gin.Context, data any) {
	result(ctx, http.StatusOK, code.Success, "success", data)
}

func ParamErr(ctx *gin.Context, msg string) {
	err := otperr.ErrParamIllegal(msg)
	R(ctx, err, "")
}

func R(ctx *gin.Context, err error, data any) {
	if err == nil {
		Success(ctx, data)
		return
	}
	log.Error("%+v", err)
	if x, ok := err.(otperr.IErr); ok {
		result(ctx, x.GetHttpCode(), x.GetCode(), x.Error(), "")
	} else {
		result(ctx, http.StatusInternalServerError, code.UnknownErr, x.Error(), "")
	}
}

func MakePagingResult[T int | int64](rows any, total T) map[string]any {
	pagingResult := make(map[string]any, 2)
	pagingResult["rows"] = rows
	pagingResult["total"] = total
	return pagingResult
}
