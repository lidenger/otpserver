package localstore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
)

type SecretStore struct {
	err error
}

// 确保SecretStore实现了store.SecretStore
var _ store.SecretStore = (*SecretStore)(nil)

func (s *SecretStore) GetStoreErr() error {
	return s.err
}

func (s *SecretStore) SetStoreErr(err error) {
	s.err = err
}

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) (store.Tx, error) {

	return nil, nil
}

func (s *SecretStore) Update(ctx context.Context, ID int64, params map[string]any) (store.Tx, error) {

	return nil, nil
}

func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {

	return nil, 0, nil
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) ([]*model.AccountSecretModel, error) {

	return nil, nil
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) (store.Tx, error) {

	return nil, nil
}
