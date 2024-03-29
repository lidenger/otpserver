package service

import (
	"encoding/hex"
	"encoding/json"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/util"
	"time"
)

// GenAccessToken 生成access token
func GenAccessToken(serverSign string, key, iv []byte) (string, error) {
	accessToken := &model.AccessToken{}
	accessToken.Sign = serverSign
	accessToken.Tim = time.Now().Unix()
	accessToken.Rn = util.Generate32Str()
	tokenJson, err := json.Marshal(accessToken)
	if err != nil {
		return "", err
	}
	at := hex.EncodeToString(tokenJson)
	// 使用业务的key加密access token
	atCipher, err := crypt.Encrypt(key, iv, []byte(at))
	return atCipher, err
}
