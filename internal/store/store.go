package store

import (
	"context"
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"gorm.io/gorm"
	"strings"
)

// SecretStore 账号密钥
type SecretStore interface {
	Insert(ctx context.Context, m *model.AccountSecretModel) (Tx, error)
	Update(ctx context.Context, ID int64, params map[string]any) (Tx, error)
	Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error)
	SelectByCondition(ctx context.Context, condition *param.SecretParam) (result []*model.AccountSecretModel, err error)
	Delete(ctx context.Context, ID int64) (Tx, error)
}

// ServerStore 接入的服务
type ServerStore interface {
	Insert(ctx context.Context, m *model.ServerModel) error
	Update(ctx context.Context, m *model.ServerModel) error
	Paging(ctx context.Context, param *param.ServerPagingParam) ([]*model.ServerModel, error)
	SelectByCondition(ctx context.Context, condition *param.ServerPagingParam) ([]*model.ServerModel, error)
	Delete(ctx context.Context, ID int64) error
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

func InitStore(conf *serverconf.Config) {
	cmd := serverconf.CMD
	if cmd.MainStore == "" {
		panic("主存储类型不能为空")
	}
	cmd.MainStore = strings.ToLower(cmd.MainStore)
	cmd.BackupStore = strings.ToLower(cmd.BackupStore)
	if cmd.MainStore == cmd.BackupStore {
		log.Warn("注意：主备存储设置一致，当前模式为弃用备存储!")
		cmd.BackupStore = ""
	}
	if cmd.MainStore == "mysql" || cmd.BackupStore == "mysql" {
		mysqlconf.InitMySQL(conf)
	} else if cmd.MainStore == "pgsql" || cmd.BackupStore == "pgsql" {
		pgsqlconf.InitPgsql(conf)
	} else {
		panic("不支持的存储类型:" + cmd.MainStore)
	}
	builder := &strings.Builder{}
	builder.WriteString("存储初始化完成,")
	builder.WriteString(fmt.Sprintf("主存储:%s,", cmd.MainStore))
	if cmd.BackupStore == "" {
		builder.WriteString("备存储:未启用")
	} else {
		builder.WriteString(fmt.Sprintf("备存储:%s,", cmd.BackupStore))
	}
	log.Info("%s", builder.String())
}
