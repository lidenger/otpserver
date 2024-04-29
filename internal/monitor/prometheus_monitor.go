package monitor

import (
	"github.com/gin-gonic/gin"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// 计算QPS： irate(otp_server_http_req_total{job="otpserver"}[1m])
// 统计4xx占比： sum(rate(otp_server_http_response_code_histogram_bucket{le="499"}[5m])) - sum(rate(otp_server_http_response_code_histogram_bucket{le="299"}[5m]))
// 查询所有otp server指标: {__name__=~"^otp_server.*"}
var (
	// http request
	httpReqTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "otp_server_http_req_total",
		Help: "The total number of http requests total(http请求总次数)",
	})
	httpReqLimitTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "otp_server_http_limit_req_total",
		Help: "The total number of http limit requests total(http请求超限总次数)",
	})
	httpReqCost = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "otp_server_http_req_cost_gauge",
		Help: "The total number of http requests cost time(http请求耗时，单位：毫秒)",
	})
	HttpCostHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "otp_server_http_cost_histogram",
		Help:    "otp server http cost histogram(http请求耗时分布，单位：毫秒)",
		Buckets: []float64{10, 100, 300, 500, 1000, 3000, 10000},
	})
	HttpRespCodeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "otp_server_http_response_code_histogram",
		Help:    "otp server http response status code histogram(http响应状态码)",
		Buckets: []float64{299, 499, 599},
	})

	// store 健康状态：0 未启用，1 正常，2 异常
	storeMySQL = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "otp_server_store_mysql_health_gauge",
		Help: "mysql store health status(mysql存储健康状态)",
	})
	storePgSQL = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "otp_server_store_pgsql_health_gauge",
		Help: "pgsql store health status(pgsql存储健康状态)",
	})
	storeLocal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "otp_server_store_local_health_gauge",
		Help: "local store health status(本地存储健康状态)",
	})
	storeMemory = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "otp_server_store_memory_health_gauge",
		Help: "memory store health status(内存存储健康状态)",
	})
)

// HttpReqTotal http请求次数
func HttpReqTotal() {
	go func() {
		httpReqTotal.Inc()
	}()
}

// HttpReqLimitTotal 被limit参数限制的http请求次数
func HttpReqLimitTotal() {
	go func() {
		httpReqLimitTotal.Inc()
	}()
}

// HttpReqCost http请求耗时
func HttpReqCost(val float64) {
	go func() {
		httpReqCost.Set(val)
		HttpCostHistogram.Observe(val)
	}()
}

func HttpRepsCode(code int) {
	go func() {
		HttpRespCodeHistogram.Observe(float64(code))
	}()
}

// store健康状态
func storeHealthStatus(storeType string, status uint8) {
	go func() {
		switch storeType {
		case enum.MySQLStore:
			storeMySQL.Set(float64(status))
		case enum.PostGreSQLStore:
			storePgSQL.Set(float64(status))
		case enum.LocalStore:
			storeLocal.Set(float64(status))
		case enum.MemoryStore:
			storeMemory.Set(float64(status))
		default:
			log.Error("未知的store类型:%s", storeType)
		}
	}()
}

// StoreStatusNotEnabled store未启用
func StoreStatusNotEnabled(storeType string) {
	storeHealthStatus(storeType, 0)
}

// StoreStatusRight store状态正常
func StoreStatusRight(storeType string) {
	storeHealthStatus(storeType, 1)
}

// StoreStatusError store状态异常
func StoreStatusError(storeType string) {
	storeHealthStatus(storeType, 2)
}

func InitPrometheusMonitor(g *gin.Engine) {
	prometheus.MustRegister(HttpCostHistogram)
	prometheus.MustRegister(HttpRespCodeHistogram)

	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	log.Info("Prometheus 初始化完成")
}
