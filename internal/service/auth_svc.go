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

var accessTokenCache1 = cache.New(3*time.Minute, 5*time.Minute)

var accessTokenCache *cache.Cache

func addTokenCache() {
	if accessTokenCache == nil {
		conf := serverconf.GetSysConf()

		accessTokenCache = cache.New(time.Duration(3*time.Minute), 5*time.Minute)
	}
}

// GenAccessToken 生成access token
// serverSign 接入服务的标识
// key,iv 系统根密钥，防伪造
func GenAccessToken(serverSign string) (accessToken string, err error) {
	accessTokenModel := &model.AccessToken{}
	accessTokenModel.Sign = serverSign
	accessTokenModel.CreateTime = time.Now().Unix()
	accessTokenModel.Rn = util.Generate32Str()
	tokenJson, err := json.Marshal(accessTokenModel)
	if err != nil {
		return "", err
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
