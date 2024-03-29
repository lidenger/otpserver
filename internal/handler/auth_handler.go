package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/model"
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
	success, m := verifySign(c)
	if !success {
		return
	}
	if !verifyToken(c, m) {
		return
	}

}

// 验证服务sign
func verifySign(c *gin.Context) (bool, *model.ServerModel) {
	serverSign, exists := c.GetQuery("sign")
	if !exists {
		result.ParamErr(c, "缺失sign参数")
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

// 验证token
func verifyToken(c *gin.Context, m *model.ServerModel) bool {
	token, exists := c.GetQuery("token")
	if !exists {
		result.ParamErr(c, "缺失token参数")
		return false
	}
	key := []byte(m.Secret)
	iv := []byte(m.IV)
	t, err := crypt.Decrypt(key, iv, token)
	if err != nil {
		result.ParamErr(c, "token不正确:"+err.Error())
		return false
	}
	clientTime, err := strconv.Atoi(string(t))
	if err != nil {
		result.ParamErr(c, "token不正确:"+err.Error())
		return false
	}
	// 时间误差在60秒之内
	if math.Abs(float64(int64(clientTime)-time.Now().Unix())) > 60 {
		result.ParamErr(c, "token不正确:"+err.Error())
		return false
	}
	return true
}
