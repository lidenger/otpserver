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

type ConfStore struct {
	Stores                  []store.ConfStore // 可用的持久化存储，例如：MySQL,PostgreSQL,本地文件
	StoreDetectionEventChan chan<- struct{}
	err                     error
}

// key | model
var confCacheMap = make(map[string]*model.SysConfModel)

var _ store.ConfStore = (*ConfStore)(nil)
var _ store.CacheStore = (*ConfStore)(nil)

func (s *ConfStore) GetStoreErr() error {
	return s.err
}

func (s *ConfStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ConfStore) GetStoreType() string {
	return enum.MemoryStore
}

func (s *ConfStore) getAvailableStore() store.ConfStore {
	for _, st := range s.Stores {
		if st.GetStoreErr() == nil {
			return st
		}
	}
	return nil
}

func (s *ConfStore) LoadAll(ctx context.Context) error {
	confStore := s.getAvailableStore()
	data, err := confStore.SelectAll(ctx)
	if err != nil {
		time.Sleep(3 * time.Second)
		confStore = s.getAvailableStore()
		data, err = confStore.SelectAll(ctx)
	}
	if err != nil {
		return err
	}
	for _, m := range data {
		confCacheMap[m.Key] = m
	}
	log.Info("Memory从[" + confStore.GetStoreType() + "]存储中获取系统配置数据成功，总数: " + strconv.Itoa(len(confCacheMap)))
	return nil
}

func (s *ConfStore) Remove(_ context.Context, p any) {
	if sp, ok := p.(*param.SysConfParam); ok {
		delete(confCacheMap, sp.Key)
	}
}

func (s *ConfStore) Refresh(ctx context.Context, par any) error {
	key := ""
	if sp, ok := par.(*param.SysConfParam); ok {
		key = sp.Key
	} else {
		return errors.New("非法参数")
	}
	// 从存储中获取一份
	p := &param.SysConfParam{Key: key}
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
		confCacheMap[p.Key] = m
	}
	return nil
}

func (s *ConfStore) Insert(ctx context.Context, m *model.SysConfModel) (store.Tx, error) {
	err := s.Refresh(ctx, &param.SysConfParam{Key: m.Key})
	return store.EmptyTxIns, err
}

// 从可用的store中获取数据并更新缓存
func (s *ConfStore) selectByIdAndRefresh(ctx context.Context, ID int64, isRefresh bool) (*model.SysConfModel, error) {
	m, err := s.getAvailableStore().SelectById(ctx, ID)
	if err != nil {
		s.StoreDetectionEventChan <- struct{}{}
		time.Sleep(3 * time.Second)
		m, err = s.getAvailableStore().SelectById(ctx, ID)
	}
	if m != nil && isRefresh {
		err = s.Refresh(ctx, &param.SysConfParam{Key: m.Key})
	}
	m2 := &model.SysConfModel{}
	err = copier.Copy(m2, m)
	return m2, err
}

func (s *ConfStore) Update(ctx context.Context, ID int64, _ map[string]any) (store.Tx, error) {
	_, err := s.selectByIdAndRefresh(ctx, ID, true)
	return store.EmptyTxIns, err
}

// Paging 暂不提供内存分页功能
func (s *ConfStore) Paging(_ context.Context, _ *param.SysConfPagingParam) (result []*model.SysConfModel, count int64, err error) {
	err = otperr.ErrServerFuncNonsupport("暂不提供Memory分页系统配置功能")
	return
}

func (s *ConfStore) SelectById(_ context.Context, ID int64) (*model.SysConfModel, error) {
	for _, m := range confCacheMap {
		if m.ID == ID {
			m2 := &model.SysConfModel{}
			err := copier.Copy(m2, m)
			return m2, err
		}
	}
	return nil, nil
}

func (s *ConfStore) SelectByCondition(_ context.Context, condition *param.SysConfParam) ([]*model.SysConfModel, error) {
	if len(condition.Key) < 1 {
		return nil, nil
	}
	m, ok := confCacheMap[condition.Key]
	if !ok {
		return nil, nil
	}
	result := make([]*model.SysConfModel, 0)
	// 不影响原始数据
	m2 := &model.SysConfModel{}
	err := copier.Copy(m2, m)
	if err != nil {
		return nil, err
	}
	result = append(result, m2)
	return result, nil
}

func (s *ConfStore) SelectAll(_ context.Context) ([]*model.SysConfModel, error) {
	result := make([]*model.SysConfModel, 0)
	for _, m := range confCacheMap {
		m2 := &model.SysConfModel{}
		err := copier.Copy(m2, m)
		if err != nil {
			return nil, err
		}
		result = append(result, m2)
	}
	return result, nil
}

func (s *ConfStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	m, err := s.SelectById(ctx, ID)
	if err != nil {
		return store.EmptyTxIns, err
	}
	if m != nil {
		delete(confCacheMap, m.Key)
	}
	return store.EmptyTxIns, err
}
