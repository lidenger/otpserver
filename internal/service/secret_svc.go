package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"strings"
	"time"
)

type SecretSvc struct {
	Store       store.SecretStore // 主存储
	StoreBackup store.SecretStore // 备存储
	RootKey     []byte            // 根密钥
	IV          []byte
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
	// 主存储
	tx, err := s.Store.Insert(ctx, m)
	if err != nil {
		tx.Rollback()
		return otperr.ErrStore(err)
	}
	// 没有备存储直接提交事务
	if s.StoreBackup == nil {
		tx.Commit()
		return nil
	}
	// 备存储
	tx2, err2 := s.StoreBackup.Insert(ctx, m)
	if err2 != nil {
		// 为了数据一致性，主存储也需要回滚
		tx.Rollback()
		tx2.Rollback()
		return otperr.ErrStoreBackup(err2)
	}
	// 主备存储成功，提交事务
	tx.Commit()
	tx2.Commit()
	return nil
}

// NewSecretModel 创建一个新的账号密钥model
func (s *SecretSvc) NewSecretModel(account string) (*model.AccountSecretModel, error) {
	m := &model.AccountSecretModel{}
	m.Account = account
	// 默认启用
	m.IsEnable = 1
	n := time.Now()
	m.CreateTime = n
	m.UpdateTime = n
	str, _ := uuid.NewUUID()
	m.SecretSeed = strings.ReplaceAll(str.String(), "-", "")
	// 密钥加密存储
	var err error
	m.SecretSeed, err = crypt.Encrypt(s.RootKey, s.IV, []byte(m.SecretSeed))
	if err != nil {
		return m, otperr.ErrEncrypt(err)
	}
	// 计算数据摘要
	m.DataCheck = s.CalcDataCheckSum(m.IsEnable, m.Account, m.SecretSeed)
	return m, nil
}

// IsExists 账号密钥是否存在
func (s *SecretSvc) IsExists(ctx context.Context, account string) (bool, error) {
	secretData, err := s.GetByAccount(ctx, account)
	if err != nil {
		return false, err
	}
	return secretData != nil, nil
}

// GetByAccount 通过账号获取密钥信息
func (s *SecretSvc) GetByAccount(ctx context.Context, account string) (*model.AccountSecretModel, error) {
	var err error
	var secretModel *model.AccountSecretModel
	secretModel, err = findByStore(ctx, account, s.Store)
	if err != nil {
		if s.StoreBackup == nil {
			return nil, otperr.ErrStore(err)
		}
		log.Warn("主存储获取账号密钥信息异常,尝试从备存储获取,主存储异常信息:%+v", err)
		var errBackup error
		secretModel, errBackup = findByStore(ctx, account, s.StoreBackup)
		if errBackup != nil {
			log.Error("主备存储都获取失败,主存储err:%+v,备存储err:%+v", err, errBackup)
			return nil, otperr.ErrStoreBackup(errBackup)
		}
	}
	if secretModel == nil {
		return nil, nil
	}
	err = s.CheckModel(ctx, secretModel)
	if err != nil {
		return nil, err
	}
	return secretModel, err
}

// CheckModel 校验数据,解密账号密钥密文
func (s *SecretSvc) CheckModel(ctx context.Context, m *model.AccountSecretModel) error {
	check := s.CalcDataCheckSum(m.IsEnable, m.Account, m.SecretSeed)
	if m.DataCheck != check {
		msg := fmt.Sprintf("账号[%s]数据被篡改，数据校验不通过，请关注", m.Account)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	secret, err := crypt.Decrypt(s.RootKey, s.IV, m.SecretSeed)
	if err != nil {
		return otperr.ErrDecrypt(err)
	}
	m.SecretSeed = string(secret)
	return nil
}

func findByStore(ctx context.Context, account string, s store.SecretStore) (*model.AccountSecretModel, error) {
	condition := &param.SecretParam{}
	condition.Account = account
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

// CalcDataCheckSum 计算数据校验和
func (s *SecretSvc) CalcDataCheckSum(isEnable uint8, account, secretSeedCipher string) string {
	data := fmt.Sprintf("%d,%s,%s", isEnable, account, secretSeedCipher)
	return crypt.HmacDigest(s.RootKey, data)
}
