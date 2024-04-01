package service

import (
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/store/mysqlstore"
	"github.com/lidenger/otpserver/internal/store/pgsqlstore"
)

var SecretSvcIns = &SecretSvc{}
var ServerSvcIns = &ServerSvc{}

func InitSvc() {

	switch cmd.P.MainStore {
	case "mysql":
		SecretSvcIns.Store = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.Store = &mysqlstore.ServerStore{DB: mysqlconf.DB}
	case "pgsql":
		SecretSvcIns.Store = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	switch cmd.P.BackupStore {
	case "mysql":
		SecretSvcIns.StoreBackup = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.StoreBackup = &mysqlstore.ServerStore{DB: mysqlconf.DB}
	case "pgsql":
		SecretSvcIns.StoreBackup = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}
