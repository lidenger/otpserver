package localstore

import (
	"context"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/enum"
	"path/filepath"
)

type ConfStore struct {
	Store    store.ConfStore
	RootPath string
	err      error
}

const ConfFileName = "otp_sys_conf"

var _ store.ConfStore = (*ConfStore)(nil)
var _ store.LocalStore[*model.SysConfModel] = (*ConfStore)(nil)

func (s *ConfStore) GetStoreErr() error {
	return s.err
}

func (s *ConfStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ConfStore) GetStoreType() string {
	return enum.LocalStore
}

func (s *ConfStore) LoadAll(ctx context.Context) error {
	err := s.Store.GetStoreErr()
	if err != nil {
		log.Error("系统配置主存储异常，无法从主存储加载数据到本地存储，使用上个版本的本地存储数据")
		return s.Store.GetStoreErr()
	}
	err = fetchAllStoreDataAndWriteToFile[*model.SysConfModel](ctx, s.Store, s.RootPath, ConfFileName)
	return nil
}

func (s *ConfStore) FetchAll(ctx context.Context) (result []*model.SysConfModel, err error) {
	filePath := filepath.Join(s.RootPath, ConfFileName)
	result, err = fetchAll[*model.SysConfModel](ctx, filePath)
	return
}

func (s *ConfStore) Insert(ctx context.Context, m *model.SysConfModel) (store.Tx, error) {
	return store.EmptyTxIns, nil
}

func (s *ConfStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return store.EmptyTxIns, nil
}

func (s *ConfStore) Paging(ctx context.Context, param *param.SysConfPagingParam) (result []*model.SysConfModel, count int64, err error) {
	if param.PageNo != 1 {
		return nil, 0, nil
	}
	data, err := s.FetchAll(ctx)
	return data, 0, err
}

func (s *ConfStore) SelectById(ctx context.Context, ID int64) (*model.SysConfModel, error) {
	return nil, nil
}

func (s *ConfStore) SelectByCondition(ctx context.Context, condition *param.SysConfParam) ([]*model.SysConfModel, error) {
	return nil, nil
}

func (s *ConfStore) SelectAll(ctx context.Context) (result []*model.SysConfModel, err error) {
	result, err = s.FetchAll(ctx)
	return
}

func (s *ConfStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return store.EmptyTxIns, nil
}
