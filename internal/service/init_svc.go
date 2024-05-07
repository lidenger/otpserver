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
var ServerIpListSvcIns = &ServerIpListSvc{}
var ConfSvcIns = &ConfSvc{}

// store type | healthFunc
var svcStoreStatusMap = make(map[string][]store.HealthFunc)
var loadAllStoreArr = make([]store.LoadAllFunc, 0)

func Initialize(storeDetectionEventChan chan<- struct{}) {
	conf := serverconf.GetSysConf()

	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan
	SecretSvcIns.storeDetectionEventChan = storeDetectionEventChan
	ServerIpListSvcIns.storeDetectionEventChan = storeDetectionEventChan
	ConfSvcIns.storeDetectionEventChan = storeDetectionEventChan

	mysqlSecretStore := &mysqlstore.SecretStore{DB: mysqlconf.DB}
	mysqlServerStore := &mysqlstore.ServerStore{DB: mysqlconf.DB}
	mysqlServerIpListStore := &mysqlstore.ServerIpListStore{DB: mysqlconf.DB}
	mysqlConfStore := &mysqlstore.ConfStore{DB: mysqlconf.DB}

	pgsqlSecretStore := &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	pgsqlServerStore := &pgsqlstore.ServerStore{DB: pgsqlconf.DB}
	pgsqlServerIpListStore := &pgsqlstore.ServerIpListStore{DB: pgsqlconf.DB}
	pgsqlConfStore := &pgsqlstore.ConfStore{DB: pgsqlconf.DB}

	localStoreRootPath := conf.Server.RootPath + conf.Store.RootPath
	localSecretStore := &localstore.SecretStore{RootPath: localStoreRootPath}
	localServerStore := &localstore.ServerStore{RootPath: localStoreRootPath}
	localServerIpListStore := &localstore.ServerIpListStore{RootPath: localStoreRootPath}
	localConfStore := &localstore.ConfStore{RootPath: localStoreRootPath}

	memoryDependSecretStoreArr := make([]store.SecretStore, 0)
	memoryDependServerStoreArr := make([]store.ServerStore, 0)
	memoryDependServerIpListStoreArr := make([]store.ServerIpListStore, 0)
	memoryDependConfStoreArr := make([]store.ConfStore, 0)

	switch conf.Store.MainStore {
	case enum.MySQLStore:
		SecretSvcIns.Store = mysqlSecretStore
		ServerSvcIns.Store = mysqlServerStore
		ServerIpListSvcIns.Store = mysqlServerIpListStore
		ConfSvcIns.Store = mysqlConfStore

		localSecretStore.Store = mysqlSecretStore
		localServerStore.Store = mysqlServerStore
		localServerIpListStore.Store = mysqlServerIpListStore
		localConfStore.Store = mysqlConfStore
		addSvcStore(enum.MySQLStore, SecretSvcIns.Store, ServerSvcIns.Store, ServerIpListSvcIns.Store, ConfSvcIns.Store)
	case enum.PostGreSQLStore:
		SecretSvcIns.Store = pgsqlSecretStore
		ServerSvcIns.Store = pgsqlServerStore
		ServerIpListSvcIns.Store = pgsqlServerIpListStore
		ConfSvcIns.Store = pgsqlConfStore

		localSecretStore.Store = pgsqlSecretStore
		localServerStore.Store = pgsqlServerStore
		localServerIpListStore.Store = pgsqlServerIpListStore
		localConfStore.Store = pgsqlConfStore
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.Store, ServerSvcIns.Store, ServerIpListSvcIns.Store, ConfSvcIns.Store)
	}

	switch conf.Store.BackupStore {
	case enum.MySQLStore:
		SecretSvcIns.StoreBackup = mysqlSecretStore
		ServerSvcIns.StoreBackup = mysqlServerStore
		ServerIpListSvcIns.StoreBackup = mysqlServerIpListStore
		ConfSvcIns.StoreBackup = mysqlConfStore
		addSvcStore(enum.MySQLStore, SecretSvcIns.StoreBackup, ServerSvcIns.StoreBackup, ServerIpListSvcIns.StoreBackup, ConfSvcIns.StoreBackup)
	case enum.PostGreSQLStore:
		SecretSvcIns.StoreBackup = pgsqlSecretStore
		ServerSvcIns.StoreBackup = pgsqlServerStore
		ServerIpListSvcIns.StoreBackup = pgsqlServerIpListStore
		ConfSvcIns.StoreBackup = pgsqlConfStore
		addSvcStore(enum.PostGreSQLStore, SecretSvcIns.StoreBackup, ServerSvcIns.StoreBackup, ServerIpListSvcIns.StoreBackup, ConfSvcIns.StoreBackup)
	}

	memoryDependSecretStoreArr = append(memoryDependSecretStoreArr, SecretSvcIns.Store, SecretSvcIns.StoreBackup)
	memoryDependServerStoreArr = append(memoryDependServerStoreArr, ServerSvcIns.Store, ServerSvcIns.StoreBackup)
	memoryDependServerIpListStoreArr = append(memoryDependServerIpListStoreArr, ServerIpListSvcIns.Store, ServerIpListSvcIns.StoreBackup)
	memoryDependConfStoreArr = append(memoryDependConfStoreArr, ConfSvcIns.Store, ConfSvcIns.StoreBackup)

	if conf.Store.IsEnableLocal {
		memoryDependSecretStoreArr = append(memoryDependSecretStoreArr, localSecretStore)
		memoryDependServerStoreArr = append(memoryDependServerStoreArr, localServerStore)
		memoryDependServerIpListStoreArr = append(memoryDependServerIpListStoreArr, localServerIpListStore)
		memoryDependConfStoreArr = append(memoryDependConfStoreArr, localConfStore)

		loadAllStoreArr = append(loadAllStoreArr, localSecretStore, localServerStore, localServerIpListStore, localConfStore)
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

		serverIpListMemory := &memorystore.ServerIpListStore{
			Stores:                  memoryDependServerIpListStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		ServerIpListSvcIns.StoreMemory = serverIpListMemory

		ConfMemory := &memorystore.ConfStore{
			Stores:                  memoryDependConfStoreArr,
			StoreDetectionEventChan: storeDetectionEventChan,
		}
		ConfSvcIns.StoreMemory = ConfMemory

		loadAllStoreArr = append(loadAllStoreArr, secretMemory, serverMemory, serverIpListMemory, ConfMemory)

		addSvcStore(enum.MemoryStore, SecretSvcIns.StoreMemory, ServerSvcIns.StoreMemory, ServerIpListSvcIns.StoreMemory, ConfSvcIns.StoreMemory)
	}

	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc", "ServerIpListSvc", "ConfSvc")
}

func LoadAllData() {
	for _, s := range loadAllStoreArr {
		if s.GetStoreType() == enum.LocalStore || s.GetStoreType() == enum.MemoryStore {
			_ = s.LoadAll(context.Background())
		}
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
