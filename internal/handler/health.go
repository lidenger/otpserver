package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/pkg/result"
)

func Health(c *gin.Context) {
	result.Success(c, "up")
}
