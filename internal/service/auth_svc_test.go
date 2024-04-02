package service

import (
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/pkg/crypt"
	"strconv"
	"testing"
	"time"
)

var (
	// 长度可选：16|24|32 分别对应AES 128|192|256
	key128 = []byte("1234567890123456")
	key192 = []byte("123456789012345612345678")
	key256 = []byte("12345678901234561234567890123456")
	// 注意长度需要是aes分组块大小
	iv = []byte("1234567890123456")
)

func initP() {
	cmd.P = &cmd.Param{}
	crypto := &cmd.Crypt{}
	crypto.RootKey128 = key128
	crypto.RootKey192 = key192
	crypto.RootKey256 = key256
	crypto.IV = iv
	cmd.P.Crypt = crypto
}

func TestGenTimeToken(t *testing.T) {
	key192 = []byte("74137809f09b11eeb9fe2cf05daf3fe5")
	iv = []byte("74154b1ef09b11ee")
	token, err := crypt.Encrypt(key192, iv, []byte(strconv.FormatInt(time.Now().Unix(), 10)))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token:%+v\n", token)
}

func TestGenAccessToken(t *testing.T) {
	initP()
	token, m, err := GenAccessToken("server1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", m)
	t.Logf("%+v\n", token)
}

func TestAnalysisAccessToken(t *testing.T) {
	initP()
	token := "94fc900c4369fba3735be9a06648c781af8d246a10fb3ddf38e3d6da7eaa3326f4a32cdf9be8a978404f33db890fd083a7318098b8f2a1792cafa6b227ba9400f98b7a4c5500c0344240c8e945e07fb5"
	m, err := AnalysisAccessToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}
