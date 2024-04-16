package memorystore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type SecretStore struct {
	Stores []store.SecretStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	err    error
}

// account | model
var secretCacheMap = make(map[string]*model.AccountSecretModel)

// 确保SecretStore实现了store.SecretStore和store.CacheStore
var _ store.SecretStore = (*SecretStore)(nil)
var _ store.CacheStore = (*SecretStore)(nil)

func (s *SecretStore) GetStoreErr() error {
	return s.err
}

func (s *SecretStore) SetStoreErr(err error) {
	s.err = err
}

func (s *SecretStore) getAvailableStore() store.SecretStore {
	for _, st := range s.Stores {
		if st.GetStoreErr() == nil {
			return st
		}
	}
	return nil
}

func (s *SecretStore) LoadAll(ctx context.Context) error {
	p := &param.SecretPagingParam{}
	p.PageNo = 1
	p.PageSize = 100
	for {
		data, _, err := s.getAvailableStore().Paging(ctx, p)
		if err != nil {
			return err
		}
		// 获取了所有数据
		if len(data) == 0 {
			break
		}
		for _, m := range data {
			secretCacheMap[m.Account] = m
		}
		p.PageNo++
	}
	return nil
}

func (s *SecretStore) Remove(ctx context.Context, account string) {
	delete(secretCacheMap, account)
}

func (s *SecretStore) Refresh(ctx context.Context, account string) error {
	// 从存储中获取一份
	p := &param.SecretParam{Account: account}
	ms, err := s.getAvailableStore().SelectByCondition(ctx, p)
	if err != nil {
		return err
	}
	m := util.GetArrFirstItem(ms)
	if m == nil {
		// 存储中没有，删除缓存
		s.Remove(ctx, account)
	} else {
		secretCacheMap[account] = m
	}
	return nil
}

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (store.Tx, error) {
	err := s.Refresh(ctx, m.Account)
	return store.EmptyTxIns, err
}

func (s *SecretStore) Update(ctx context.Context, ID int64, _ map[string]any) (store.Tx, error) {
	m, err := s.getAvailableStore().SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return store.EmptyTxIns, nil
	}
	secretCacheMap[m.Account] = m
	return store.EmptyTxIns, err
}

// Paging 暂不提供内存分页功能
func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	err = otperr.ErrServerFuncNonsupport("暂不提供Memory分页Secret功能")
	return
}

func (s *SecretStore) SelectById(ctx context.Context, ID int64) (*model.AccountSecretModel, error) {
	for _, m := range secretCacheMap {
		if m.ID == ID {
			return m, nil
		}
	}
	return nil, nil
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) ([]*model.AccountSecretModel, error) {
	if len(condition.Account) < 1 {
		return nil, nil
	}
	m := secretCacheMap[condition.Account]
	result := make([]*model.AccountSecretModel, 0)
	result = append(result, m)
	return result, nil
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	m, err := s.SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m != nil {
		s.Remove(ctx, m.Account)
	}
	return store.EmptyTxIns, nil
}
