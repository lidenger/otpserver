package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"time"
)

type ServerSvc struct {
	Crypt
	Store       store.ServerStore // 主存储
	StoreBackup store.ServerStore // 备存储
}

// Add 增加接入的服务
func (s *ServerSvc) Add(ctx context.Context, p *param.ServerParam) error {
	exists, err := s.IsExists(ctx, p.Sign)
	if err != nil {
		return err
	}
	if exists {
		msg := fmt.Sprintf("服务%s已存在不能重复添加", p.Sign)
		return otperr.ErrRepeatAdd(errors.New(msg))
	}
	// 创建一个新的model
	m, err := s.NewServerModel(p)
	if err != nil {
		return err
	}
	// 主备写
	var backupExec doubleWriteFunc = nil
	if s.StoreBackup != nil {
		backupExec = func() (store.Tx, error) {
			return s.StoreBackup.Insert(ctx, m)
		}
	}
	err = DoubleWrite(func() (store.Tx, error) {
		return s.Store.Insert(ctx, m)
	}, backupExec)
	return err
}

func (s *ServerSvc) NewServerModel(p *param.ServerParam) (*model.ServerModel, error) {
	m := &model.ServerModel{}
	m.Name = p.Name
	m.Sign = p.Sign
	m.Remark = p.Remark
	if p.IsEnable == 0 {
		// 默认启用
		m.IsEnable = 1
	} else {
		m.IsEnable = p.IsEnable
	}
	n := time.Now()
	m.CreateTime = n
	m.UpdateTime = n
	var err error
	m.Secret, err = genSecret(s.RootKey, s.IV)
	if err != nil {
		return nil, err
	}
	// 计算数据摘要
	m.DataCheck = s.CalcDataCheckSum(m.IsEnable, m.Sign, m.Secret)
	return m, nil
}

// CalcDataCheckSum 计算数据校验和
func (s *ServerSvc) CalcDataCheckSum(isEnable uint8, sign, secretCipher string) string {
	data := fmt.Sprintf("%d,%s,%s", isEnable, sign, secretCipher)
	return crypt.HmacDigest(s.RootKey, data)
}

// IsExists 账号密钥是否存在
func (s *ServerSvc) IsExists(ctx context.Context, sign string) (bool, error) {
	secretData, err := s.GetBySign(ctx, sign)
	if err != nil {
		return false, err
	}
	return secretData != nil, nil
}

func (s *ServerSvc) GetBySign(ctx context.Context, sign string) (*model.ServerModel, error) {
	var err error
	var serverModel *model.ServerModel
	serverModel, err = findServerByStore(ctx, sign, s.Store)
	if err != nil {
		if s.StoreBackup == nil {
			return nil, otperr.ErrStore(err)
		}
		log.Warn("主存储获取服务信息异常,尝试从备存储获取,主存储异常信息:%+v", err)
		var errBackup error
		serverModel, errBackup = findServerByStore(ctx, sign, s.StoreBackup)
		if errBackup != nil {
			log.Error("主备存储都获取失败,主存储err:%+v,备存储err:%+v", err, errBackup)
			return nil, otperr.ErrStoreBackup(errBackup)
		}
	}
	if serverModel == nil {
		return nil, nil
	}
	err = s.CheckModel(ctx, serverModel)
	if err != nil {
		return nil, err
	}
	return serverModel, err
}

func findServerByStore(ctx context.Context, sign string, s store.ServerStore) (*model.ServerModel, error) {
	condition := &param.ServerParam{}
	condition.Sign = sign
	data, err := s.SelectByCondition(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(data) > 0 {
		return data[0], nil
	} else {
		return nil, nil
	}
}

// CheckModel 校验数据,解密账号密钥密文
func (s *ServerSvc) CheckModel(ctx context.Context, m *model.ServerModel) error {
	check := s.CalcDataCheckSum(m.IsEnable, m.Sign, m.Secret)
	if m.DataCheck != check {
		msg := fmt.Sprintf("服务[%s]数据[ID:%d]校验不通过，请关注", m.Sign, m.ID)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	secret, err := crypt.Decrypt(s.RootKey, s.IV, m.Secret)
	if err != nil {
		return otperr.ErrDecrypt(err)
	}
	m.Secret = string(secret)
	return nil
}
