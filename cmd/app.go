package main

import (
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/router"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/internal/store"
)

func main() {
	serverconf.InitCmdParam()
	conf := serverconf.InitConfig()
	log.InitLog(conf)
	store.InitStore(conf)
	service.InitSvc()
	g := router.InitRouter()
	log.Info("Http服务已启动,端口:%d", conf.Server.Port)
	err := g.Run(fmt.Sprintf("0.0.0.0:%d", conf.Server.Port))
	if err != nil {
		panic(err)
	}
}
