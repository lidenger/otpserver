package gormstore

import "gorm.io/gorm"

// ConfigPagingParam 设置db分页参数
func ConfigPagingParam(pageNo, pageSize int, db *gorm.DB) *gorm.DB {
	if pageNo != 0 && pageSize != 0 {
		offset := (pageNo - 1) * pageSize
		db = db.Offset(offset).Limit(pageSize)
	}
	return db
}
