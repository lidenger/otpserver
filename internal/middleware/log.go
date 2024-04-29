package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/monitor"
	"go.uber.org/zap"
	"time"
)

func RequestLog(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	c.Next()
	cost := time.Since(start)
	monitor.HttpReqCost(float64(cost.Milliseconds()))
	code := c.Writer.Status()
	monitor.HttpRepsCode(code)
	log.HttpZapLogger.Info(path,
		zap.Int("code", code),
		zap.String("method", c.Request.Method),
		zap.String("path", path),
		zap.String("ip", c.ClientIP()),
		zap.String("ua", c.Request.UserAgent()),
		zap.Duration("cost", cost),
	)
}
