package main

import (
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/internal/router"
)

func main() {
	conf := serverconf.InitConfig("dev")
	log.InitLog(conf)
	initStore(conf)
	g := router.InitRouter()
	err := g.Run(fmt.Sprintf("0.0.0.0:%d", conf.Server.Port))
	if err != nil {
		panic(err)
	}
}

// 初始化store
func initStore(conf *serverconf.Config) {
	cmd := serverconf.GetCmdParam()
	if cmd.MainStore == "mysql" || cmd.BackupStore == "mysql" {
		mysqlconf.InitMySQL(conf)
	}
}
