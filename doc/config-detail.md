## app配置文件详细描述
### 说明
配置文件为app.toml内容，可以选择在本地加载app.toml文件也可以将内容配置到nacos中，通过配置中心启动服务

### [server] 系统配置
```shell
[server]
# 服务启动端口
port = 8066

# 系统工作主目录，里面包含了日志和local store数据
rootPath = "./otpwork"

# 限速QPS
reqLimit = 100

# access token有效期(小时)，外部接入到该系统的认证凭证
accessTokenValidHour = 4

# time token有效期(分钟)，外部接入到该系统请求获取认证凭证时最大允许的时间差
# 例如这里设置了2分钟代表otp服务接收到外部服务的accessToken请求最大允许差2分钟
# 也有可能存在client和server端时钟差距超过2分钟导致获取accessToken失败的情况
# 合理设置该值可以在保证安全的情况下外部服务正常接入
timeTokenValidMinute = 2

# 服务域名(正式环境应该替换为相应的域名)，这里配置的域名被用于设置cookie时所属的域名属性，是一个重要的安全属性
# 异常配置会导致无法登陆admin管理后台，但不会影响外部系统的调用
domain = "localhost"

# admin登陆态有效期(小时)，admin管理后台应该在不使用时主动退出登陆
# 这个配置可以在忘记登出时的12小时后让登陆态失效，清除浏览器cookie和登录态验证都会使用到这个参数
adminLoginValidHour = 12
```

### [log] 日志配置
```shell
[log]
# 日志输出级别，debug info error
level = "info"

# 日志文件根目录，在生成目录时会主动拼接[server]配置的[rootPath]路径
rootPath = "/log/"

# 系统日志文件名
appFile = "app.log"

# http请求日志文件名
httpFile = "request.log"

# 每个日志文件最大大小(MB)，超过大小会自动切分日志文件
maxSize = 10

# 旧日志保留的个数，历史日志个数超过该值时删除旧的
maxBackups = 100

# 旧日志保留的最大天数，历史日志文件生成日期超过该值时删除旧的
maxAge = 30

# 是否压缩旧日志文件
compress = true

```

### [store] 存储配置
```shell
[store]
# 主存储,当前可选: mysql,pgsql
mainStore = "mysql"

# 备存储,如果选择了备存储可提高系统容灾能力和可用性，当主存储不可用时备存储可提供读的能力,可不选择,当前可选: "",mysql,pgsql
backupStore = "" 

# 存储的目录，在生成目录时会主动拼接[server]配置的[rootPath]路径
rootPath = "/store/"

# 是否启用本地存储（数据冗余一份到本地磁盘） true | false
isEnableLocal = true

# 周期备份数据到本地，注意只有启用本地存储才有效，单位：小时
cycleBakToLocalHour = 4

# 是否启用内存存储（数据冗余一份到内存，优先查询） true | false
isEnableMemory = true
```

### [mysql] MySQL配置
```shell
[mysql]
# 数据库地址
address = "127.0.0.1:3306"

# 账号
userName = "root"

# 密码，使用[app -encrypt -data "xxx"]加密
# 如果不启用MySQL,可不配置密码，例如："password@cipher" = ""
"password@cipher" = "ac1b926895e6982a7060f1ca10248659"

# 数据库名称
dbName = "otp"

# 连接可复用的最大时间,分钟
connMaxLifeTime = 30

# 空闲连接池中连接的最大数量
maxIdleConn = 30

# 打开数据库连接的最大数量
maxOpenConn = 200

# 连接超时时间,秒
connMaxWaitTime = "20s"
```

### [pgsql] PostgreSQL配置
```shell
[pgsql]
host = "127.0.0.1"
port = "5432"
userName = "postgres"
# 密码，使用[app -encrypt -data "xxx"]加密
# 如果不启用PostgerSQL,可不配置密码，例如："password@cipher" = ""
"password@cipher" = ""
dbName = "otp"
```