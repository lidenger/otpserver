package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/result"
	"math"
	"strconv"
	"time"
)

// GetAccessToken 获取AccessToken
// sign 服务标识，添加服务时指定的标识
// token AES(KEY+IV, UNIX时间戳)
func GetAccessToken(c *gin.Context) {
	serverSign, exists := c.GetQuery("sign")
	if !exists {
		result.ParamErr(c, "缺失sign参数")
		return
	}
	s, err := service.ServerSvcIns.GetBySign(c, serverSign)
	if err != nil {
		result.R(c, err, "")
		return
	}
	if s == nil {
		result.ParamErr(c, "服务不存在:"+serverSign)
		return
	}
	if s.IsEnable == 2 {
		result.ParamErr(c, "server已禁用")
		return
	}
	token, exists := c.GetQuery("token")
	if !exists {
		result.ParamErr(c, "缺失token参数")
		return
	}
	rootKey := service.ServerSvcIns.RootKey
	iv := service.ServerSvcIns.IV
	t, err := crypt.Decrypt(rootKey, iv, token)
	if err != nil {
		result.ParamErr(c, "token不正确:"+err.Error())
		return
	}
	clientTime, err := strconv.Atoi(string(t))
	if err != nil {
		result.ParamErr(c, "token不正确:"+err.Error())
		return
	}
	// 时间误差在60秒之内
	if math.Abs(float64(int64(clientTime)-time.Now().Unix())) > 60 {
		result.ParamErr(c, "token不正确:"+err.Error())
		return
	}

}
