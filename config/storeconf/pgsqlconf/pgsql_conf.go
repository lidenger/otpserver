package pgsqlconf

import (
	"fmt"
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/enum"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var DB *gorm.DB

type PgSQLConf struct {
}

var PgSQLConfIns = &PgSQLConf{}

// Initialize https://github.com/go-gorm/postgres
func Initialize(conf *config.M) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		conf.PgSQL.Host,
		conf.PgSQL.UserName,
		conf.PgSQL.Password,
		conf.PgSQL.DbName,
		conf.PgSQL.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})

	if err != nil {
		panic(fmt.Sprintf("PostgreSQL配置失败:%+v", err))
	}
	_db, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("PostgreSQL初始化失败:%+v", err))
	}
	_db.SetConnMaxLifetime(time.Duration(conf.MySQL.ConnMaxLifeTime) * time.Minute)
	_db.SetMaxIdleConns(conf.MySQL.MaxIdleConn)
	_db.SetMaxOpenConns(conf.MySQL.MaxOpenConn)
	DB = db
	return db
}

func (p *PgSQLConf) GetStoreType() string {
	return enum.PostGreSQLStore
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
