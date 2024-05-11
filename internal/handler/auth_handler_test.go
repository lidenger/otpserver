package handler

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/lidenger/otpserver/pkg/crypt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func genTimeToken(key, iv string) (string, error) {
	// 使用接入服务的密钥和IV生成time token
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	// pkcs7 padding: https://github.com/lidenger/cryptology/tree/main/padding/pkcs7
	padData := crypt.Pad([]byte(now), aes.BlockSize)
	cipherText := make([]byte, len(padData))
	mode := cipher.NewCBCEncrypter(block, []byte(iv))
	mode.CryptBlocks(cipherText, padData)
	timeToken := hex.EncodeToString(cipherText)
	return timeToken, nil
}

func TestVerifyTimeToken(t *testing.T) {
	// 服务密钥和IV
	key := "0c8441ba0ec011efbb1e2cf05daf3fe5"
	iv := "0c8441ba0ec011ef"

	timeToken, err := genTimeToken(key, iv)
	if err != nil {
		t.Fatal(err)
	}
	msg := verifyTimeTokenInner(timeToken, key, iv, 1)
	if len(msg) == 0 {
		t.Fatal("time token验证成功")
	} else {
		t.Fatal("time token验证失败")
	}
}

type Result struct {
	Code int    `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func TestRequestAccessToken(t *testing.T) {
	url := "http://127.0.0.1:8066/v1/access-token"
	// 服务密钥和IV
	key := "0c8441ba0ec011efbb1e2cf05daf3fe5"
	iv := "0c8441ba0ec011ef"
	accessToken, err := getAccessToken(url, key, iv)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("access_token: %s", accessToken)
	// 验证access token
	verifyUrl := "http://127.0.0.1:8066/v1/access-token/verify"
	err = verifyAccessToken(verifyUrl, accessToken)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("access_token验证成功: %s", accessToken)
	}
}

func getAccessToken(url, key, iv string) (string, error) {
	timeToken, err := genTimeToken(key, iv)
	if err != nil {
		return "", err
	}
	params := make(map[string]string)
	params["serverSign"] = "server1"
	params["timeToken"] = timeToken

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonParams))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	result := &Result{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return "", err
	}
	if result.Code != 200000 {
		return "", errors.New(result.Msg)
	} else {
		return result.Data, nil
	}
}

func verifyAccessToken(baseUrl, accessToken string) error {
	// 验证access token
	fullUrl := baseUrl + "?accessToken=" + accessToken
	resp, err := http.Get(fullUrl)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	result := &Result{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	if result.Code != 200000 {
		return errors.New(result.Msg)
	} else {
		return nil
	}
}
