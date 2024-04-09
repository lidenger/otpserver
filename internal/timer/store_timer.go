package timer

import (
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"time"
)

var storeHealthCheckTicker = time.NewTicker(5 * time.Second)

func StoreHealthCheckTickerStart() {
	go func() {
		for range storeHealthCheckTicker.C {
			fmt.Println("storeHealthCheckTicker..")
		}
	}()
}

func StoreHealthCheckTickerStop() {
	storeHealthCheckTicker.Stop()
	log.Info("storeHealthCheckTicker已关闭")
}
