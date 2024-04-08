package service

import (
	"context"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
)

type doubleWriteFunc func() (store.Tx, error)

// DoubleWrite 主备写
func DoubleWrite(exec, execBackup doubleWriteFunc) error {
	// 主存储
	tx, err := exec()
	if err != nil {
		tx.Rollback()
		return otperr.ErrStore(err)
	}
	// 没有备存储直接提交事务
	if execBackup == nil {
		tx.Commit()
		return nil
	}
	// 备存储
	tx2, err2 := execBackup()
	if err2 != nil {
		// 为了数据一致性，主存储也需要回滚
		tx.Rollback()
		tx2.Rollback()
		return otperr.ErrStoreBackup(err2)
	}
	// 主备存储成功，提交事务
	tx.Commit()
	tx2.Commit()
	return nil
}

// MultiStoreInsert 多store insert
func MultiStoreInsert[T any](ctx context.Context, m T, main, backup store.InsertFunc[T]) error {
	var backupExec doubleWriteFunc = nil
	if backup != nil {
		backupExec = func() (store.Tx, error) {
			return backup.Insert(ctx, m)
		}
	}
	err := DoubleWrite(func() (store.Tx, error) {
		return main.Insert(ctx, m)
	}, backupExec)
	return err
}

// MultiStoreRead 多store读取
func MultiStoreRead() {

}
