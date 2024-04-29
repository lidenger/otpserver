package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/pkg/result"
	"time"
)

func Health(c *gin.Context) {
	data := make(map[string]string, 2)
	data["status"] = "ok"
	data["time"] = time.Now().Format(time.DateTime)
	result.ParamErr(c, "test param error")
	//result.Success(c, data)
}
