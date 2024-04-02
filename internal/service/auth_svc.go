package service

import (
	"encoding/json"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
	"github.com/patrickmn/go-cache"
	"time"
)

var accessTokenCache *cache.Cache

// AddAccessTokenCache 增加access token到缓存
func AddAccessTokenCache(accessToken string, m *model.AccessToken) {
	if accessTokenCache == nil {
		conf := serverconf.GetSysConf()
		minute := time.Duration(conf.Server.TimeTokenValidMinute)
		accessTokenCache = cache.New(minute*time.Minute, minute*time.Minute)
	}
	accessTokenCache.SetDefault(accessToken, m)
}

// GetAccessTokenInCache 在缓存中获取access token
func GetAccessTokenInCache(accessToken string) *model.AccessToken {
	if accessTokenCache == nil {
		return nil
	}
	m, exists := accessTokenCache.Get(accessToken)
	if !exists {
		return nil
	}
	if x, ok := m.(*model.AccessToken); ok {
		return x
	} else {
		return nil
	}
}

// GenAccessToken 生成access token
// serverSign 接入服务的标识
// key,iv 系统根密钥，防伪造
func GenAccessToken(serverSign string) (accessToken string, m *model.AccessToken, err error) {
	m = &model.AccessToken{}
	m.Sign = serverSign
	m.CreateTime = time.Now().Unix()
	m.Rn = util.Generate32Str()
	tokenJson, err := json.Marshal(m)
	if err != nil {
		return "", nil, err
	}
	accessToken, err = crypt.Encrypt(cmd.P.RootKey128, cmd.P.IV, tokenJson)
	return
}

// AnalysisAccessToken 解析access token
func AnalysisAccessToken(accessToken string) (*model.AccessToken, error) {
	tokenJson, err := crypt.Decrypt(cmd.P.RootKey128, cmd.P.IV, accessToken)
	if err != nil {
		return nil, err
	}
	accessTokenModel := &model.AccessToken{}
	err = json.Unmarshal(tokenJson, accessTokenModel)
	if err != nil {
		return nil, err
	}
	if len(accessTokenModel.Sign) == 0 ||
		len(accessTokenModel.Rn) == 0 ||
		accessTokenModel.CreateTime == 0 {
		return nil, otperr.ErrParamIllegal("access token不正确")
	}
	return accessTokenModel, nil
}

// VerifyAccessTokenM 验证access token
// 验证有效期
func VerifyAccessTokenM(tokenM *model.AccessToken) error {
	diff := time.Now().Unix() - tokenM.CreateTime
	conf := serverconf.GetSysConf()
	validHour := int64(conf.Server.AccessTokenValidHour)
	if diff > validHour*3600 {
		msg := fmt.Sprintf("access token已过期,有效期:%d", validHour)
		return otperr.ErrParamIllegal(msg)
	}
	return nil
}
