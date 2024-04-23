package gormstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
)

func ServerIpListInsert(_ context.Context, DB *gorm.DB, m *model.ServerIpListModel) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Create(m)
	return getTx(tx), tx.Error
}

func ServerIpListSelectByCondition(_ context.Context, DB *gorm.DB, condition *param.ServerIpListParam) (result []*model.ServerIpListModel, err error) {
	db := DB
	if condition.Sign != "" {
		db = db.Where("server_sign = ?", condition.Sign)
	}
	err = db.Order("update_time desc").Find(&result).Error
	return
}

func ServerIpListSelectById(_ context.Context, DB *gorm.DB, ID int64) (result *model.ServerIpListModel, err error) {
	err = DB.First(&result, "id = ?", ID).Error
	return
}

func ServerIpListSelectAll(_ context.Context, DB *gorm.DB) (result []*model.ServerIpListModel, err error) {
	err = DB.Order("update_time desc").Find(&result).Error
	return
}

func ServerIPListDelete(_ context.Context, DB *gorm.DB, ID int64) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Delete(&model.ServerIpListModel{}, ID)
	return getTx(tx), tx.Error
}
