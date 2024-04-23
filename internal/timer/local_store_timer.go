package timer

import (
	"context"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/service"
	"time"
)

// 每4小时备份一次数据到本地
var localStoreCheckTicker *time.Ticker

func StartLocalStoreCheckTicker() {
	if localStoreCheckTicker == nil {
		conf := serverconf.GetSysConf()
		hour := time.Duration(conf.Store.CycleBakToLocalHour)
		localStoreCheckTicker = time.NewTicker(hour * time.Hour)
	}
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
