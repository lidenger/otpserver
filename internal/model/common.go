package model

import (
	"gorm.io/gorm"
	"time"
)

type Time struct {
	CreateTime time.Time `gorm:"column:create_time" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time" json:"update_time"`
}

func (t *Time) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	t.CreateTime = now
	t.UpdateTime = now
	return nil
}

func (t *Time) BeforeUpdate(tx *gorm.DB) (err error) {
	t.UpdateTime = time.Now()
	return nil
}

type Common struct {
	IsEnable uint8 `gorm:"column:is_enable" json:"is_enable"` // 是否启用，1启用，2禁用
	Time
}

func (c *Common) BeforeCreate(tx *gorm.DB) (err error) {
	return c.Time.BeforeCreate(tx)
}

func (c *Common) BeforeUpdate(tx *gorm.DB) (err error) {
	return c.Time.BeforeUpdate(tx)
}
