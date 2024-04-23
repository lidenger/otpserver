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

func wrapStrErr(param any) error {
	var err error
	if x, ok := param.(error); ok {
		err = x
	}
	if x, ok := param.(string); ok {
		err = errors.New(x)
	}
	return err
}

// ErrParamIllegal 非法参数
func ErrParamIllegal(param any) IErr {
	return &Err{
		error:    wrapStrErr(param),
		httpCode: http.StatusBadRequest,
		code:     code.ParamIllegal,
	}
}

// ErrRepeatAdd 重复添加
func ErrRepeatAdd(err any) IErr {
	return &Err{
		error:    wrapStrErr(err),
		httpCode: http.StatusBadRequest,
		code:     code.ParamRepeatAdd,
	}
}

// ErrDataNotExists 数据不存在
func ErrDataNotExists(format string, arg ...any) IErr {
	msg := fmt.Sprintf(format, arg)
	return &Err{
		error:    errors.New(msg),
		httpCode: http.StatusBadRequest,
		code:     code.DataNotExists,
	}
}

// ErrUnauthorized 权限不足
func ErrUnauthorized(format string, arg ...any) IErr {
	msg := format
	if len(arg) != 0 {
		msg = fmt.Sprintf(format, arg)
	}
	return &Err{
		error:    errors.New(msg),
		httpCode: http.StatusUnauthorized,
		code:     code.Unauthorized,
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
func ErrServerUnReady(param any) IErr {
	return &Err{
		error:    wrapStrErr(param),
		httpCode: http.StatusServiceUnavailable,
		code:     code.ServerUnready,
	}
}

// ErrServerFuncNonsupport 服务未支持的功能
func ErrServerFuncNonsupport(param any) IErr {
	return &Err{
		error:    wrapStrErr(param),
		httpCode: http.StatusNotImplemented,
		code:     code.ServerFuncNonsupport,
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

// ErrStoreUnready 备存储异常
func ErrStoreUnready(err error) IErr {
	return &Err{
		error:    err,
		httpCode: http.StatusInternalServerError,
		code:     code.StoreUnready,
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
