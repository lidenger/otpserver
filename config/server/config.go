package server

import (
	_ "embed"
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
		RootPath string `toml:"rootPath"`
	}
	Log struct {
		Level      string `toml:"level"`
		RootPath   string `toml:"rootPath"`
		AppFile    string `toml:"appFile"`
		MaxSize    int    `toml:"maxSize"`
		MaxBackups int    `toml:"maxBackups"`
		MaxAge     int    `toml:"maxAge"`
		Compress   bool   `toml:"compress"`
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
