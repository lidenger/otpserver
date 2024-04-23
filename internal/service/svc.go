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
func DoubleWrite(storeDetectionEventChan chan<- struct{}, exec, execBackup doubleWriteFunc) error {
	// 主存储
	tx, err := exec()
	if err != nil {
		go func() {
			storeDetectionEventChan <- struct{}{}
		}()
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
		go func() {
			storeDetectionEventChan <- struct{}{}
		}()
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
func MultiStoreInsert[T any](ctx context.Context, storeDetectionEventChan chan<- struct{}, m T, stores ...store.InsertFunc[T]) error {
	main := stores[0]
	backup := stores[1]
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
	err := DoubleWrite(storeDetectionEventChan, func() (store.Tx, error) {
		return main.Insert(ctx, m)
	}, backupExec)
	// 其他store
	for _, s := range stores[2:] {
		tx, _ := s.Insert(ctx, m)
		if tx != nil {
			tx.Commit()
		}
	}
	return err
}

// MultiStoreSelectByCondition 多store 条件查询
func MultiStoreSelectByCondition[P any, R any](ctx context.Context, storeDetectionEventChan chan<- struct{}, p P, stores ...store.SelectFunc[P, R]) (result []R, err error) {
	var cacheStore store.CacheStore
	for _, s := range stores {
		err = storeHealthCheck(s)
		if err != nil {
			continue
		}
		isCacheStore := false
		if x, ok := s.(store.CacheStore); ok {
			isCacheStore = ok
			cacheStore = x
		} else {
			isCacheStore = ok
		}
		result, err = s.SelectByCondition(ctx, p)
		if err != nil {
			go func() {
				storeDetectionEventChan <- struct{}{}
			}()
		}
		if err == nil && result != nil {
			// cache中没有，存储中存在，更新cache
			if !isCacheStore && cacheStore != nil {
				_ = cacheStore.Refresh(ctx, p)
			}
			return
		}
	}
	return
}

func MultiStorePaging[P any, R any](ctx context.Context, storeDetectionEventChan chan<- struct{}, p P, stores ...store.PagingFunc[P, R]) (result []R, count int64, err error) {
	for _, s := range stores {
		err = storeHealthCheck(s)
		if err != nil {
			continue
		}
		result, count, err = s.Paging(ctx, p)
		if err != nil {
			go func() {
				storeDetectionEventChan <- struct{}{}
			}()
		}
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

func MultiStoreUpdate(ctx context.Context, storeDetectionEventChan chan<- struct{}, ID int64, params map[string]any, stores ...store.UpdateFunc) error {
	main := stores[0]
	backup := stores[1]
	if main.GetStoreErr() != nil {
		return otperr.ErrStore(main.GetStoreErr())
	}
	var backupExec doubleWriteFunc = nil
	if backup != nil {
		if backup.GetStoreErr() != nil {
			return otperr.ErrStore(backup.GetStoreErr())
		}
		backupExec = func() (store.Tx, error) {
			return backup.Update(ctx, ID, params)
		}
	}
	err := DoubleWrite(storeDetectionEventChan, func() (store.Tx, error) {
		return main.Update(ctx, ID, params)
	}, backupExec)
	// 其他store
	for _, s := range stores[2:] {
		tx, _ := s.Update(ctx, ID, params)
		if tx != nil {
			tx.Commit()
		}
	}
	return err
}

func MultiStoreDelete(ctx context.Context, storeDetectionEventChan chan<- struct{}, ID int64, stores ...store.DeleteFunc) error {
	main := stores[0]
	backup := stores[1]
	if main.GetStoreErr() != nil {
		return otperr.ErrStore(main.GetStoreErr())
	}
	var backupExec doubleWriteFunc = nil
	if backup != nil {
		if backup.GetStoreErr() != nil {
			return otperr.ErrStore(backup.GetStoreErr())
		}
		backupExec = func() (store.Tx, error) {
			return backup.Delete(ctx, ID)
		}
	}
	err := DoubleWrite(storeDetectionEventChan, func() (store.Tx, error) {
		return main.Delete(ctx, ID)
	}, backupExec)
	// 其他store
	for _, s := range stores[2:] {
		tx, _ := s.Delete(ctx, ID)
		if tx != nil {
			tx.Commit()
		}
	}
	return err
}
