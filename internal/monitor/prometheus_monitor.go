package monitor

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"time"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "otp_server_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func InitPrometheusMonitor(g *gin.Engine) {
	recordMetrics()
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Info("Prometheus 初始化完成")
}
