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
var loadAllStoreArr = make([]store.LoadAllFunc, 0)

func Initialize(storeDetectionEventChan chan<- struct{}) {
	conf := serverconf.GetSysConf()

	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan
	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan

	mysqlSecretStore := &mysqlstore.SecretStore{DB: mysqlconf.DB}
	mysqlServerStore := &mysqlstore.ServerStore{DB: mysqlconf.DB}

	pgsqlSecretStore := &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	pgsqlServerStore := &pgsqlstore.ServerStore{DB: pgsqlconf.DB}

	localStoreRootPath := conf.Server.RootPath + conf.Store.RootPath
	localSecretStore := &localstore.SecretStore{RootPath: localStoreRootPath}
	localServerStore := &localstore.ServerStore{RootPath: localStoreRootPath}

	memoryDependSecretStoreArr := make([]store.SecretStore, 0)
	memoryDependServerStoreArr := make([]store.ServerStore, 0)

	switch conf.Store.MainStore {
	case enum.MySQLStore:
		SecretSvcIns.Store = mysqlSecretStore
		ServerSvcIns.Store = mysqlServerStore
		localSecretStore.Store = mysqlSecretStore
		localServerStore.Store = mysqlServerStore
		addSvcStore(enum.MySQLStore, SecretSvcIns.Store, ServerSvcIns.Store)
	case enum.PostGreSQLStore:
		SecretSvcIns.Store = pgsqlSecretStore
		ServerSvcIns.Store = pgsqlServerStore
		localSecretStore.Store = pgsqlSecretStore
		localServerStore.Store = pgsqlServerStore
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

		loadAllStoreArr = append(loadAllStoreArr, localSecretStore, localServerStore)
	}

	if conf.Store.IsEnableMemory {
		secretMemory := &memorystore.SecretStore{
			Stores:                  memoryDependSecretStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		SecretSvcIns.StoreMemory = secretMemory

		serverMemory := &memorystore.ServerStore{
			Stores:                  memoryDependServerStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		ServerSvcIns.StoreMemory = serverMemory

		loadAllStoreArr = append(loadAllStoreArr, secretMemory, serverMemory)

		addSvcStore(enum.MemoryStore, SecretSvcIns.StoreMemory, ServerSvcIns.StoreMemory)
	}

	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}

func LoadAllData() {
	for _, cache := range loadAllStoreArr {
		_ = cache.LoadAll(context.Background())
	}
}

func FetchLocalStore() []store.LoadAllFunc {
	ss := make([]store.LoadAllFunc, 0)
	for _, s := range loadAllStoreArr {
		if s.GetStoreType() == enum.LocalStore {
			ss = append(ss, s)
		}
	}
	return ss
}

func addSvcStore(typ string, funcs ...store.HealthFunc) {
	fs := svcStoreStatusMap[typ]
	if fs == nil {
		fs = make([]store.HealthFunc, 0)
	}
	fs = append(fs, funcs...)
	svcStoreStatusMap[typ] = fs
}

// FetchSvcStores 通过类型获取svc store
func FetchSvcStores(typ string) []store.HealthFunc {
	fs := svcStoreStatusMap[typ]
	return fs
}
