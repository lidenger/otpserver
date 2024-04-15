package gormstore

import (
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
)

type gromTx struct {
	tx *gorm.DB
}

func (t *gromTx) Commit() {
	t.tx.Commit()
}
func (t *gromTx) Rollback() {
	t.tx.Rollback()
}

func getTx(db *gorm.DB) store.Tx {
	return &gromTx{tx: db}
}
