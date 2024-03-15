package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"strings"
	"time"
)

type SecretSvc struct {
	store   store.SecretStore
	rootKey []byte // 根密钥
	IV      []byte
}

// Add 添加账号密钥
func (s *SecretSvc) Add(ctx context.Context, account string) error {
	m := &model.AccountSecretModel{}
	m.Account = account
	// 默认启用
	m.IsEnable = 1
	n := time.Now()
	m.CreateTime = n
	m.UpdateTime = n
	str, _ := uuid.NewUUID()
	m.SecretSeed = strings.ReplaceAll(str.String(), "-", "")
	var err error = nil
	m.SecretSeed, err = crypt.Encrypt(s.rootKey, s.IV, []byte(m.SecretSeed))

	err = s.store.Insert(ctx, m)
	if err != nil {
		return otperr.ErrStore(err)
	}
	return nil
}

// CalcDataCheckSum 计算数据校验和
func (s *SecretSvc) CalcDataCheckSum(isEnable uint8, account, secretSeedCipher string) string {
	data := fmt.Sprintf("%d,%s,%s", isEnable, account, secretSeedCipher)
	return crypt.HmacDigest(s.rootKey, data)
}
