package service

import (
	"context"
	"errors"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
)

type doubleWriteFunc func() (store.Tx, error)

// DoubleWrite 主备写
func DoubleWrite(exec, execBackup doubleWriteFunc) error {
	// 主存储
	tx, err := exec()
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
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
		if tx2 != nil {
			tx2.Rollback()
		}
		return otperr.ErrStoreBackup(err2)
	}
	// 主备存储成功，提交事务
	tx.Commit()
	tx2.Commit()
	return nil
}

// MultiStoreInsert 多store insert
func MultiStoreInsert[T any](ctx context.Context, m T, main, backup store.InsertFunc[T]) error {
	if main.GetStoreErr() != nil {
		return otperr.ErrStore(main.GetStoreErr())
	}
	var backupExec doubleWriteFunc = nil
	if backup != nil {
		if backup.GetStoreErr() != nil {
			return otperr.ErrStore(backup.GetStoreErr())
		}
		backupExec = func() (store.Tx, error) {
			return backup.Insert(ctx, m)
		}
	}
	err := DoubleWrite(func() (store.Tx, error) {
		return main.Insert(ctx, m)
	}, backupExec)
	return err
}

// MultiStoreSelectByCondition 多store 条件查询
func MultiStoreSelectByCondition[P any, R any](ctx context.Context, p P, stores ...store.SelectByConditionFunc[P, R]) (result []R, err error) {
	for _, s := range stores {
		err = storeHealthCheck(s)
		if err != nil {
			continue
		}
		result, err = s.SelectByCondition(ctx, p)
		if err == nil && result != nil {
			return
		}
	}
	return
}

func MultiStorePaging[P any, R any](ctx context.Context, p P, stores ...store.PagingFunc[P, R]) (result []R, count int64, err error) {
	for _, s := range stores {
		err = storeHealthCheck(s)
		if err != nil {
			continue
		}
		result, count, err = s.Paging(ctx, p)
		if err == nil && result != nil {
			return
		}
	}
	return
}

func storeHealthCheck(f store.HealthFunc) error {
	if f == nil {
		return errors.New("f is blank")
	}
	if f.GetStoreErr() != nil {
		log.Error("Store异常:%s", f.GetStoreErr())
		return otperr.ErrStore(f.GetStoreErr())
	}
	return nil
}
