package service

import (
	"encoding/json"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
	"time"
)

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
