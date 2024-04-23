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

type SecretStore struct {
	Store    store.SecretStore
	RootPath string
	err      error
}

const SecretFileName = "otp_account_secret"

// 确保SecretStore实现了store.SecretStore
var _ store.SecretStore = (*SecretStore)(nil)
var _ store.LocalStore[*model.AccountSecretModel] = (*SecretStore)(nil)

func (s *SecretStore) GetStoreErr() error {
	return s.err
}

func (s *SecretStore) SetStoreErr(err error) {
	s.err = err
}

func (s *SecretStore) GetStoreType() string {
	return enum.LocalStore
}

func (s *SecretStore) LoadAll(ctx context.Context) error {
	err := s.Store.GetStoreErr()
	if err != nil {
		log.Error("账号密钥主存储异常，无法从主存储加载数据到本地存储，使用上个版本的本地存储数据")
		return s.Store.GetStoreErr()
	}
	err = fetchAllStoreDataAndWriteToFile[*model.AccountSecretModel](ctx, s.Store, s.RootPath, SecretFileName)
	return nil
}

func (s *SecretStore) FetchAll(ctx context.Context) (result []*model.AccountSecretModel, err error) {
	filePath := filepath.Join(s.RootPath, SecretFileName)
	result, err = fetchAll[*model.AccountSecretModel](ctx, filePath)
	return
}

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (store.Tx, error) {
	return store.EmptyTxIns, nil
}

func (s *SecretStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {
	return store.EmptyTxIns, nil
}

func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	if param.PageNo != 1 {
		return nil, 0, nil
	}
	data, err := s.FetchAll(ctx)
	return data, 0, err
}

func (s *SecretStore) SelectById(ctx context.Context, ID int64) (*model.AccountSecretModel, error) {
	return nil, nil
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) ([]*model.AccountSecretModel, error) {
	return nil, nil
}

func (s *SecretStore) SelectAll(ctx context.Context) (result []*model.AccountSecretModel, err error) {
	result, err = s.FetchAll(ctx)
	return
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {
	return store.EmptyTxIns, nil
}
