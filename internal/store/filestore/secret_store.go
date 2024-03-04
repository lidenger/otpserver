package filestore

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
)

type SecretStore struct {
}

// 确保SecretStore实现了store.SecretStore
var _ store.SecretStore = (*SecretStore)(nil)

func (s *SecretStore) Insert(ctx context.Context, m *model.AccountSecretModel) error {

	return nil
}

func (s *SecretStore) Update(ctx context.Context, m *model.AccountSecretModel) error {

	return nil
}

func (s *SecretStore) Paging(ctx context.Context, param *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {

	return nil, 0, nil
}

func (s *SecretStore) SelectByCondition(ctx context.Context, condition *param.SecretParam) ([]*model.AccountSecretModel, error) {

	return nil, nil
}

func (s *SecretStore) Delete(ctx context.Context, ID int64) error {

	return nil
}
