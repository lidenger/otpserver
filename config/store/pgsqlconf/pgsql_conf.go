package pgsqlconf

import (
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/config/log"
	"gorm.io/gorm"
)

var DB *gorm.DB

type PgSQLConf struct {
}

var PgSQLConfIns = &PgSQLConf{}

func Initialize(conf *config.M) *gorm.DB {

	return nil
}

func (p *PgSQLConf) CloseStore() {
	if DB == nil {
		return
	}
	_db, _ := DB.DB()
	_ = _db.Close()
	log.Info("PostgreSQL已关闭")
}

func (p *PgSQLConf) TestStore() error {
	var db = DB
	var x uint8
	db = db.Raw("select 1").Scan(&x)
	return db.Error
}
