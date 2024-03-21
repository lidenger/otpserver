package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/internal/handler"
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
	g.GET("/health", handler.Health)

	v1 := g.Group("/v1")
	{
		secretv1 := v1.Group("/secret")
		{
			secretv1.POST("", handler.AddAccountSecret)        // POST /v1/secret
			secretv1.GET(":account", handler.GetAccountSecret) // GET /v1/secret/zhangsan
		}
	}

}
