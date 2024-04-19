package localstore

import (
	"context"
	"encoding/json"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/enum"
	"os"
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
	filePath := filepath.Join(s.RootPath, SecretFileName)
	err = manageLocalStoreFile(filePath)
	if err != nil {
		return err
	}
	p := &param.SecretPagingParam{}
	p.PageNo = 1
	p.PageSize = 100
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	secretArr := make([]*model.AccountSecretModel, 0)
	for {
		data, _, err := s.Store.Paging(ctx, p)
		if err != nil {
			return err
		}
		// 获取了所有数据
		if len(data) == 0 {
			break
		}
		for _, m := range data {
			secretArr = append(secretArr, m)
		}
		p.PageNo++
	}
	js, err := json.Marshal(secretArr)
	if err != nil {
		return err
	}
	_, err = file.Write(js)
	if err != nil {
		return err
	}
	return nil
}

func (s *SecretStore) FetchAll(_ context.Context) (result []*model.AccountSecretModel, err error) {
	filePath := filepath.Join(s.RootPath, SecretFileName)
	js, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(js, &result)
	if err != nil {
		return nil, err
	}
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
