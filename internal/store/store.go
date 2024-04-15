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
	"gorm.io/gorm"
	"strings"
)

// HealthFunc store健康状态
type HealthFunc interface {
	GetStoreErr() error
	SetStoreErr(error)
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

type SelectByConditionFunc[P any, R any] interface {
	HealthFunc
	SelectByCondition(ctx context.Context, condition P) (result []R, err error)
}

// SecretStore 账号密钥
type SecretStore interface {
	InsertFunc[*model.AccountSecretModel]
	UpdateFunc
	DeleteFunc
	PagingFunc[*param.SecretPagingParam, *model.AccountSecretModel]
	SelectByConditionFunc[*param.SecretParam, *model.AccountSecretModel]
}

// ServerStore 接入的服务
type ServerStore interface {
	InsertFunc[*model.ServerModel]
	UpdateFunc
	DeleteFunc
	PagingFunc[*param.ServerPagingParam, *model.ServerModel]
	SelectByConditionFunc[*param.ServerParam, *model.ServerModel]
}

// Tx 事务，这里定义事务接口，不依赖于具体的框架实现，降低耦合
type Tx interface {
	Commit()
	Rollback()
}

// ConfigPagingParam 设置db分页参数
func ConfigPagingParam(pageNo, pageSize int, db *gorm.DB) *gorm.DB {
	if pageNo != 0 && pageSize != 0 {
		offset := (pageNo - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	return db
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
		err := mysqlconf.MySQLConfIns.TestStore()
		if err != nil {
			log.Error("MySQL检测不通过，err:%s", err.Error())
		}
		activeStore = append(activeStore, mysqlconf.MySQLConfIns)
		isKnownStore = true
	}
	if util.Eqs(enum.PostGreSQLStore, c.MainStore, c.BackupStore) {
		pgsqlconf.Initialize(conf)
		err := pgsqlconf.PgSQLConfIns.TestStore()
		if err != nil {
			log.Error("PostgreSQL检测不通过，err:%s", err.Error())
		}
		activeStore = append(activeStore, pgsqlconf.PgSQLConfIns)
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
	if conf.Server.IsEnableLocalStore {
		localconf.Initialize(conf)
		activeStore = append(activeStore, localconf.LocalConfIns)
	}
	// 启用memory存储
	if conf.Server.IsEnableMemoryStore {
		memoryconf.Initialize(conf)
		activeStore = append(activeStore, memoryconf.MemoryConfIns)
	}
}

// GetAllActiveStore 获取所有启用的store
func GetAllActiveStore() []storeconf.Status {
	return activeStore
}
