package otperr

import (
	"errors"
	"fmt"
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

// ErrRepeatAdd 重复添加
func ErrRepeatAdd(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusBadRequest,
		code:     code.ParamRepeatAdd,
	}
}

// ErrDataNotExists 重复添加
func ErrDataNotExists(format string, arg ...any) IErr {
	msg := fmt.Sprintf(format, arg)
	return &Err{
		error:    errors.New(msg),
		httpCode: http.StatusBadRequest,
		code:     code.DataNotExists,
	}
}

// ErrReqOverLimit 请求超限
func ErrReqOverLimit(format string, arg ...any) IErr {
	msg := fmt.Sprintf(format, arg)
	return &Err{
		error:    errors.New(msg),
		httpCode: http.StatusTooManyRequests,
		code:     code.ReqOverLimit,
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

// ErrStore 存储异常
func ErrStore(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.StoreErr,
	}
}

// ErrStoreBackup 备存储异常
func ErrStoreBackup(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.StoreBackupErr,
	}
}

// ErrEncrypt 加密异常
func ErrEncrypt(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.EncryptErr,
	}
}

// ErrDecrypt 解密异常
func ErrDecrypt(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.DecryptErr,
	}
}

// ErrAccountSecretDataCheck 账号密钥数据校验失败
func ErrAccountSecretDataCheck(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.AccountSecretDataCheckErr,
	}
}

// ErrGenCode 生成动态令牌异常
func ErrGenCode(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.GenCodeErr,
	}
}
