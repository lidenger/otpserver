package serverconf

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lidenger/otpserver/cmd"
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
		ReqLimit int    `toml:"reqLimit"`
	} `toml:"server"`
	Log struct {
		Level      string `toml:"level"`
		RootPath   string `toml:"rootPath"`
		AppFile    string `toml:"appFile"`
		HttpFile   string `toml:"httpFile"`
		MaxSize    int    `toml:"maxSize"`
		MaxBackups int    `toml:"maxBackups"`
		MaxAge     int    `toml:"maxAge"`
		Compress   bool   `toml:"compress"`
	} `toml:"log"`
	Store struct {
		MainStore   string `toml:"mainStore"`
		BackupStore string `toml:"backupStore"`
	} `toml:"store"`
	MySQL struct {
		Address         string `toml:"address"`
		UserName        string `toml:"userName"`
		Password        string `toml:"password@cipher"`
		DbName          string `toml:"dbName"`
		ConnMaxLifeTime int    `toml:"connMaxLifeTime"`
		MaxIdleConn     int    `toml:"maxIdleConn"`
		MaxOpenConn     int    `toml:"maxOpenConn"`
		ConnMaxWaitTime string `toml:"connMaxWaitTime"`
	} `toml:"mysql"`
}

func InitConfig() *Config {
	var conf = devConf
	if cmd.P.Env == "prod" {
		conf = prodConf
	}
	config := &Config{}
	_, err := toml.Decode(conf, &config)
	if err != nil {
		panic(fmt.Sprintf("加载%s配置文件失败:%+v", cmd.P.Env, err))
	}
	return config
}
