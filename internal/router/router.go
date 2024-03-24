package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/handler"
	"github.com/lidenger/otpserver/internal/middleware"
)

func InitRouter(conf *serverconf.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(middleware.ReqLimit(conf.Server.ReqLimit))
	g.Use(middleware.RequestLog)
	api(g)
	return g
}

func api(g *gin.Engine) {
	g.GET("/health", handler.Health)

	v1 := g.Group("/v1")
	{
		{
			v1.GET("/valid", handler.ValidCode) // 验证动态令牌
		}

		secretv1 := v1.Group("/secret") // 账号密钥
		{
			secretv1.POST("", handler.AddAccountSecret)        // POST /v1/secret
			secretv1.GET(":account", handler.GetAccountSecret) // GET /v1/secret/zhangsan
		}

		serverv1 := v1.Group("/server") // 接入服务
		{
			serverv1.POST("", handler.AddServer)
			serverv1.GET(":sign", handler.GetServer)
		}
	}

}
