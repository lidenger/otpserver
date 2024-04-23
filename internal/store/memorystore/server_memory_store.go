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

type ServerStore struct {
	Stores                  []store.ServerStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	StoreDetectionEventChan chan<- struct{}
	err                     error
}

// sign | model
var serverCacheMap = make(map[string]*model.ServerModel)

// 确保SecretStore实现了store.SecretStore
var _ store.ServerStore = (*ServerStore)(nil)
var _ store.CacheStore = (*ServerStore)(nil)

func (s *ServerStore) GetStoreErr() error {
	return s.err
}

func (s *ServerStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ServerStore) GetStoreType() string {
	return enum.MemoryStore
}

func (s *ServerStore) getAvailableStore() store.ServerStore {
	for _, st := range s.Stores {
		if st.GetStoreErr() == nil {
			return st
		}
	}
	return nil
}

// LoadAll cache store func
func (s *ServerStore) LoadAll(ctx context.Context) error {
	serverStore := s.getAvailableStore()
	data, err := serverStore.SelectAll(ctx)
	if err != nil {
		time.Sleep(3 * time.Second)
		serverStore = s.getAvailableStore()
		data, err = serverStore.SelectAll(ctx)
	}
	if err != nil {
		return err
	}
	for _, m := range data {
		serverCacheMap[m.Sign] = m
	}
	log.Info("Memory从[" + serverStore.GetStoreType() + "]存储中获取接入服务数据成功，总数: " + strconv.Itoa(len(serverCacheMap)))
	return nil
}

// Remove cache store func
func (s *ServerStore) Remove(_ context.Context, p any) {
	if sp, ok := p.(*param.ServerParam); ok {
		delete(serverCacheMap, sp.Sign)
	}
}

// Refresh cache store func
func (s *ServerStore) Refresh(ctx context.Context, par any) error {
	serverSign := ""
	if sp, ok := par.(*param.ServerParam); ok {
		serverSign = sp.Sign
	} else {
		return errors.New("非法参数")
	}
	p := &param.ServerParam{Sign: serverSign}
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
		s.Remove(ctx, par)
	} else {
		// 更新缓存
		serverCacheMap[serverSign] = m
	}
	return nil
}

func (s *ServerStore) Insert(ctx context.Context, m *model.ServerModel) (store.Tx, error) {
	err := s.Refresh(ctx, &param.ServerParam{Sign: m.Sign})
	return store.EmptyTxIns, err
}

func (s *ServerStore) Update(ctx context.Context, ID int64, _ map[string]any) (store.Tx, error) {
	m, err := s.getAvailableStore().SelectById(ctx, ID)
	if err != nil {
		s.StoreDetectionEventChan <- struct{}{}
		time.Sleep(3 * time.Second)
		m, err = s.getAvailableStore().SelectById(ctx, ID)
	}
	if err != nil {
		return nil, err
	}
	if m != nil {
		serverCacheMap[m.Sign] = m
	}
	return store.EmptyTxIns, err
}

func (s *ServerStore) Paging(_ context.Context, _ *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	err = otperr.ErrServerFuncNonsupport("暂不提供Memory分页Server功能")
	return
}

func (s *ServerStore) SelectById(_ context.Context, ID int64) (*model.ServerModel, error) {
	for _, m := range serverCacheMap {
		if m.ID == ID {
			m2 := &model.ServerModel{}
			err := copier.Copy(m2, m)
			if err != nil {
				return nil, err
			}
			return m2, nil
		}
	}
	return nil, nil
}

func (s *ServerStore) SelectByCondition(_ context.Context, condition *param.ServerParam) ([]*model.ServerModel, error) {
	if len(condition.Sign) < 1 {
		return nil, nil
	}
	m, ok := serverCacheMap[condition.Sign]
	if !ok {
		return nil, nil
	}
	result := make([]*model.ServerModel, 0)
	// 不影响原始数据
	m2 := &model.ServerModel{}
	err := copier.Copy(m2, m)
	if err != nil {
		return nil, err
	}
	result = append(result, m2)
	return result, nil
}

func (s *ServerStore) SelectAll(_ context.Context) ([]*model.ServerModel, error) {
	result := make([]*model.ServerModel, 0)
	for _, m := range serverCacheMap {
		m2 := &model.ServerModel{}
		err := copier.Copy(m2, m)
		if err != nil {
			return nil, err
		}
		result = append(result, m2)
	}
	return result, nil
}

func (s *ServerStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	m, err := s.SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m != nil {
		delete(serverCacheMap, m.Sign)
	}
	return store.EmptyTxIns, nil
}
