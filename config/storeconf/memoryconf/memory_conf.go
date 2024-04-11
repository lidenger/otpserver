package memoryconf

import (
	"github.com/lidenger/otpserver/config"
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

}

func (m *MemoryConf) TestStore() error {

	return nil
}
