package service

import (
	"encoding/base32"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type Crypt struct {
	RootKey []byte // 根密钥
	IV      []byte
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

type doubleWriteFunc func() (store.Tx, error)

// DoubleWrite 主备双写
func DoubleWrite(exec, execBackup doubleWriteFunc) error {
	// 主存储
	tx, err := exec()
	if err != nil {
		tx.Rollback()
		return otperr.ErrStore(err)
	}
	// 没有备存储直接提交事务
	if execBackup == nil {
		tx.Commit()
		return nil
	}
	// 备存储
	tx2, err2 := execBackup()
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
