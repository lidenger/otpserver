## 外部服务接入

### 外部接入服务规划
```text
服务标识: server1
服务名：接入服务1
服务描述：接入服务1描述
敏感信息：不涉及
IP白名单：启用
IP白名单列表：10.0.0.1,10.0.0.2
```

### otp server admin管理员
1. 人工审批并在管理平台创建server1服务
<img src="/doc/res/add_server.png" alt="增加server1服务"/>
2. 增加server1 IP白名单
<img src="/doc/res/add_server_iplist.png" alt="增加server1服务IP白名单"/>

### 外部服务接入方
Go语言接入
```go
// 服务密钥和IV
key := "0c8441ba0ec011efbb1e2cf05daf3fe5"
iv := "0c8441ba0ec011ef"

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

// 生成time token
timeToken := genTimeToken(key, iv)


```
注: pkcs7 padding: https://github.com/lidenger/cryptology/tree/main/padding/pkcs7


Java语言接入
```java

```

