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

func (s *ServerStore) LoadAll(ctx context.Context) error {
	err := s.Store.GetStoreErr()
	if err != nil {
		log.Error("接入服务主存储异常，无法从主存储加载数据到本地存储，使用上个版本的本地存储数据")
		return s.Store.GetStoreErr()
	}
	filePath := filepath.Join(s.RootPath, ServerFileName)
	err = manageLocalStoreFile(filePath)
	if err != nil {
		return err
	}
	p := &param.ServerPagingParam{}
	p.PageNo = 1
	p.PageSize = 100
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	secretArr := make([]*model.ServerModel, 0)
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

func (s *ServerStore) FetchAll(ctx context.Context) ([]*model.ServerModel, error) {
	return nil, nil
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
