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
	if cmd.P.Init {
		cmd.GenKeyFile()
		fmt.Println("系统初始化模式")
		return
	}
	crypt := cmd.AnalysisKeyFile(cmd.P.StartKeyFile)
	cmd.P.Crypt = crypt
	fmt.Printf("%+v\n", cmd.P)
	conf := serverconf.InitConfig()
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
