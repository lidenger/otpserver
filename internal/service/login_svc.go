package service

import (
	"encoding/json"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
	"time"
)

// GenLoginToken 生成Login token
func GenLoginToken(account string) (string, error) {
	m := &model.AdminToken{
		Account: account,
		Nonce:   util.Generate32Str(),
		Time:    time.Now().Unix(),
	}
	data := fmt.Sprintf("account:%s;nonce:%s;time:%d", account, m.Nonce, m.Time)
	dataCheck := crypt.HmacDigest(cmd.P.RootKey128, data)
	m.CheckData = dataCheck
	t, err := json.Marshal(m)
	if err != nil {
		return "", otperr.ErrServer(err)
	}
	token, err := crypt.Encrypt(cmd.P.RootKey192, cmd.P.IV, t)
	if err != nil {
		return "", otperr.ErrServer(err)
	}
	return token, nil
}

// AnalysisLoginToken 解析Login token
func AnalysisLoginToken(token string, maxValidMill int64) (*model.AdminToken, error) {
	// 解密获取明文
	tokenJsonBytes, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, token)
	if err != nil || len(tokenJsonBytes) == 0 {
		log.Warn("解密token失败:%s,err:%s", string(tokenJsonBytes), err)
		return nil, otperr.ErrParamIllegal("无效的token[1]")
	}
	// 解析为结构体
	m := &model.AdminToken{}
	err = json.Unmarshal(tokenJsonBytes, m)
	if err != nil {
		log.Warn("json.Unmarshal解析token失败:%s,err:", string(tokenJsonBytes), err)
		return nil, otperr.ErrParamIllegal("无效的token[2]")
	}
	// 验证摘要
	data := fmt.Sprintf("account:%s;nonce:%s;time:%d", m.Account, m.Nonce, m.Time)
	dataCheck := crypt.HmacDigest(cmd.P.RootKey128, data)
	if m.CheckData != dataCheck {
		log.Warn("比对数据摘要不匹配:%s,err:%s", string(tokenJsonBytes), err)
		return nil, otperr.ErrParamIllegal("无效的token[3]")
	}
	// 验证有效期
	if m.Time+maxValidMill < time.Now().Unix() {
		log.Warn("token已过期:%s", string(tokenJsonBytes))
		return nil, otperr.ErrParamIllegal("无效的token[4]")
	}
	return m, nil
}
