package memoryconf

import (
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/enum"
)

type MemoryConf struct {
}

var MemoryConfIns = &MemoryConf{}

func Initialize(conf *config.M) {

}

func (m *MemoryConf) GetStoreType() string {
	return enum.MemoryStore
}

func (m *MemoryConf) CloseStore() {
	log.Info("memory存储已关闭")
}

func (m *MemoryConf) TestStore() error {
	tmp := make([]int, 0)
	tmp = append(tmp, 1)
	return nil
}
