package crypt

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"testing"
)

var (
	key128 = []byte("1234567890123456")
	key192 = []byte("123456789012345612345678")
	key256 = []byte("12345678901234561234567890123456")
	iv     = []byte("1234567890123456")
)

func TestEncryptAndDecrypt128(t *testing.T) {
	data := "123abc加解密测试&#!22"
	cipher, err := Encrypt(key128, iv, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	origin, err := Decrypt(key128, iv, cipher)
	if err != nil {
		t.Fatal(err)
	}
	if string(origin) != data {
		t.Fatalf("加解密测试失败，期望:%s,实际:%s", data, string(origin))
	}
}

func TestEncryptAndDecrypt192(t *testing.T) {
	data := "123abc加解密测试&#!22"
	cipher, err := Encrypt(key192, iv, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	origin, err := Decrypt(key192, iv, cipher)
	if err != nil {
		t.Fatal(err)
	}
	if string(origin) != data {
		t.Fatalf("加解密测试失败，期望:%s,实际:%s", data, string(origin))
	}
}

func TestEncryptAndDecrypt256(t *testing.T) {
	data := "123abc加解密测试&#!22"
	cipher, err := Encrypt(key256, iv, []byte(data))
	if err != nil {
		t.Fatal(err)
	}
	origin, err := Decrypt(key256, iv, cipher)
	if err != nil {
		t.Fatal(err)
	}
	if string(origin) != data {
		t.Fatalf("加解密测试失败，期望:%s,实际:%s", data, string(origin))
	}
}

func TestHmacDigest(t *testing.T) {
	files := []string{"otpserver-linux-amd64.zip", "otpserver-macos-amd64.zip", "otpserver-macos-arm64.zip", "otpserver-windows-amd64.zip"}
	for _, file := range files {
		content, err := os.ReadFile("../../doc/download/" + file)
		if err != nil {
			t.Fatal(err)
		}

		digested := sha256.Sum256(content)
		hexStr := hex.EncodeToString(digested[:])
		size := float64(len(content)) / 1024 / 1024
		fmt.Printf("%s:%.1fMB | %s\n", file, size, hexStr)
	}
}

type downloadM struct {
	filename string
	os       string
	cpu      string
	size     string
	sum      string
}

func TestGenDownloadContent(t *testing.T) {
	itmes := []downloadM{
		{"otpserver-linux-amd64.zip", "Linux", "x86-64", "", ""},
		{"otpserver-macos-amd64.zip", "macOS", "x86-64", "", ""},
		{"otpserver-macos-arm64.zip", "macOS", "M系列", "", ""},
		{"otpserver-windows-amd64.zip", "Windows", "x86-64", "", ""},
	}
	for _, item := range itmes {
		content, err := os.ReadFile("../../doc/download/" + item.filename)
		if err != nil {
			t.Fatal(err)
		}
		digested := sha256.Sum256(content)
		item.sum = hex.EncodeToString(digested[:])
		item.size = fmt.Sprintf("%.1fMB", float64(len(content))/1024/1024)
		fmt.Printf("| [%s](download/%s)     | %s   | %s | %s | %s | \n", item.filename, item.filename, item.os, item.cpu, item.size, item.sum)
	}
}
