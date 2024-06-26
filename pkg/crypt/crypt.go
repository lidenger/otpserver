package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// HmacDigest hmac数据摘要
func HmacDigest(key []byte, param any) string {
	h := hmac.New(sha256.New, key)
	var data []byte
	if t1, ok := param.(string); ok {
		data = []byte(t1)
	} else if t2, ok2 := param.([]byte); ok2 {
		data = t2
	}
	h.Write(data)
	digested := h.Sum(nil)
	return hex.EncodeToString(digested)
}

// Encrypt AES/CBC/PKCS#7
func Encrypt(key, iv, data []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	padData := Pad(data, aes.BlockSize)
	cipherText := make([]byte, len(padData))
	mode.CryptBlocks(cipherText, padData)
	return hex.EncodeToString(cipherText), nil
}

// Decrypt AES/CBC/PKCS#7
func Decrypt(key, iv []byte, cipherText string) ([]byte, error) {
	cipherBytes, err := hex.DecodeString(cipherText)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	data := make([]byte, len(cipherBytes))
	mode.CryptBlocks(data, cipherBytes)
	originData, err := Unpad(data)
	if err != nil {
		return nil, err
	}
	return originData, nil
}
