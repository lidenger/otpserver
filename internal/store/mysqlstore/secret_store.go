package mysqlstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
	"time"
)

type SecretStore struct {
	db *gorm.DB
}

// 确保SecretStore实现了store.SecretStore
var _ store.SecretStore = (*SecretStore)(nil)

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (*gorm.DB, error) {
	tx := s.db.Begin()
	tx = tx.Create(m)
	return tx, tx.Error
}

func (s *SecretStore) Update(ctx context.Context, ID int64, params map[string]any) error {
	db := s.db.Model(&model.AccountSecretModel{})
	db = db.Where("id = ?", ID)
	params["update_time"] = time.Now()
	db = db.Updates(params)
	return db.Error
}

func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	db := store.ConfigPagingParam(param.PageNo, param.PageSize, s.db)
	if param.SearchTxt != "" {
		db = db.Where("account like ?", "%"+param.SearchTxt+"%")
	}
	err = db.Order("update_time desc").Find(&result).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) (result []*model.AccountSecretModel, err error) {
	db := s.db
	if condition.ID != 0 {
		db.Where("ID = ?", condition.ID)
	}
	if condition.IsEnable != 0 {
		db.Where("is_enable = ?", condition.IsEnable)
	}
	if condition.Account != "" {
		db.Where("account = ?", condition.Account)
	}
	err = db.Order("update_time desc").Find(&result).Error
	return
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) error {
	db := s.db
	db = db.Delete(&model.AccountSecretModel{}, ID)
	return db.Error
}
