package service

import (
	"context"
	"encoding/base32"
	"errors"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type SecretSvc struct {
	Store       store.SecretStore // 主存储
	StoreBackup store.SecretStore // 备存储
	StoreMemory store.SecretStore // 内存存储
	StoreLocal  store.SecretStore // 本地存储
}

// Add 添加账号密钥
func (s *SecretSvc) Add(ctx context.Context, account string) error {
	exists, err := s.IsExists(ctx, account)
	if err != nil {
		return err
	}
	if exists {
		msg := fmt.Sprintf("账号%s已存在不能重复添加", account)
		return otperr.ErrRepeatAdd(errors.New(msg))
	}
	// 创建一个新的model
	m, err := s.NewSecretModel(account)
	if err != nil {
		return err
	}
	// 主备写
	err = MultiStoreInsert[*model.AccountSecretModel](ctx, m, s.Store, s.StoreBackup)
	return err
}

// NewSecretModel 创建一个新的账号密钥model
func (s *SecretSvc) NewSecretModel(account string) (*model.AccountSecretModel, error) {
	m := &model.AccountSecretModel{}
	m.Account = account
	m.IsEnable = 1
	// 密钥加密存储
	var err error
	m.SecretSeed, err = genSecret(cmd.P.RootKey192, cmd.P.IV)
	if err != nil {
		return nil, err
	}
	// 计算数据摘要
	m.DataCheck = s.CalcDataCheckSum(m.IsEnable, m.Account, m.SecretSeed)
	return m, nil
}

// 生成密钥
func genSecret(rootKey, iv []byte) (string, error) {
	secret := util.GenerateStr()
	secretEncode := base32.StdEncoding.EncodeToString([]byte(secret))
	secretCipher, err := crypt.Encrypt(rootKey, iv, []byte(secretEncode))
	if err != nil {
		return "", otperr.ErrEncrypt(err)
	}
	return secretCipher, nil
}

// IsExists 账号密钥是否存在
func (s *SecretSvc) IsExists(ctx context.Context, account string) (bool, error) {
	secretData, err := s.GetByAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return secretData != nil, nil
}

// GetByAccount 通过账号获取密钥信息,密钥已解密
func (s *SecretSvc) GetByAccount(ctx context.Context, account string) (*model.AccountSecretModel, error) {
	var err error
	p := &param.SecretParam{Account: account}
	data, err := MultiStoreSelectByCondition[*param.SecretParam, *model.AccountSecretModel](ctx, p, s.StoreMemory, s.Store, s.StoreBackup, s.StoreLocal)
	if err != nil {
		return nil, err
	}
	secretModel := util.GetArrFirstItem(data)
	if secretModel == nil {
		return nil, nil
	}
	err = s.CheckModel(secretModel)
	if err != nil {
		return nil, err
	}
	return secretModel, err
}

// CheckModel 校验数据,解密账号密钥密文
func (s *SecretSvc) CheckModel(m *model.AccountSecretModel) error {
	check := s.CalcDataCheckSum(m.IsEnable, m.Account, m.SecretSeed)
	if m.DataCheck != check {
		msg := fmt.Sprintf("账号密钥数据校验不通过,疑似被篡改,请关注(ID:%d,账号:%s)", m.ID, m.Account)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	secret, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.SecretSeed)
	if err != nil {
		return otperr.ErrDecrypt(err)
	}
	m.SecretSeed = string(secret)
	return nil
}

// CalcDataCheckSum 计算数据校验和
func (s *SecretSvc) CalcDataCheckSum(isEnable uint8, account, secretSeedCipher string) string {
	data := fmt.Sprintf("%d,%s,%s", isEnable, account, secretSeedCipher)
	return crypt.HmacDigest(cmd.P.RootKey192, data)
}

func (s *SecretSvc) Paging(ctx context.Context, p *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	result, count, err = MultiStorePaging[*param.SecretPagingParam, *model.AccountSecretModel](ctx, p, s.StoreMemory, s.Store, s.StoreBackup, s.StoreLocal)
	return
}
