package memorystore

import (
	"context"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
	"strconv"
	"time"
)

type SecretStore struct {
	Stores                  []store.SecretStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	StoreDetectionEventChan chan<- struct{}
	err                     error
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

func (s *SecretStore) GetStoreType() string {
	return enum.MemoryStore
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
	secretStore := s.getAvailableStore()
	for {
		secretStore = s.getAvailableStore()
		data, _, err := secretStore.Paging(ctx, p)
		// 查询异常做一次store的检测，重新查询一次
		if err != nil {
			s.StoreDetectionEventChan <- struct{}{}
			time.Sleep(3 * time.Second)
			secretStore = s.getAvailableStore()
			data, _, err = secretStore.Paging(ctx, p)
		}
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
	log.Info("Memory从[" + secretStore.GetStoreType() + "]存储中获取数据成功，总数: " + strconv.Itoa(len(secretCacheMap)))
	return nil
}

func (s *SecretStore) Remove(_ context.Context, p any) {
	if sp, ok := p.(*param.SecretParam); ok {
		delete(secretCacheMap, sp.Account)
	}
}

func (s *SecretStore) Refresh(ctx context.Context, par any) error {
	account := ""
	if sp, ok := par.(*param.SecretParam); ok {
		account = sp.Account
	} else {
		return errors.New("非法参数")
	}
	// 从存储中获取一份
	p := &param.SecretParam{Account: account}
	ms, err := s.getAvailableStore().SelectByCondition(ctx, p)
	if err != nil {
		s.StoreDetectionEventChan <- struct{}{}
		time.Sleep(3 * time.Second)
		ms, err = s.getAvailableStore().SelectByCondition(ctx, p)
	}
	if err != nil {
		return err
	}
	m := util.GetArrFirstItem(ms)
	if m == nil {
		// 存储中没有，删除缓存
		s.Remove(ctx, p)
	} else {
		secretCacheMap[p.Account] = m
	}
	return nil
}

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (store.Tx, error) {
	err := s.Refresh(ctx, &param.SecretParam{Account: m.Account})
	return store.EmptyTxIns, err
}

func (s *SecretStore) Update(ctx context.Context, ID int64, _ map[string]any) (store.Tx, error) {
	m, err := s.getAvailableStore().SelectById(ctx, ID)
	if err != nil {
		s.StoreDetectionEventChan <- struct{}{}
		time.Sleep(3 * time.Second)
		m, err = s.getAvailableStore().SelectById(ctx, ID)
	}
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
func (s *SecretStore) Paging(_ context.Context, _ *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	err = otperr.ErrServerFuncNonsupport("暂不提供Memory分页Secret功能")
	return
}

func (s *SecretStore) SelectById(_ context.Context, ID int64) (*model.AccountSecretModel, error) {
	for _, m := range secretCacheMap {
		if m.ID == ID {
			// 不影响原始数据
			m2 := &model.AccountSecretModel{}
			err := copier.Copy(m2, m)
			if err != nil {
				return nil, err
			}
			return m2, nil
		}
	}
	return nil, nil
}

func (s *SecretStore) SelectByCondition(_ context.Context, condition *param.SecretParam) ([]*model.AccountSecretModel, error) {
	if len(condition.Account) < 1 {
		return nil, nil
	}
	m, ok := secretCacheMap[condition.Account]
	if !ok {
		return nil, nil
	}
	result := make([]*model.AccountSecretModel, 0)
	// 不影响原始数据
	m2 := &model.AccountSecretModel{}
	err := copier.Copy(m2, m)
	if err != nil {
		return nil, err
	}
	result = append(result, m2)
	return result, nil
}

func (s *SecretStore) SelectAll(ctx context.Context) ([]*model.AccountSecretModel, error) {
	result := make([]*model.AccountSecretModel, 0)
	for _, m := range secretCacheMap {
		result = append(result, m)
	}
	return result, nil
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	m, err := s.SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m != nil {
		s.Remove(ctx, &param.SecretParam{Account: m.Account})
	}
	return store.EmptyTxIns, nil
}
