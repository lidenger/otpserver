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

1. 人工审批并在管理平台创建server1服务 <br>
   <img src="/doc/res/add_server.png" alt="增加server1服务"/> <br>
2. 增加server1 IP白名单 <br>
   <img src="/doc/res/add_server_iplist.png" alt="增加server1服务IP白名单"/> <br>

### 外部服务接入方

#### Go语言接入

1. 获取依赖

```shell
go get github.com/lidenger/otp-sdk-go
```

2. 配置server参数

```go
// otp server地址
address := "http://127.0.0.1:8066"
// 服务标识
serverSign := "server1"
// 服务密钥
serverKey := "0c8441ba0ec011efbb1e2cf05daf3fe5"
serverIV := "0c8441ba0ec011ef"
otpsdk.Conf(serverSign, address, serverKey, serverIV)
```

3. 开始使用

```go
// 验证动态令牌
success, err := otpsdk.VerifyCode("liweiyi2", "431775")

// 获取密钥二维码内容，可根据内容生成二维码，app扫描后直接添加
content, err := otpsdk.GetQRCodeContent("liweiyi2")

// 增加账号密钥
err := otpsdk.AddAccountSecret("liweiyi2")

// 获取账号密钥
secret, err := otpsdk.GetAccountSecret("liweiyi2")
```

<br>

#### 其他语言接入

请求otp server接口时需要在header中携带access token参数

##### 计算timeToken

```
aes.encrypt("CBC", key, iv, hex.encoding(now()))
```

1. 取当前UTC时间（从1970年到当前时间的毫秒值）
2. 使用Hex编码为字符串
3. 使用aes cbc加密字符串（key和iv在生成服务时已创建），密文为time token

<br/>

##### 获取access token

```http
POST /v1/access-token
```

参数示例:

```json
{
  "serverSign": "server1",
  "timeToken": "96d5286e03fc60a1bf3faecd921e250c"
}
```

返回示例：

```json
{
  "code": 200000,
  "data": "2edefa52af1e848c56a2749d25c653acc68607ffeba7a1420c17a10f12de1909a37fb871030935fb9114a573e3dace044ef5d582325c7ad0ef04eef40c2ce2313100d0a3adbd6872f09e5d9e769f1fed",
  "msg": "success"
}
```

<br>

##### 新增账号密钥

```http
POST /v1/secret

## header
Content-Type : application/json
Authorization : 2edefa52af1e848c56a2749d25c653acc68607ffeba7a1420c17a10f12de1909a37fb871030935fb9114a573e3dace044ef5d582325c7ad0ef04eef40c2ce2313100d0a3adbd6872f09e5d9e769f1fed

```

参数示例:

```json
{
  "account": "liweiyi3",
  "isEnable": 1
}
```

返回示例：

```json
{
  "code": 200000,
  "data": "",
  "msg": "success"
}
```

<br>

##### 获取账号密钥

GET /v1/secret/{account}

```http
GET /v1/secret/liweiyi3

## header
Authorization : 2edefa52af1e848c56a2749d25c653acc68607ffeba7a1420c17a10f12de1909a37fb871030935fb9114a573e3dace044ef5d582325c7ad0ef04eef40c2ce2313100d0a3adbd6872f09e5d9e769f1fed

```

返回示例：

```json
{
  "code": 200000,
  "data": {
    "id": 30,
    "secret": "GFQWCYRXMQYDEMJQMVTDCMLFMY4DGZJTGJRWMMBVMRQWMM3GMU2Q",
    "account": "liweiyi3",
    "dataCheck": "9b2737693d57fd538c756c6fa473766a31c6c01055631b0d76460bb8f0ef1c7c",
    "isEnable": 1,
    "createTime": "2024-05-13T14:07:40+08:00",
    "updateTime": "2024-05-13T14:07:40+08:00"
  },
  "msg": "success"
}
```

<br>

##### 验证动态令牌

GET /v1/secret/valid?account={account}&code={code}

```http
GET /v1/secret/valid?account=liweiyi3&code=431775

## header
Authorization : 2edefa52af1e848c56a2749d25c653acc68607ffeba7a1420c17a10f12de1909a37fb871030935fb9114a573e3dace044ef5d582325c7ad0ef04eef40c2ce2313100d0a3adbd6872f09e5d9e769f1fed

```

返回示例：<br>
data: true验证通过，false验证失败

```json
{
  "code": 200000,
  "data": true,
  "msg": "success"
}
```

<br>

##### 获取密钥二维码内容

可根据内容生成二维码，app扫描后直接添加

GET /v1/secret/{account}/qrcode-content

```http
GET /v1/secret/liweiyi2/qrcode-content

## header
Authorization : 2edefa52af1e848c56a2749d25c653acc68607ffeba7a1420c17a10f12de1909a37fb871030935fb9114a573e3dace044ef5d582325c7ad0ef04eef40c2ce2313100d0a3adbd6872f09e5d9e769f1fed

```

返回示例：<br>
data: true验证通过，false验证失败

```json
{
  "code": 200000,
  "data": "otpauth://totp/liweiyi2?secret=HAZWIZTEGRQTGMJQMY4TCMLFMZQWCZJRGJRWMMBVMRQWMM3GMU2Q====\u0026issuer=otpserver",
  "msg": "success"
}
```


