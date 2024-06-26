package store

import (
	"context"
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/storeconf"
	"github.com/lidenger/otpserver/config/storeconf/localconf"
	"github.com/lidenger/otpserver/config/storeconf/memoryconf"
	"github.com/lidenger/otpserver/config/storeconf/mysqlconf"
	"github.com/lidenger/otpserver/config/storeconf/pgsqlconf"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/lidenger/otpserver/pkg/util"
	"strings"
)

// HealthFunc store健康状态
type HealthFunc interface {
	GetStoreErr() error
	SetStoreErr(error)
	GetStoreType() string
}

type InsertFunc[T any] interface {
	HealthFunc
	Insert(ctx context.Context, m T) (Tx, error)
}

type UpdateFunc interface {
	HealthFunc
	Update(ctx context.Context, ID int64, params map[string]any) (Tx, error)
}

type DeleteFunc interface {
	HealthFunc
	Delete(ctx context.Context, ID int64) (Tx, error)
}

type PagingFunc[P any, R any] interface {
	HealthFunc
	Paging(ctx context.Context, param P) (result []R, count int64, err error)
}

type SelectAllFunc[R any] interface {
	HealthFunc
	SelectAll(ctx context.Context) (result []R, err error)
}

type SelectFunc[P any, R any] interface {
	HealthFunc
	SelectByCondition(ctx context.Context, condition P) (result []R, err error)
	SelectById(ctx context.Context, ID int64) (R, error)
	SelectAllFunc[R]
}

// SecretStore 账号密钥
type SecretStore interface {
	InsertFunc[*model.AccountSecretModel]
	UpdateFunc
	DeleteFunc
	PagingFunc[*param.SecretPagingParam, *model.AccountSecretModel]
	SelectFunc[*param.SecretParam, *model.AccountSecretModel]
}

// ServerStore 接入的服务
type ServerStore interface {
	InsertFunc[*model.ServerModel]
	UpdateFunc
	DeleteFunc
	PagingFunc[*param.ServerPagingParam, *model.ServerModel]
	SelectFunc[*param.ServerParam, *model.ServerModel]
}

// ServerIpListStore 接入服务IP白名单
type ServerIpListStore interface {
	InsertFunc[*model.ServerIpListModel]
	DeleteFunc
	SelectFunc[*param.ServerIpListParam, *model.ServerIpListModel]
}

// ConfStore 系统配置
type ConfStore interface {
	InsertFunc[*model.SysConfModel]
	UpdateFunc
	DeleteFunc
	PagingFunc[*param.SysConfPagingParam, *model.SysConfModel]
	SelectFunc[*param.SysConfParam, *model.SysConfModel]
}

// Tx 事务，这里定义事务接口，不依赖于具体的框架实现，降低耦合
type Tx interface {
	Commit()
	Rollback()
}

type EmptyTx struct{}

func (e *EmptyTx) Commit()   {}
func (e *EmptyTx) Rollback() {}

var EmptyTxIns = &EmptyTx{}

type LoadAllFunc interface {
	GetStoreType() string
	LoadAll(ctx context.Context) error
}

// CacheStore 缓存存储
type CacheStore interface {
	LoadAllFunc
	Remove(ctx context.Context, param any)
	Refresh(ctx context.Context, param any) error
}

// LocalStore 本地存储
type LocalStore[M any] interface {
	LoadAllFunc
	FetchAll(ctx context.Context) ([]M, error)
}

var activeStore = make([]storeconf.Status, 0)

func Initialize() {
	conf := serverconf.GetSysConf()
	c := conf.Store
	if c.MainStore == "" {
		panic("主存储类型不能为空")
	}
	c.MainStore = strings.ToLower(c.MainStore)
	c.BackupStore = strings.ToLower(c.BackupStore)
	if c.MainStore == c.BackupStore {
		log.Warn("注意：主备存储设置一致，当前模式为弃用备存储!")
		c.BackupStore = ""
	}
	isKnownStore := false
	if util.Eqs(enum.MySQLStore, c.MainStore, c.BackupStore) {
		mysqlconf.Initialize(conf)
		activeStore = append(activeStore, mysqlconf.Ins)
		isKnownStore = true
	}
	if util.Eqs(enum.PostGreSQLStore, c.MainStore, c.BackupStore) {
		pgsqlconf.Initialize(conf)
		activeStore = append(activeStore, pgsqlconf.Ins)
		isKnownStore = true
	}
	if !isKnownStore {
		panic("不支持的存储类型:" + c.MainStore)
	}
	builder := &strings.Builder{}
	builder.WriteString("存储初始化完成,")
	builder.WriteString(fmt.Sprintf("主存储:%s,", c.MainStore))
	if c.BackupStore == "" {
		builder.WriteString("备存储:未启用")
	} else {
		builder.WriteString(fmt.Sprintf("备存储:%s", c.BackupStore))
	}
	log.Info("%s", builder.String())
	// 启用本地存储
	if conf.Store.IsEnableLocal {
		localconf.Initialize(conf)
		activeStore = append(activeStore, localconf.Ins)
	}
	// 启用memory存储
	if conf.Store.IsEnableMemory {
		memoryconf.Initialize(conf)
		activeStore = append(activeStore, memoryconf.Ins)
	}
}

// GetAllActiveStore 获取所有启用的store
func GetAllActiveStore() []storeconf.Status {
	return activeStore
}
