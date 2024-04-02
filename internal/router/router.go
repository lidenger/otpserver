package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/handler"
	"github.com/lidenger/otpserver/internal/middleware"
)

func InitRouter() *gin.Engine {
	conf := serverconf.GetSysConf()
	gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(middleware.ReqLimit(conf.Server.ReqLimit))
	g.Use(middleware.RequestLog)
	api(g)
	return g
}

// RESTful API : GET 查询、获取，幂等操作; POST 新增; PUT 更新; DELETE 删除
func api(g *gin.Engine) {
	g.GET("/health", handler.Health)

	v1 := g.Group("/v1")
	{
		v1.POST("/access-token", handler.GenAccessToken)          // 生成access token
		v1.GET("/access-token/verify", handler.VerifyAccessToken) // 验证access token

		secretv1 := v1.Group("/secret") // 账号密钥
		secretv1.Use(middleware.ServerAuth)
		{
			secretv1.GET("/valid", handler.ValidCode)          // 验证动态令牌
			secretv1.POST("", handler.AddAccountSecret)        // 新增账号密钥 POST /v1/secret
			secretv1.GET(":account", handler.GetAccountSecret) // 获取指定账号密钥 GET /v1/secret/zhangsan
		}

		adminv1 := v1.Group("/admin") // 管理平台
		adminv1.Use(middleware.AdminAuth)
		{
			serverv1 := adminv1.Group("/server") // 接入服务
			{
				serverv1.POST("", handler.AddServer)     // 新增服务
				serverv1.GET(":sign", handler.GetServer) // 获取服务信息
			}
		}
	}

}
