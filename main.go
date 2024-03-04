package main

import (
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/server"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/internal/router"
)

func main() {
	conf := server.InitConfig("dev")
	log.InitLog(conf)
	mysqlconf.InitMySQL(conf)
	g := router.InitRouter()
	err := g.Run(fmt.Sprintf("0.0.0.0:%d", conf.Server.Port))
	if err != nil {
		panic(err)
	}
}
