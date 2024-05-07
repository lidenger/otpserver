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
	Store                   store.SecretStore // 主存储
	StoreBackup             store.SecretStore // 备存储
	StoreMemory             store.SecretStore // 内存存储
	storeDetectionEventChan chan<- struct{}
}

// Add 添加账号密钥
func (s *SecretSvc) Add(ctx context.Context, account string, isEnable uint8) error {
	exists, err := s.IsExists(ctx, account)
	if err != nil {
		return err
	}
	if exists {
		return otperr.ErrRepeatAdd("账号%s已存在不能重复添加", account)
	}
	// 创建一个新的model
	m, err := s.NewSecretModel(account, isEnable)
	if err != nil {
		return err
	}
	err = MultiStoreInsert[*model.AccountSecretModel](ctx, s.storeDetectionEventChan, m, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}

// NewSecretModel 创建一个新的账号密钥model
func (s *SecretSvc) NewSecretModel(account string, isEnable uint8) (*model.AccountSecretModel, error) {
	m := &model.AccountSecretModel{}
	m.Account = account
	m.IsEnable = isEnable
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
	secretData, err := s.GetByAccount(ctx, account, false)
	if err != nil {
		return false, err
	}
	return secretData != nil, nil
}

// GetByAccount 通过账号获取密钥信息,密钥已解密
// isDecrypt 是否解密账号密钥
func (s *SecretSvc) GetByAccount(ctx context.Context, account string, isDecrypt bool) (*model.AccountSecretModel, error) {
	var err error
	p := &param.SecretParam{Account: account}
	data, err := MultiStoreSelectByCondition[*param.SecretParam, *model.AccountSecretModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return nil, err
	}
	secretModel := util.GetArrFirstItem(data)
	if secretModel == nil {
		return nil, nil
	}
	err = s.CheckModel(secretModel, isDecrypt)
	if err != nil {
		return nil, err
	}
	return secretModel, err
}

// CheckModel 校验数据,解密账号密钥密文
func (s *SecretSvc) CheckModel(m *model.AccountSecretModel, isDecrypt bool) error {
	check := s.CalcDataCheckSum(m.IsEnable, m.Account, m.SecretSeed)
	if m.DataCheck != check {
		msg := fmt.Sprintf("账号密钥数据校验不通过,疑似被篡改,请关注(ID:%d,账号:%s)", m.ID, m.Account)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	if isDecrypt {
		secret, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.SecretSeed)
		if err != nil {
			return otperr.ErrDecrypt(err)
		}
		m.SecretSeed = string(secret)
	}
	return nil
}

// CalcDataCheckSum 计算数据校验和
func (s *SecretSvc) CalcDataCheckSum(isEnable uint8, account, secretSeedCipher string) string {
	data := fmt.Sprintf("%d,%s,%s", isEnable, account, secretSeedCipher)
	return crypt.HmacDigest(cmd.P.RootKey192, data)
}

// Paging 分页
func (s *SecretSvc) Paging(ctx context.Context, p *param.SecretPagingParam) (result []*model.AccountSecretModel, count int64, err error) {
	result, count, err = MultiStorePaging[*param.SecretPagingParam, *model.AccountSecretModel](ctx, s.storeDetectionEventChan, p, s.Store, s.StoreBackup)
	return
}

func (s *SecretSvc) SetEnable(ctx context.Context, account string, isEnable uint8) error {
	m, err := s.GetByAccount(ctx, account, false)
	if err != nil {
		return err
	}
	if m == nil {
		return otperr.ErrParamIllegal("账号不存在:" + account)
	}
	// 数据一致无需更新
	if m.IsEnable == isEnable {
		return nil
	}
	checkSum := s.CalcDataCheckSum(isEnable, account, m.SecretSeed)
	params := make(map[string]any)
	params["is_enable"] = isEnable
	params["data_check"] = checkSum
	err = MultiStoreUpdate(ctx, s.storeDetectionEventChan, m.ID, params, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}
