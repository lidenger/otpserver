package pgsqlconf

import (
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/config/log"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Initialize(conf *config.M) *gorm.DB {

	return nil
}

func Close() {
	if DB == nil {
		return
	}
	_db, _ := DB.DB()
	_ = _db.Close()
	log.Info("PostgreSQL已关闭")
}

func Test() error {
	var db = DB
	var x uint8
	db = db.Raw("select 1").Scan(&x)
	return db.Error
}
