package main

import (
	"context"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/store/localconf"
	"github.com/lidenger/otpserver/config/store/memoryconf"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/router"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/timer"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cmd.InitParam()
	// 系统初始化模式，生成[app.key]
	if cmd.P.IsInitMode {
		fmt.Println("系统初始化模式")
		cmd.InitMode()
		return
	}
	// 加载解析[app.key]
	crypt := cmd.AnalysisKeyFile(cmd.P.AppKeyFile)
	cmd.P.Crypt = crypt
	// 工具模式
	if cmd.P.IsToolMode {
		fmt.Println("工具模式")
		cmd.ToolMode()
		return
	}
	// 正常启动Http服务
	serverconf.Initialize()
	log.Initialize()
	store.Initialize()
	service.Initialize()
	g := router.Initialize()
	httpPort := serverconf.GetSysConf().Server.Port
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", httpPort),
		Handler: g,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Error("http服务异常:%+v", err)
		}
	}()
	log.Info("Http服务已启动,端口:%d", httpPort)
	// 启动store定期检测
	timer.StoreHealthCheckTickerStart()
	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 最多等待http shutdown 10秒
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("http服务关闭异常,err:%+v", err)
	} else {
		log.Info("http服务成功关闭,端口:%d", httpPort)
	}
	// 关闭其他资源
	closeRes()
	log.Info("所有资源已关闭")
}

func closeRes() {
	// 关闭store定期检测
	timer.StoreHealthCheckTickerStop()
	// 关闭MySQL
	mysqlconf.Close()
	// 关闭PostgreSQL
	pgsqlconf.Close()
	// 关闭本地存储
	localconf.Close()
	// 关闭memory存储
	memoryconf.Close()
}
