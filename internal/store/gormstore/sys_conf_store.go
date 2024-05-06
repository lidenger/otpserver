package gormstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
	"time"
)

func ConfInsert(_ context.Context, DB *gorm.DB, m *model.SysConfModel) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Create(m)
	return getTx(tx), tx.Error
}

func ConfUpdate(_ context.Context, DB *gorm.DB, ID int64, params map[string]any) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Model(&model.SysConfModel{})
	tx = tx.Where("id = ?", ID)
	params["update_time"] = time.Now()
	tx = tx.Updates(params)
	return getTx(tx), tx.Error
}

func ConfPaging(_ context.Context, DB *gorm.DB, param *param.SysConfPagingParam) (result []*model.SysConfModel, count int64, err error) {
	db := ConfigPagingParam(param.PageNo, param.PageSize, DB)
	if param.SearchTxt != "" {
		db = db.Where("(sys_key like ? or remark like ?)", "%"+param.SearchTxt+"%", "%"+param.SearchTxt+"%")
	}
	err = db.Order("update_time desc").Find(&result).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return
}

func ConfSelectByCondition(_ context.Context, DB *gorm.DB, condition *param.SysConfParam) (result []*model.SysConfModel, err error) {
	db := DB
	if condition.ID != 0 {
		db = db.Where("ID = ?", condition.ID)
	}
	if condition.Key != "" {
		db = db.Where("sys_key = ?", condition.Key)
	}
	err = db.Order("update_time desc").Find(&result).Error
	return
}

func ConfSelectAll(_ context.Context, DB *gorm.DB) (result []*model.SysConfModel, err error) {
	err = DB.Order("update_time desc").Find(&result).Error
	return
}

func ConfSelectById(_ context.Context, DB *gorm.DB, ID int64) (result *model.SysConfModel, err error) {
	err = DB.First(&result, "id = ?", ID).Error
	return
}

func ConfDelete(_ context.Context, DB *gorm.DB, ID int64) (store.Tx, error) {
	tx := DB.Begin()
	tx = tx.Delete(&model.SysConfModel{}, ID)
	return getTx(tx), tx.Error
}
