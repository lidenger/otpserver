package pgsqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/gormstore"
	"gorm.io/gorm"
)

type ServerStore struct {
	DB  *gorm.DB
	err error
}

// 确保SecretStore实现了store.SecretStore
var _ store.ServerStore = (*ServerStore)(nil)

func (s *ServerStore) GetStoreErr() error {
	return s.err
}

func (s *ServerStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ServerStore) Insert(ctx context.Context, m *model.ServerModel) (store.Tx, error) {
	return gormstore.ServerInsert(ctx, s.DB, m)
}

func (s *ServerStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return gormstore.ServerUpdate(ctx, s.DB, ID, params)
}

func (s *ServerStore) Paging(ctx context.Context, param *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	return gormstore.ServerPaging(ctx, s.DB, param)
}

func (s *ServerStore) SelectById(ctx context.Context, ID int64) (*model.ServerModel, error) {
	return gormstore.ServerSelectById(ctx, s.DB, ID)
}

func (s *ServerStore) SelectByCondition(ctx context.Context, condition *param.ServerParam) (result []*model.ServerModel, err error) {
	return gormstore.ServerSelectByCondition(ctx, s.DB, condition)
}

func (s *ServerStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return gormstore.ServerDelete(ctx, s.DB, ID)
}
