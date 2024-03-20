package mysqlstore

import (
	"github.com/lidenger/otpserver/internal/store"
	"gorm.io/gorm"
)

type mysqlTx struct {
	tx *gorm.DB
}

func (t *mysqlTx) Commit() {
	t.tx.Commit()
}
func (t *mysqlTx) Rollback() {
	t.tx.Rollback()
}

func getTx(db *gorm.DB) store.Tx {
	return &mysqlTx{tx: db}
}
