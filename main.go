package main

import (
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/server"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
)

func main() {
	conf := server.InitConfig("dev")
	log.InitLog(conf)
	mysqlconf.InitMySQL(conf)

}
