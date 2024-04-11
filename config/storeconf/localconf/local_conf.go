package localconf

import (
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/pkg/enum"
)

type LocalConf struct {
}

var LocalConfIns = &LocalConf{}

func Initialize(conf *config.M) {

}

func (l *LocalConf) GetStoreType() string {
	return enum.LocalStore
}

func (l *LocalConf) CloseStore() {

}

func (l *LocalConf) TestStore() error {

	return nil
}
