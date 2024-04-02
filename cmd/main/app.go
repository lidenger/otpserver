package main

import (
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/router"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/internal/store"
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
	// 加载配置文件
	conf := serverconf.GetSysConf()
	// 正常启动Http服务
	log.InitLog(conf)
	store.InitStore(conf)
	service.InitSvc()
	g := router.InitRouter(conf)
	log.Info("Http服务已启动,端口:%d", conf.Server.Port)
	err := g.Run(fmt.Sprintf("0.0.0.0:%d", conf.Server.Port))
	if err != nil {
		panic(err)
	}
}
