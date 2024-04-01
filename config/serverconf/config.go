package serverconf

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lidenger/otpserver/cmd"
	conf "github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/pkg/crypt"
	"reflect"
	"strings"
)

//go:embed ../localdev/app.toml
var appConfig string

func InitConfig() *conf.M {
	config := &conf.M{}
	_, err := toml.Decode(appConfig, &config)
	if err != nil {
		panic(fmt.Sprintf("加载配置文件[app.toml]失败:%+v", err))
	}
	reValue := reflect.ValueOf(config).Elem()
	reType := reflect.TypeOf(config).Elem()
	recursionDecrypt(reType, reValue, "")
	return config
}

// 递归解密密文
func recursionDecrypt(typ reflect.Type, val reflect.Value, tomlTag string) {
	// 不是结构体，结束递归
	if val.Kind().String() != "struct" {
		// 解密
		if strings.Contains(tomlTag, "@cipher") {
			key := []byte(cmd.P.RootKey256)
			iv := []byte(cmd.P.IV)
			data, err := crypt.Decrypt(key, iv, val.String())
			if err != nil {
				fmt.Printf("密文解密失败,key:%s,value:%s\n", tomlTag, val.String())
				panic(err)
			}
			val.SetString(string(data))
		}
		return
	}
	for i := 0; i < val.NumField(); i++ {
		// 只解析toml tag
		tag := typ.Field(i).Tag.Get("toml")
		if tag == "" {
			continue
		}
		fieldVal := val.Field(i)
		recursionDecrypt(typ.Field(i).Type, fieldVal, tag)
	}
}
