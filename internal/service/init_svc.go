package service

import (
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/mysqlstore"
	"github.com/lidenger/otpserver/internal/store/pgsqlstore"
	"github.com/lidenger/otpserver/pkg/enum"
)

var SecretSvcIns = &SecretSvc{}
var ServerSvcIns = &ServerSvc{}

// store type | healthFunc
var svcStoreStatusMap = make(map[string][]store.HealthFunc)

func InitSvc() {
	switch cmd.P.MainStore {
	case enum.MySQLStore:
		SecretSvcIns.Store = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.Store = &mysqlstore.ServerStore{DB: mysqlconf.DB}
		addSvcStore(enum.MySQLStore, SecretSvcIns.Store, ServerSvcIns.Store)
	case enum.PostGreSQLStore:
		SecretSvcIns.Store = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.Store)
	}
	switch cmd.P.BackupStore {
	case enum.MySQLStore:
		SecretSvcIns.StoreBackup = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.StoreBackup = &mysqlstore.ServerStore{DB: mysqlconf.DB}
		addSvcStore(enum.MySQLStore, SecretSvcIns.StoreBackup, ServerSvcIns.StoreBackup)
	case enum.PostGreSQLStore:
		SecretSvcIns.StoreBackup = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.StoreBackup)
	}
	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}

func addSvcStore(typ string, funcs ...store.HealthFunc) {
	fs := svcStoreStatusMap[typ]
	if fs == nil {
		fs = make([]store.HealthFunc, 0)
	}
	fs = append(fs, funcs...)
	svcStoreStatusMap[typ] = fs
}

// GetSvcStores 通过类型svc store
func GetSvcStores(typ string) []store.HealthFunc {
	fs := svcStoreStatusMap[typ]
	return fs
}
