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

## 外部服务接入
[接入详情](doc%2Fexternal-server.md)







