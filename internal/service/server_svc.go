package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type ServerSvc struct {
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
	m.IsEnable = p.IsEnable
	// 默认启用
	if m.IsEnable == 0 {
		m.IsEnable = 1
	}
	// 默认不启用操作敏感信息
	if m.IsOperateSensitiveData == 0 {
		m.IsOperateSensitiveData = 2
	}
	// 服务密钥
	secret := util.Generate32Str()
	var err error
	secretCipher, err := crypt.Encrypt(cmd.P.RootKey192, cmd.P.IV, []byte(secret))
	if err != nil {
		return nil, err
	}
	m.Secret = secretCipher
	// 服务密钥IV
	iv := util.Generate16Str()
	ivCipher, err := crypt.Encrypt(cmd.P.RootKey192, cmd.P.IV, []byte(iv))
	if err != nil {
		return nil, err
	}
	m.IV = ivCipher
	// 计算数据摘要
	m.DataCheck = s.CalcDataCheckSum(m)
	return m, nil
}

// CalcDataCheckSum 计算数据校验和
func (s *ServerSvc) CalcDataCheckSum(m *model.ServerModel) string {
	data := fmt.Sprintf("%d,%d,%s,%s,%s", m.IsEnable, m.IsOperateSensitiveData, m.Sign, m.Secret, m.IV)
	return crypt.HmacDigest(cmd.P.RootKey192, data)
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
	return util.GetArrFirstItem(data), nil
}

// CheckModel 校验数据,解密账号密钥密文
func (s *ServerSvc) CheckModel(ctx context.Context, m *model.ServerModel) error {
	check := s.CalcDataCheckSum(m)
	if m.DataCheck != check {
		msg := fmt.Sprintf("服务数据校验不通过,疑似被篡改,请关注(ID:%d,sign:%s)", m.ID, m.Sign)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	// 服务密钥
	secret, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.Secret)
	if err != nil {
		return otperr.ErrDecrypt(err)
	}
	m.Secret = string(secret)
	// 服务密钥IV
	iv, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.IV)
	if err != nil {
		return otperr.ErrDecrypt(err)
	}
	m.IV = string(iv)
	return nil
}
