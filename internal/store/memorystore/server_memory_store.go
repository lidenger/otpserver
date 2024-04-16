package memorystore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type ServerStore struct {
	Stores []store.ServerStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	err    error
}

// account | model
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

func (s *ServerStore) getAvailableStore() store.ServerStore {
	for _, st := range s.Stores {
		if st.GetStoreErr() == nil {
			return st
		}
	}
	return nil
}

func (s *ServerStore) LoadAll(ctx context.Context) error {
	p := &param.ServerPagingParam{}
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
			serverCacheMap[m.Sign] = m
		}
		p.PageNo++
	}
	return nil
}

func (s *ServerStore) Remove(ctx context.Context, serverSign string) {
	delete(serverCacheMap, serverSign)
}

func (s *ServerStore) Refresh(ctx context.Context, serverSign string) error {
	p := &param.ServerParam{Sign: serverSign}
	ms, err := s.getAvailableStore().SelectByCondition(ctx, p)
	if err != nil {
		return err
	}
	m := util.GetArrFirstItem(ms)
	if m == nil {
		// 存储中没有，删除缓存
		s.Remove(ctx, serverSign)
	} else {
		serverCacheMap[serverSign] = m
	}
	return nil
}

func (s *ServerStore) Insert(ctx context.Context, m *model.ServerModel) (store.Tx, error) {
	err := s.Refresh(ctx, m.Sign)
	return store.EmptyTxIns, err
}

func (s *ServerStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	m, err := s.getAvailableStore().SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return store.EmptyTxIns, nil
	}
	serverCacheMap[m.Sign] = m
	return store.EmptyTxIns, err
}

func (s *ServerStore) Paging(ctx context.Context, param *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	err = otperr.ErrServerFuncNonsupport("暂不提供Memory分页Server功能")
	return
}

func (s *ServerStore) SelectById(ctx context.Context, ID int64) (*model.ServerModel, error) {
	for _, m := range serverCacheMap {
		if m.ID == ID {
			return m, nil
		}
	}
	return nil, nil
}

func (s *ServerStore) SelectByCondition(ctx context.Context, condition *param.ServerParam) ([]*model.ServerModel, error) {
	if len(condition.Sign) < 1 {
		return nil, nil
	}
	m := serverCacheMap[condition.Sign]
	result := make([]*model.ServerModel, 0)
	result = append(result, m)
	return result, nil
}

func (s *ServerStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	m, err := s.SelectById(ctx, ID)
	if err != nil {
		return nil, err
	}
	if m != nil {
		s.Remove(ctx, m.Sign)
	}
	return store.EmptyTxIns, nil
}
