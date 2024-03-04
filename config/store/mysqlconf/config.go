package mysqlconf

import (
	"fmt"
	"github.com/lidenger/otpserver/config/server"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitMySQL(conf *server.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s",
		conf.MySQL.UserName,
		conf.MySQL.Password,
		conf.MySQL.Address,
		conf.MySQL.DbName,
		conf.MySQL.ConnMaxWaitTime)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         1024,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	if err != nil {
		panic(fmt.Sprintf("MySQL配置失败:%+v", err))
	}
	_db, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("MySQL初始化失败:%+v", err))
	}
	_db.SetConnMaxLifetime(time.Duration(conf.MySQL.ConnMaxLifeTime) * time.Minute)
	_db.SetMaxIdleConns(conf.MySQL.MaxIdleConn)
	_db.SetMaxOpenConns(conf.MySQL.MaxOpenConn)
	return db
}
