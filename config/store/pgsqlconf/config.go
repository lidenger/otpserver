package pgsqlconf

import (
	"github.com/lidenger/otpserver/config/serverconf"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitPgsql(conf *serverconf.Config) {

}
