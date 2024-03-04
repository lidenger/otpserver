package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/controller"
	"github.com/lidenger/otpserver/internal/middleware"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(middleware.RequestLog())
	api(g)
	return g
}

func api(g *gin.Engine) {
	g.GET("/health", controller.Health)
}
