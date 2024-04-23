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

type ServerIpListStore struct {
	Store    store.ServerIpListStore
	RootPath string
	err      error
}

const ServerIpListFileName = "otp_server_ip_whitelist"

// 确保ServerStore实现了store.ServerStore
var _ store.ServerIpListStore = (*ServerIpListStore)(nil)
var _ store.LocalStore[*model.ServerIpListModel] = (*ServerIpListStore)(nil)

func (s *ServerIpListStore) GetStoreType() string {
	return enum.LocalStore
}

func (s *ServerIpListStore) GetStoreErr() error {
	return s.err
}

func (s *ServerIpListStore) SetStoreErr(err error) {
	s.err = err
}

func (s *ServerIpListStore) LoadAll(ctx context.Context) error {
	err := s.Store.GetStoreErr()
	if err != nil {
		log.Error("接入服务IP白名单列表主存储异常，无法从主存储加载数据到本地存储，使用上个版本的本地存储数据")
		return s.Store.GetStoreErr()
	}
	err = fetchAllStoreDataAndWriteToFile[*model.ServerIpListModel](ctx, s.Store, s.RootPath, ServerIpListFileName)
	return nil
}

func (s *ServerIpListStore) FetchAll(ctx context.Context) (result []*model.ServerIpListModel, err error) {
	filePath := filepath.Join(s.RootPath, ServerIpListFileName)
	result, err = fetchAll[*model.ServerIpListModel](ctx, filePath)
	return
}

func (s *ServerIpListStore) Insert(ctx context.Context, m *model.ServerIpListModel) (store.Tx, error) {
	return nil, nil
}

func (s *ServerIpListStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return nil, nil
}

func (s *ServerIpListStore) Paging(ctx context.Context, param *param.ServerIpListParam) (result []*model.ServerIpListModel, count int64, err error) {
	return nil, 0, nil
}

func (s *ServerIpListStore) SelectById(ctx context.Context, ID int64) (*model.ServerIpListModel, error) {
	return nil, nil
}

func (s *ServerIpListStore) SelectByCondition(ctx context.Context, condition *param.ServerIpListParam) ([]*model.ServerIpListModel, error) {
	return nil, nil
}

func (s *ServerIpListStore) SelectAll(ctx context.Context) (result []*model.ServerIpListModel, err error) {
	result, err = s.FetchAll(ctx)
	return
}

func (s *ServerIpListStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return nil, nil
}
