package gormstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
	"time"
)

func SecretInsert(_ context.Context, DB *gorm.DB, m *model.AccountSecretModel) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Create(m)
	return getTx(tx), tx.Error
}

func SecretUpdate(_ context.Context, DB *gorm.DB, ID int64, params map[string]any) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Model(&model.AccountSecretModel{})
	tx = tx.Where("id = ?", ID)
	params["update_time"] = time.Now()
	tx = tx.Updates(params)
	return getTx(tx), tx.Error
}

func SecretPaging(_ context.Context, DB *gorm.DB, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	db := ConfigPagingParam(param.PageNo, param.PageSize, DB)
	if param.SearchTxt != "" {
		if param.SearchTxt == "启用" {
			db = db.Where("is_enable = 1")
		} else if param.SearchTxt == "禁用" {
			db = db.Where("is_enable = 2")
		} else {
			db = db.Where("account like ?", "%"+param.SearchTxt+"%")
		}
	}
	err = db.Order("update_time desc").Find(&result).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return
}

func SecretSelectById(_ context.Context, DB *gorm.DB, ID int64) (result *model.AccountSecretModel, err error) {
	err = DB.First(&result, "id = ?", ID).Error
	return
}

func SecretSelectByCondition(_ context.Context, DB *gorm.DB, condition *param.SecretParam) (result []*model.AccountSecretModel, err error) {
	db := DB
	if condition.ID != 0 {
		db = db.Where("ID = ?", condition.ID)
	}
	if condition.IsEnable != 0 {
		db = db.Where("is_enable = ?", condition.IsEnable)
	}
	if condition.Account != "" {
		db = db.Where("account = ?", condition.Account)
	}
	err = db.Order("update_time desc").Find(&result).Error
	return
}

func SecretSelectAll(_ context.Context, DB *gorm.DB) (result []*model.AccountSecretModel, err error) {
	err = DB.Order("update_time desc").Find(&result).Error
	return
}

func SecretDelete(_ context.Context, DB *gorm.DB, ID int64) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Delete(&model.AccountSecretModel{}, ID)
	return getTx(tx), tx.Error
}
