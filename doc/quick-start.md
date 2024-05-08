## 快速开始使用

### 说明

提供了开箱即用的可执行文件，通过简单配置即可直接线上运行

### 下载

| 文件                                                | 系统      | CPU架构  | 大小     | SHA256 Checksum                                                  |
|:--------------------------------------------------|:--------|:-------|:-------|:-----------------------------------------------------------------|
| [linux-amd64.zip](download%2Flinux-amd64.zip)     | Linux   | x86-64 | 10.8MB | 5f60e7f83fe4a2f9367e51f814ba8440d3bba86ccac5258d4b5dc1babe7fc1ef |
| [macos-amd64.zip](download%2Fmacos-amd64.zip)     | macOS   | x86-64 | 11.3MB | 22cb779336c8a651819464dfdd60fefa45b66d8fea04beb29f17b0397de192b0 |
| [macos-arm64.zip](download%2Fmacos-arm64.zip)     | macOS   | M系列    | 11.0MB | fc2812771ad39d10d4dadbdb3b06a2b98fc3f1ede020528164f1f3e91bd3a652 |
| [windows-amd64.zip](download%2Fwindows-amd64.zip) | Windows | x86-64 | 10.9MB | 20287b63ab997860f8a0580286368b9e33a142c74536879de7fd3c946543d482 |

示例：wget下载

``` shell
wget https://github.com/lidenger/otpserver/raw/main-doc/doc/download/linux-amd64.zip
```

### 解压
1. 解压（这里以linux/mac系统为例）
```shell
unzip linux-amd64.zip

-rw-r--r-- 1 root root  30M May  8 11:22 app
-rw-r--r-- 1 root root 1.6K May  7 20:33 app.toml
-rw-r--r-- 1 root root  354 Apr 22 14:50 nacos.toml
```
2. app授权可执行：(Windows系统解压即可执行，无需操作)
```shell
chmod +x app
```
3. 文件说明
- app 可执行文件
- app.toml 系统启动配置文件，如果从nacos配置中心启动可删除该文件
- nacos.toml 如果需要从nacos配置中心加载配置需要在该文件中设置nacos server地址

### 配置
这里最简洁的方式配置系统
如果我们选择了MySQL存储


