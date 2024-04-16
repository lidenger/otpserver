package service

import (
	"context"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/storeconf/mysqlconf"
	"github.com/lidenger/otpserver/config/storeconf/pgsqlconf"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/localstore"
	"github.com/lidenger/otpserver/internal/store/memorystore"
	"github.com/lidenger/otpserver/internal/store/mysqlstore"
	"github.com/lidenger/otpserver/internal/store/pgsqlstore"
	"github.com/lidenger/otpserver/pkg/enum"
)

var SecretSvcIns = &SecretSvc{}
var ServerSvcIns = &ServerSvc{}

// store type | healthFunc
var svcStoreStatusMap = make(map[string][]store.HealthFunc)

func Initialize() {
	conf := serverconf.GetSysConf()

	mysqlSecretStore := &mysqlstore.SecretStore{DB: mysqlconf.DB}
	mysqlServerStore := &mysqlstore.ServerStore{DB: mysqlconf.DB}

	pgsqlSecretStore := &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	pgsqlServerStore := &pgsqlstore.ServerStore{DB: pgsqlconf.DB}

	localSecretStore := &localstore.SecretStore{}
	localServerStore := &localstore.ServerStore{}

	memoryDependSecretStoreArr := make([]store.SecretStore, 0)
	memoryDependServerStoreArr := make([]store.ServerStore, 0)

	switch conf.Store.MainStore {
	case enum.MySQLStore:
		SecretSvcIns.Store = mysqlSecretStore
		ServerSvcIns.Store = mysqlServerStore
		addSvcStore(enum.MySQLStore, SecretSvcIns.Store, ServerSvcIns.Store)
	case enum.PostGreSQLStore:
		SecretSvcIns.Store = pgsqlSecretStore
		ServerSvcIns.Store = pgsqlServerStore
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.Store, ServerSvcIns.Store)
	}

	switch conf.Store.BackupStore {
	case enum.MySQLStore:
		SecretSvcIns.StoreBackup = mysqlSecretStore
		ServerSvcIns.StoreBackup = mysqlServerStore
		addSvcStore(enum.MySQLStore, SecretSvcIns.StoreBackup, ServerSvcIns.Store)
	case enum.PostGreSQLStore:
		SecretSvcIns.StoreBackup = pgsqlSecretStore
		ServerSvcIns.StoreBackup = pgsqlServerStore
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.StoreBackup, ServerSvcIns.StoreBackup)
	}

	memoryDependSecretStoreArr = append(memoryDependSecretStoreArr, SecretSvcIns.Store, SecretSvcIns.StoreBackup)
	memoryDependServerStoreArr = append(memoryDependServerStoreArr, ServerSvcIns.Store, ServerSvcIns.StoreBackup)

	if conf.Server.IsEnableLocalStore {
		memoryDependSecretStoreArr = append(memoryDependSecretStoreArr, localSecretStore)
		memoryDependServerStoreArr = append(memoryDependServerStoreArr, localServerStore)
	}

	if conf.Server.IsEnableMemoryStore {
		secretMemory := &memorystore.SecretStore{Stores: memoryDependSecretStoreArr}
		SecretSvcIns.StoreMemory = secretMemory
		cacheLoadData(secretMemory)

		serverMemory := &memorystore.ServerStore{Stores: memoryDependServerStoreArr}
		ServerSvcIns.StoreMemory = serverMemory
		cacheLoadData(serverMemory)

		addSvcStore(enum.MemoryStore, SecretSvcIns.StoreMemory, ServerSvcIns.StoreMemory)
	}

	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}

func cacheLoadData(cache store.CacheStore) {
	err := cache.LoadAll(context.Background())
	if err != nil {
		log.Error("memory存储获取所有数据异常:%s", err.Error())
	}
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
