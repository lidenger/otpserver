package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/result"
	"time"
)

// ServerAuth 接入服务鉴权
func ServerAuth(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) == 0 {
		result.R(c, otperr.ErrUnauthorized("Authorization不正确"), "")
		c.Abort()
		return
	}
	// 缓存中存在，验证通过
	tokenM := service.GetAccessTokenInCache(accessToken)
	if tokenM != nil {
		c.Next()
		return
	}
	// 解析access token
	m, err := service.AnalysisAccessToken(accessToken)
	if err != nil {
		result.R(c, otperr.ErrUnauthorized("Authorization不正确"), "")
		c.Abort()
		return
	}
	// 验证access token
	err = service.VerifyAccessTokenM(m)
	if err != nil {
		result.R(c, otperr.ErrUnauthorized("Authorization不正确"), "")
		c.Abort()
		return
	}
	// 验证生效记录到缓存中
	service.AddAccessTokenToCache(accessToken, m)
	c.Next()
}

// AdminAuth admin管理平台鉴权
func AdminAuth(c *gin.Context) {
	token, err := c.Cookie("otp_login_token")
	if err != nil || len(token) == 0 {
		log.Warn("获取otp_login_token失败,err:%v", err)
		result.R(c, otperr.ErrUnauthorized("未获取到token"), "")
		c.Abort()
		return
	}
	conf := serverconf.GetSysConf()
	maxValidTime := time.Duration(conf.Server.AdminLoginValidHour) * time.Hour
	_, err = service.AnalysisLoginToken(token, maxValidTime.Milliseconds())
	if err != nil {
		log.Warn("admin login token验证失败,err:%v", err)
		result.R(c, otperr.ErrUnauthorized("验证token失败"), "")
		c.Abort()
		return
	}
	c.Next()
}
