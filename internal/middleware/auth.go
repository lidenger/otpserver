package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
)

// ServerAuth 接入服务鉴权
func ServerAuth(c *gin.Context) {
	log.Info("ServerAuth%s", "request")
	c.Next()
}

// AdminAuth admin管理平台鉴权
func AdminAuth(c *gin.Context) {
	log.Info("AdminAuth%s", "request")
	c.Next()
}
