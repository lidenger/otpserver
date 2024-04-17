package timer

import (
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/storeconf"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/internal/store"
	"time"
)

var storeHealthCheckTicker = time.NewTicker(10 * time.Second)

// StartStoreHealthCheckTicker 定期检测store健康状态
func StartStoreHealthCheckTicker(detectionChan <-chan struct{}) {
	// 开启定期检测ticker
	go func() {
		for range storeHealthCheckTicker.C {
			TriggerDetectionStore()
		}
	}()
	// 监听检测信号
	go func() {
		for {
			<-detectionChan
			TriggerDetectionStore()
		}
	}()
}

// TriggerDetectionStore 触发一次store检测
func TriggerDetectionStore() {
	storeArr := store.GetAllActiveStore()
	for _, s := range storeArr {
		testStore(s)
	}
}

// 检测store健康状态
func testStore(s storeconf.Status) {
	err := s.TestStore()
	if err != nil {
		log.Error(s.GetStoreType()+"存储检测异常:%s", err.Error())
		// TODO 发送存储异常报警消息
	}
	svcStores := service.GetSvcStores(s.GetStoreType())
	for _, svcStore := range svcStores {
		svcStore.SetStoreErr(err)
	}
}

func StopStoreHealthCheckTicker() {
	storeHealthCheckTicker.Stop()
	log.Info("storeHealthCheckTicker已停止")
}
