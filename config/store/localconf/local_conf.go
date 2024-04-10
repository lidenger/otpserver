package localconf

import (
	"github.com/lidenger/otpserver/config"
)

type LocalConf struct {
}

var LocalConfIns = &LocalConf{}

func Initialize(conf *config.M) {

}

func (l *LocalConf) CloseStore() {

}

func (l *LocalConf) TestStore() error {

	return nil
}
