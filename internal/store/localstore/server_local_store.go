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

type ServerStore struct {
	Store    store.ServerStore
	RootPath string
	err      error
}

const ServerFileName = "otp_server"

// 确保ServerStore实现了store.ServerStore
var _ store.ServerStore = (*ServerStore)(nil)
var _ store.LocalStore[*model.ServerModel] = (*ServerStore)(nil)

func (s *ServerStore) GetStoreType() string {
	return enum.LocalStore
}

func (s *ServerStore) GetStoreErr() error {
	return s.err
}

func (s *ServerStore) SetStoreErr(err error) {
	s.err = err
}

// LoadAll 从store中获取数据到local存储
func (s *ServerStore) LoadAll(ctx context.Context) error {
	err := s.Store.GetStoreErr()
	if err != nil {
		log.Error("接入服务主存储异常，无法从主存储加载数据到本地存储，使用上个版本的本地存储数据")
		return s.Store.GetStoreErr()
	}
	err = fetchAllStoreDataAndWriteToFile[*model.ServerModel](ctx, s.Store, s.RootPath, ServerFileName)
	return nil
}

// FetchAll 从local存储获取所有数据
func (s *ServerStore) FetchAll(ctx context.Context) (result []*model.ServerModel, err error) {
	filePath := filepath.Join(s.RootPath, ServerFileName)
	result, err = fetchAll[*model.ServerModel](ctx, filePath)
	return
}

func (s *ServerStore) Insert(ctx context.Context, m *model.ServerModel) (store.Tx, error) {
	return nil, nil
}

func (s *ServerStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return nil, nil
}

func (s *ServerStore) Paging(ctx context.Context, param *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	return nil, 0, nil
}

func (s *ServerStore) SelectById(ctx context.Context, ID int64) (*model.ServerModel, error) {
	return nil, nil
}

func (s *ServerStore) SelectByCondition(ctx context.Context, condition *param.ServerParam) ([]*model.ServerModel, error) {
	return nil, nil
}

func (s *ServerStore) SelectAll(ctx context.Context) (result []*model.ServerModel, err error) {
	result, err = s.FetchAll(ctx)
	return
}

func (s *ServerStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return nil, nil
}
