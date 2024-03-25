package main

import (
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config"
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
	conf := loadConfig()
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

func loadConfig() *config.M {
	var conf *config.M
	confFile := cmd.P.ConfFile
	switch cmd.P.ConfSource {
	case "nacos":
		if len(confFile) == 0 {
			confFile = "nacos.toml"
		}
		conf = cmd.LoadConfByNacos(confFile)
		fmt.Println("从nacos加载配置完成")
	case "local":
		if len(confFile) == 0 {
			confFile = "app.toml"
		}
		conf = cmd.LoadConfByLocalFile(confFile)
		fmt.Println("从本地文件加载配置完成")
	case "default":
		conf = serverconf.InitConfig()
		fmt.Println("从系统默认配置文件加载配置完成")
	default:
		panic(fmt.Sprintf("未知的配置来源:%s", cmd.P.ConfSource))
	}
	return conf
}
