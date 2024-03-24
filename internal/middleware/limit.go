package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/result"
	"golang.org/x/time/rate"
)

// ReqLimit 令牌桶算法
func ReqLimit(limit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)
	return func(c *gin.Context) {
		if limiter.Allow() {
			c.Next()
		} else {
			err := otperr.ErrReqOverLimit("请求超限:%d", limit)
			result.R(c, err, "")
			c.Abort()
		}
	}
}
