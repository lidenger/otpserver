package pgsqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/gormstore"
	"github.com/lidenger/otpserver/pkg/enum"
	"gorm.io/gorm"
)

type ConfStore struct {
	DB  *gorm.DB
	err error
}

// 确保SecretStore实现了store.SecretStore
var _ store.ConfStore = (*ConfStore)(nil)

func (s *ConfStore) GetStoreErr() error {
	return s.err
}

func (s *ConfStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ConfStore) GetStoreType() string {
	return enum.PostGreSQLStore
}

func (s *ConfStore) GetDB() *gorm.DB {
	return s.DB
}

func (s *ConfStore) Insert(ctx context.Context, m *model.SysConfModel) (store.Tx, error) {
	return gormstore.ConfInsert(ctx, s.DB, m)
}

func (s *ConfStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return gormstore.ConfUpdate(ctx, s.DB, ID, params)
}

func (s *ConfStore) Paging(ctx context.Context, param *param.SysConfPagingParam) (result []*model.SysConfModel, count int64, err error) {
	return gormstore.ConfPaging(ctx, s.DB, param)
}

func (s *ConfStore) SelectAll(ctx context.Context) ([]*model.SysConfModel, error) {
	return gormstore.ConfSelectAll(ctx, s.DB)
}

func (s *ConfStore) SelectById(ctx context.Context, ID int64) (*model.SysConfModel, error) {
	return gormstore.ConfSelectById(ctx, s.DB, ID)
}

func (s *ConfStore) SelectByCondition(ctx context.Context, condition *param.SysConfParam) (result []*model.SysConfModel, err error) {
	return gormstore.ConfSelectByCondition(ctx, s.DB, condition)
}

func (s *ConfStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return gormstore.ConfDelete(ctx, s.DB, ID)
}
