package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"go.uber.org/zap"
	"time"
)

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		cost := time.Since(start)
		log.HttpZapLogger.Info(path,
			zap.Int("code", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("ip", c.ClientIP()),
			zap.String("ua", c.Request.UserAgent()),
			zap.Duration("cost", cost),
		)
	}
}
