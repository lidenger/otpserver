package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/result"
)

// ServerAuth 接入服务鉴权
func ServerAuth(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if len(accessToken) == 0 {
		result.R(c, otperr.ErrUnauthorized("Authorization不正确"), "")
		c.Abort()
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
	c.Next()
}

// AdminAuth admin管理平台鉴权
func AdminAuth(c *gin.Context) {
	log.Info("AdminAuth%s", "request")
	c.Next()
}
