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

func initPFromLocalAppKeyFile() {
	crypto := cmd.AnalysisKeyFile("H:/lidenger/otpserver/app.key")
	if cmd.P == nil {
		cmd.P = &cmd.Param{}
	}
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
	initPFromLocalAppKeyFile()
	token, m, err := GenAccessToken("server1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", m)
	t.Logf("%+v\n", token)
}

func TestAnalysisAccessToken(t *testing.T) {
	initPFromLocalAppKeyFile()
	token := "2edefa52af1e848c56a2749d25c653ac6dd55ec5748b259bf66c39b5f0cda70e5cb95d3b27a161a14efc1b5ddd265a88cd367927365aa6a78b2ccb48dcaee2c3ec4d0ad56cfb01ae3acffbf4acfad4e1"
	m, err := AnalysisAccessToken(token)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", m)
}
