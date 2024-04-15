package pgsqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/gormstore"
	"gorm.io/gorm"
)

type SecretStore struct {
	DB  *gorm.DB
	err error
}

// 确保SecretStore实现了store.SecretStore
var _ store.SecretStore = (*SecretStore)(nil)

func (s *SecretStore) GetStoreErr() error {
	return s.err
}

func (s *SecretStore) SetStoreErr(err error) {
	s.err = err
}

func (s *SecretStore) GetDB() *gorm.DB {
	return s.DB
}

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (store.Tx, error) {
	return gormstore.SecretInsert(ctx, s.DB, m)
}

func (s *SecretStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return gormstore.SecretUpdate(ctx, s.DB, ID, params)
}

func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	return gormstore.SecretPaging(ctx, s.DB, param)
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) (result []*model.AccountSecretModel, err error) {
	return gormstore.SecretSelectByCondition(ctx, s.DB, condition)
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return gormstore.SecretDelete(ctx, s.DB, ID)
}
