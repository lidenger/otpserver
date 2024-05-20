# otpserver

## 简述

OTP : One-Time Password 一次性密码<br/>
otpserver围绕着一次性密码功能搭建的完整的服务，特点：
- 通过简单的配置即可运行
- 高数据安全，配置文件中的密码等敏感信息，数据库中的密钥都做了加密保护，数据库中的数据也做了防篡改防护
- 高可用/主动降级，主存储 -> 备存储 -> 本地存储 -> 内存存储，共建立4道防线，主动摘除和上线store<br/>
- 多模式启动，支持本地配置文件和配置中心启动服务

## 系统架构
<img src="/doc/res/system.png" alt="otpserver系统架构图"/>

## 开始使用
- [快速开始](doc%2Fquick-start.md)
- [配置详情](doc%2Fconfig-detail.md)

#### 启动方式
- 启动方式1：本地启动
```shell
# 方式一 （默认本地启动）
./app

# 方式二 （本地启动并指定配置文件路径）
./app --confSource local --confFile /etc/conf/app.toml
```

- 启动方式2：配置中心Nacos启动
1. 将app.toml文件内容配置到Nacos Server端
<img src="/doc/res/nacos_conf.png" alt="Nacos Server配置"/>
2. Nacos配置文件
```toml
[client]
namespaceId = ""
dataId = "otpserverv1"
group = "DEFAULT_GROUP"
timeoutMs = 5000
logDir = "/nacos/log/"
cacheDir = "/nacos/cache/"
logLevel = "info"
isListenChange = true
# server1
[[server]]
ip = "127.0.0.1"
port = 8848
contextPath = "/nacos"

# server 2
#[[server]]
#ip = "127.0.0.1"
#port = 8848
#contextPath = "/nacos"
```

3. 从nacos启动服务
```shell
./app --confSource nacos
```

## 外部服务接入
[接入详情](doc%2Fexternal-server.md)


## 监控-Prometheus
- otp_server_http_req_total  : http请求总次数
- otp_server_http_limit_req_total : http请求超限总次数
- otp_server_http_req_cost_gauge : http请求耗时，单位：毫秒
- otp_server_http_cost_histogram : http请求耗时分布，单位：毫秒
- otp_server_http_response_code_histogram : http响应状态码
- otp_server_store_mysql_health_gauge : mysql存储健康状态
- otp_server_store_pgsql_health_gauge : pgsql存储健康状态
- otp_server_store_local_health_gauge : 本地存储健康状态
- otp_server_store_memory_health_gauge : 内存存储健康状态






