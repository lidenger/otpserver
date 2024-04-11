package timer

import (
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/storeconf"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/internal/store"
	"time"
)

var storeHealthCheckTicker = time.NewTicker(5 * time.Second)

// StoreHealthCheckTickerStart 定期检测store健康状态
func StoreHealthCheckTickerStart() {
	go func() {
		for range storeHealthCheckTicker.C {
			storeArr := store.GetAllActiveStore()
			for _, s := range storeArr {
				testStore(s)
			}
		}
	}()
}

// 检测单个store健康状态
func testStore(s storeconf.Status) {
	err := s.TestStore()
	if err != nil {
		// TODO 发送存储异常报警消息
	}
	svcStores := service.GetSvcStores(s.GetStoreType())
	for _, svcStore := range svcStores {
		svcStore.SetStoreErr(err)
	}
}

func StoreHealthCheckTickerStop() {
	storeHealthCheckTicker.Stop()
	log.Info("storeHealthCheckTicker已关闭")
}
