package store

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"gorm.io/gorm"
)

type SecretStore interface {
	Insert(ctx context.Context, m *model.AccountSecretModel) error
	Update(ctx context.Context, m *model.AccountSecretModel) error
	Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error)
	SelectByCondition(ctx context.Context, condition *param.SecretParam) (result []*model.AccountSecretModel, err error)
	Delete(ctx context.Context, ID int64) error
}

type ServerStore interface {
	Insert(ctx context.Context, m *model.ServerModel) error
	Update(ctx context.Context, m *model.ServerModel) error
	Paging(ctx context.Context, param *param.ServerPagingParam) ([]*model.ServerModel, error)
	SelectByCondition(ctx context.Context, condition *param.ServerPagingParam) ([]*model.ServerModel, error)
	Delete(ctx context.Context, ID int64) error
}

// ConfigPagingParam 设置db分页参数
func ConfigPagingParam(pageNo, pageSize int, db *gorm.DB) *gorm.DB {
	if pageNo != 0 && pageSize != 0 {
		offset := (pageNo - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	return db
}
