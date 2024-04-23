package mysqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/internal/store/gormstore"
	"github.com/lidenger/otpserver/pkg/enum"
	"gorm.io/gorm"
)

type ServerIpListStore struct {
	DB  *gorm.DB
	err error
}

// 确保ServerIpListStore实现了store.ServerIpListStore
var _ store.ServerIpListStore = (*ServerIpListStore)(nil)

func (s *ServerIpListStore) GetStoreErr() error {
	return s.err
}

func (s *ServerIpListStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ServerIpListStore) GetStoreType() string {
	return enum.MySQLStore
}

func (s *ServerIpListStore) Insert(ctx context.Context, m *model.ServerIpListModel) (store.Tx, error) {
	return gormstore.ServerIpListInsert(ctx, s.DB, m)
}

func (s *ServerIpListStore) SelectById(ctx context.Context, ID int64) (*model.ServerIpListModel, error) {
	return gormstore.ServerIpListSelectById(ctx, s.DB, ID)
}

func (s *ServerIpListStore) SelectByCondition(ctx context.Context, condition *param.ServerIpListParam) (result []*model.ServerIpListModel, err error) {
	return gormstore.ServerIpListSelectByCondition(ctx, s.DB, condition)
}

func (s *ServerIpListStore) SelectAll(ctx context.Context) ([]*model.ServerIpListModel, error) {
	return gormstore.ServerIpListSelectAll(ctx, s.DB)
}

func (s *ServerIpListStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return gormstore.ServerIPListDelete(ctx, s.DB, ID)
}
