package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/result"
	"math"
	"strconv"
	"time"
)

// GetAccessToken 获取AccessToken
// serverSign 服务标识，添加服务时指定的标识
// timeToken AES(KEY+IV, UNIX时间戳)
func GetAccessToken(c *gin.Context) {
	success, m := verifySign(c)
	if !success {
		return
	}
	if !VerifyTimeToken(c, m) {
		return
	}
	token, err := service.GenAccessToken(m.Sign)
	result.R(c, err, token)
}

// VerifyAccessToken 验证Token
func VerifyAccessToken(c *gin.Context) {
	accessToken, exists := c.GetQuery("accessToken")
	if !exists || len(accessToken) == 0 {
		result.ParamErr(c, "缺少accessToken参数")
		return
	}
	tokenM, err := service.AnalysisAccessToken(accessToken)
	if err != nil {
		result.ParamErr(c, "access token错误:"+err.Error())
		return
	}
	err = service.VerifyAccessTokenM(tokenM)
	if err != nil {
		result.ParamErr(c, err.Error())
		return
	}
	result.R(c, nil, "1")
}

// 验证服务sign
func verifySign(c *gin.Context) (bool, *model.ServerModel) {
	serverSign, exists := c.GetQuery("serverSign")
	if !exists {
		result.ParamErr(c, "缺失serverSign参数")
		return false, nil
	}
	s, err := service.ServerSvcIns.GetBySign(c, serverSign)
	if err != nil {
		result.R(c, err, "")
		return false, nil
	}
	if s == nil {
		result.ParamErr(c, "服务不存在:"+serverSign)
		return false, nil
	}
	if s.IsEnable == 2 {
		result.ParamErr(c, "server已禁用")
		return false, nil
	}
	return true, s
}

// VerifyTimeToken 验证客户端时间token
func VerifyTimeToken(c *gin.Context, m *model.ServerModel) bool {
	token, exists := c.GetQuery("timeToken")
	if !exists {
		result.ParamErr(c, "缺失timeToken参数")
		return false
	}
	key := []byte(m.Secret)
	iv := []byte(m.IV)
	t, err := crypt.Decrypt(key, iv, token)
	if err != nil {
		result.ParamErr(c, "timeToken不正确:"+err.Error())
		return false
	}
	clientTime, err := strconv.Atoi(string(t))
	if err != nil {
		result.ParamErr(c, "timeToken不正确:"+err.Error())
		return false
	}
	// 检测时间误差
	conf := serverconf.GetSysConf()
	validMinute := float64(conf.Server.TimeTokenValidMinute)
	if math.Abs(float64(int64(clientTime)-time.Now().Unix())) > (validMinute * 60) {
		msg := fmt.Sprintf("timeToken不正确,和服务端时间差大于%f分钟,client time:%d,server time:%d", validMinute, clientTime, time.Now().Unix())
		result.ParamErr(c, msg)
		return false
	}
	return true
}
