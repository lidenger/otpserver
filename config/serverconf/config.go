package serverconf

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
)

//go:embed dev.toml
var devConf string

//go:embed prod.toml
var prodConf string

type Config struct {
	Server struct {
		Env      string `toml:"env"`
		Port     int    `toml:"port"`
		RootPath string `toml:"rootPath"`
	}
	Log struct {
		Level      string `toml:"level"`
		RootPath   string `toml:"rootPath"`
		AppFile    string `toml:"appFile"`
		HttpFile   string `toml:"httpFile"`
		MaxSize    int    `toml:"maxSize"`
		MaxBackups int    `toml:"maxBackups"`
		MaxAge     int    `toml:"maxAge"`
		Compress   bool   `toml:"compress"`
	}
	Store struct {
		MainStore   string `toml:"mainStore"`
		BackupStore string `toml:"backupStore"`
	}
	MySQL struct {
		Address         string `toml:"address"`
		UserName        string `toml:"userName"`
		Password        string `toml:"password@cipher"`
		DbName          string `toml:"dbName"`
		ConnMaxLifeTime int    `toml:"connMaxLifeTime"`
		MaxIdleConn     int    `toml:"maxIdleConn"`
		MaxOpenConn     int    `toml:"maxOpenConn"`
		ConnMaxWaitTime string `toml:"connMaxWaitTime"`
	}
}

func InitConfig(env string) *Config {
	var conf = devConf
	if env == "prod" {
		conf = prodConf
	}
	config := &Config{}
	_, err := toml.Decode(conf, &config)
	if err != nil {
		panic(fmt.Sprintf("加载%s配置文件失败:%+v", env, err))
	}
	return config
}

// CmdParam 命令参数
type CmdParam struct {
	Env         string // dev,prod
	Port        int    // 服务启动端口
	MainStore   string // 主存储 mysql,pgsql,oracle
	BackupStore string // 备存储 mysql,pgsql,oracle
}

var cmdParam *CmdParam

func GetCmdParam() *CmdParam {
	if cmdParam != nil {
		return cmdParam
	}
	cmdParam = &CmdParam{}
	flag.StringVar(&cmdParam.Env, "env", "dev", "当前环境[dev,prod]")
	flag.IntVar(&cmdParam.Port, "port", 8080, "服务启动端口")
	flag.StringVar(&cmdParam.MainStore, "mainStore", "mysql", "主存储[mysql,pgsql,oracle]")
	flag.StringVar(&cmdParam.BackupStore, "backupStore", "", "备存储[mysql,pgsql,oracle]")
	flag.Parse()
	return cmdParam
}
