package timer

import (
	"context"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/service"
	"time"
)

var localStoreCheckTicker = time.NewTicker(1 * time.Minute)

func StartLocalStoreCheckTicker() {
	go func() {
		for range localStoreCheckTicker.C {
			for _, ls := range service.FetchLocalStore() {
				err := ls.LoadAll(context.Background())
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}()
}

func StopLocalStoreCheckTicker() {
	localStoreCheckTicker.Stop()
	log.Info("localStoreCheckTicker已停止")
}
