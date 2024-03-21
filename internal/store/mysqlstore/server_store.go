package mysqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
)

type ServerStore struct {
	DB *gorm.DB
}

// 确保ServerStore实现了store.ServerStore
var _ store.ServerStore = (*ServerStore)(nil)

func (s *ServerStore) Insert(ctx context.Context, m *model.ServerModel) (store.Tx, error) {
	tx := s.DB.Begin()
	tx = tx.Create(m)
	return getTx(tx), tx.Error
}

func (s *ServerStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerStore) Paging(ctx context.Context, param *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *ServerStore) SelectByCondition(ctx context.Context, condition *param.ServerParam) (result []*model.ServerModel, err error) {
	db := s.DB
	if condition.ID != 0 {
		db = db.Where("ID = ?", condition.ID)
	}
	if condition.IsEnable != 0 {
		db = db.Where("is_enable = ?", condition.IsEnable)
	}
	if condition.Name != "" {
		db = db.Where("server_name = ?", condition.Name)
	}
	if condition.Sign != "" {
		db = db.Where("server_sign = ?", condition.Sign)
	}
	err = db.Order("update_time desc").Find(&result).Error
	return
}

func (s *ServerStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	//TODO implement me
	panic("implement me")
}
