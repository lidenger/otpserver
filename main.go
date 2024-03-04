package main

import (
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/server"
)

func main() {
	server.InitConfig("dev")
	log.InitLog()
}
