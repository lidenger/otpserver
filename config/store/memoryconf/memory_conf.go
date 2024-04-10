package memoryconf

import "github.com/lidenger/otpserver/config"

type MemoryConf struct {
}

var MemoryConfIns = &MemoryConf{}

func Initialize(conf *config.M) {

}

func (m *MemoryConf) CloseStore() {

}

func (m *MemoryConf) TestStore() error {

	return nil
}
