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

var Conf = new(Config)

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
		Address  string `toml:"address"`
		UserName string `toml:"userName"`
	}
}

func InitConfig(env string) {
	var conf = devConf
	if env == "prod" {
		conf = prodConf
	}
	_, err := toml.Decode(conf, &Conf)
	if err != nil {
		panic(fmt.Sprintf("加载%s配置文件失败:%+v", env, err))
	}
}
