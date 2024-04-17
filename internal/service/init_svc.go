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
var cacheStoreArr = make([]store.CacheStore, 0)

func Initialize(storeDetectionEventChan chan<- struct{}) {
	conf := serverconf.GetSysConf()

	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan
	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan

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

	if conf.Store.IsEnableLocal {
		memoryDependSecretStoreArr = append(memoryDependSecretStoreArr, localSecretStore)
		memoryDependServerStoreArr = append(memoryDependServerStoreArr, localServerStore)
	}

	if conf.Store.IsEnableMemory {
		secretMemory := &memorystore.SecretStore{
			Stores:                  memoryDependSecretStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		SecretSvcIns.StoreMemory = secretMemory
		cacheStoreArr = append(cacheStoreArr, secretMemory)

		serverMemory := &memorystore.ServerStore{
			Stores:                  memoryDependServerStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		ServerSvcIns.StoreMemory = serverMemory
		cacheStoreArr = append(cacheStoreArr, serverMemory)

		addSvcStore(enum.MemoryStore, SecretSvcIns.StoreMemory, ServerSvcIns.StoreMemory)
	}

	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}

func LoadAllCacheData() {
	for _, cache := range cacheStoreArr {
		err := cache.LoadAll(context.Background())
		if err != nil {
			log.Error("memory存储获取所有数据异常:%s", err.Error())
		}
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
