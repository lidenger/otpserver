package timer

import (
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"time"
)

var storeHealthCheckTicker = time.NewTicker(5 * time.Second)

func StoreHealthCheckTickerStart() {
	go func() {
		for range storeHealthCheckTicker.C {
			// 检测MySQL
			mysqlconf.Test()
		}
	}()
}

func StoreHealthCheckTickerStop() {
	storeHealthCheckTicker.Stop()
	log.Info("storeHealthCheckTicker已关闭")
}
