package crypt

import "testing"

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
