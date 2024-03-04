package otperr

import (
	"github.com/lidenger/otpserver/pkg/code"
	"net/http"
)

type IErr interface {
	GetHttpCode() int
	GetCode() code.CODE
	Error() string
}

type Err struct {
	error
	httpCode int
	code     code.CODE
}

func (e *Err) GetHttpCode() int {
	return e.httpCode
}

func (e *Err) GetCode() code.CODE {
	return e.code
}

// ErrParamIllegal 非法参数
func ErrParamIllegal(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusBadRequest,
		code:     code.ParamIllegal,
	}
}

// ErrServerUnReady 服务未准备就绪
func ErrServerUnReady(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusServiceUnavailable,
		code:     code.ServerUnready,
	}
}
