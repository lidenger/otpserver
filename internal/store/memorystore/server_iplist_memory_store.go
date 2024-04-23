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
	"strconv"
	"time"
)

type ServerIpListStore struct {
	Stores                  []store.ServerIpListStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	StoreDetectionEventChan chan<- struct{}
	err                     error
}

// sign | []model
var serverIpListCacheMap = make(map[string][]*model.ServerIpListModel)

// 确保SecretStore实现了store.SecretStore
var _ store.ServerIpListStore = (*ServerIpListStore)(nil)
var _ store.CacheStore = (*ServerIpListStore)(nil)

func (s *ServerIpListStore) GetStoreErr() error {
	return s.err
}

func (s *ServerIpListStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ServerIpListStore) GetStoreType() string {
	return enum.MemoryStore
}

func (s *ServerIpListStore) getAvailableStore() store.ServerIpListStore {
	for _, st := range s.Stores {
		if st.GetStoreErr() == nil {
			return st
		}
	}
	return nil
}

// LoadAll cache store func
func (s *ServerIpListStore) LoadAll(ctx context.Context) error {
	serverIpListStore := s.getAvailableStore()
	data, err := s.getAvailableStore().SelectAll(ctx)
	if err != nil {
		time.Sleep(3 * time.Second)
		serverIpListStore = s.getAvailableStore()
		data, err = serverIpListStore.SelectAll(ctx)
	}
	if err != nil {
		return err
	}
	for _, m := range data {
		ips := serverIpListCacheMap[m.ServerSign]
		if ips == nil {
			ips = []*model.ServerIpListModel{}
		}
		m2 := &model.ServerIpListModel{}
		err = copier.Copy(m2, m)
		if err != nil {
			return err
		}
		ips = append(ips, m2)
		serverIpListCacheMap[m.ServerSign] = ips
	}
	log.Info("Memory从[" + serverIpListStore.GetStoreType() + "]存储中获取服务IP白名单数据成功，总数: " + strconv.Itoa(len(serverIpListCacheMap)))
	return nil
}

// Remove cache store func
func (s *ServerIpListStore) Remove(ctx context.Context, p any) {
	_ = s.Refresh(ctx, p)
}

// Refresh cache store func 从可用的store中获取并更新到缓存中
func (s *ServerIpListStore) Refresh(ctx context.Context, par any) error {
	serverSign := ""
	if sp, ok := par.(*param.ServerIpListParam); ok {
		serverSign = sp.Sign
	} else {
		return errors.New("非法参数")
	}
	p := &param.ServerIpListParam{Sign: serverSign}
	ms, err := s.getAvailableStore().SelectByCondition(ctx, p)
	if err != nil {
		s.StoreDetectionEventChan <- struct{}{}
		time.Sleep(3 * time.Second)
		ms, err = s.getAvailableStore().SelectByCondition(ctx, p)
	}
	if err != nil {
		return err
	}
	if ms == nil {
		// 服务没有对应的IP白名单列表，直接删除缓存
		delete(serverIpListCacheMap, serverSign)
	} else {
		// 更新服务对应的IP白名单列表缓存
		serverIpListCacheMap[serverSign] = ms
	}
	return nil
}

func (s *ServerIpListStore) Insert(ctx context.Context, m *model.ServerIpListModel) (store.Tx, error) {
	err := s.Refresh(ctx, &param.ServerIpListParam{Sign: m.ServerSign})
	return store.EmptyTxIns, err
}

// Update IP白名单只有增加和删除
func (s *ServerIpListStore) Update(_ context.Context, _ int64, _ map[string]any) (store.Tx, error) {
	err := otperr.ErrServerFuncNonsupport("暂不提供Memory更新服务IP白名单功能: IP白名单只有增加和删除")
	return store.EmptyTxIns, err
}

// Delete 从cache的ID比对获取服务标识（直接通过store可能因删除而获取不到），通过服务标识更新缓存
func (s *ServerIpListStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	serverSign := ""
	for _, ms := range serverIpListCacheMap {
		for _, m := range ms {
			if m.ID == ID {
				serverSign = m.ServerSign
				goto found
			}
		}
	}
found:
	err := s.Refresh(ctx, &param.ServerIpListParam{Sign: serverSign})
	return store.EmptyTxIns, err
}

func (s *ServerIpListStore) SelectById(_ context.Context, ID int64) (*model.ServerIpListModel, error) {
	for _, ms := range serverIpListCacheMap {
		for _, m := range ms {
			if m.ID == ID {
				m2 := &model.ServerIpListModel{}
				err := copier.Copy(m2, m)
				return m2, err
			}
		}
	}
	return nil, nil
}

func (s *ServerIpListStore) SelectByCondition(_ context.Context, condition *param.ServerIpListParam) ([]*model.ServerIpListModel, error) {
	if len(condition.Sign) == 0 {
		return nil, nil
	}
	ms, ok := serverIpListCacheMap[condition.Sign]
	if !ok {
		return nil, nil
	}
	result := make([]*model.ServerIpListModel, 0)
	for _, m := range ms {
		if len(condition.IP) != 0 && m.IP != condition.IP {
			continue
		}
		m2 := &model.ServerIpListModel{}
		err := copier.Copy(m2, m)
		if err != nil {
			return nil, err
		}
		result = append(result, m2)
	}
	return result, nil
}

func (s *ServerIpListStore) SelectAll(ctx context.Context) ([]*model.ServerIpListModel, error) {
	result := make([]*model.ServerIpListModel, 0)
	for k := range serverIpListCacheMap {
		ms, err := s.SelectByCondition(ctx, &param.ServerIpListParam{Sign: k})
		if err != nil {
			return nil, err
		}
		result = append(result, ms...)
	}
	return result, nil
}
