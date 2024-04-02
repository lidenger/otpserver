package serverconf

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"os"
	"reflect"
	"strings"
)

//go:embed app.toml
var appConfig string

// 系统配置
var sysConf *config.M

func InitConfig() *config.M {
	conf := &config.M{}
	_, err := toml.Decode(appConfig, &conf)
	if err != nil {
		panic(fmt.Sprintf("加载配置文件[app.toml]失败:%+v", err))
	}
	reValue := reflect.ValueOf(conf).Elem()
	reType := reflect.TypeOf(conf).Elem()
	recursionDecrypt(reType, reValue, "")
	return conf
}

// 递归解密密文
func recursionDecrypt(typ reflect.Type, val reflect.Value, tomlTag string) {
	// 不是结构体，结束递归
	if val.Kind().String() != "struct" {
		// 解密
		if strings.Contains(tomlTag, "@cipher") {
			key := cmd.P.RootKey256
			iv := cmd.P.IV
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

func GetSysConf() *config.M {
	if sysConf == nil {
		sysConf = LoadConfigByConfSource()
	}
	return sysConf
}

func InitSysConf() {
	GetSysConf()
}

// LoadConfigByConfSource 根据选择的方式加载配置文件
func LoadConfigByConfSource() *config.M {
	var conf *config.M
	confFile := cmd.P.ConfFile
	switch cmd.P.ConfSource {
	case "nacos":
		if len(confFile) == 0 {
			confFile = "nacos.toml"
		}
		conf = LoadConfByNacos(confFile)
		fmt.Println("从nacos加载配置完成")
	case "local":
		if len(confFile) == 0 {
			confFile = "app.toml"
		}
		conf = LoadConfByLocalFile(confFile)
		fmt.Println("从本地文件加载配置完成")
	case "default": // 需要明确指定从默认配置加载而不是自动使用default分支
		conf = InitConfig()
		fmt.Println("从系统默认配置文件加载配置完成")
	default:
		panic(fmt.Sprintf("未知的配置来源:%s", cmd.P.ConfSource))
	}
	return conf
}

// LoadConfByNacos 从Nacos配置中心加载配置
// https://github.com/nacos-group/nacos-sdk-go
func LoadConfByNacos(confFile string) *config.M {
	content, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Printf("Nacos配置文件[%s]不正确(读取文件失败):%+v", confFile, err)
		panic(err)
	}
	conf := &config.NacosM{}
	_, err = toml.Decode(string(content), &conf)
	// 配置client
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(conf.Client.NamespaceId),
		constant.WithTimeoutMs(conf.Client.TimeoutMs),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir(conf.Client.LogDir),
		constant.WithCacheDir(conf.Client.CacheDir),
		constant.WithLogLevel(conf.Client.LogLevel),
	)
	// 配置server
	scSlice := make([]constant.ServerConfig, 0)
	for _, server := range conf.ServerArr {
		sone := *constant.NewServerConfig(server.Ip, server.Port, constant.WithContextPath(server.ContextPath))
		scSlice = append(scSlice, sone)
	}
	sc := scSlice[:]
	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}

	c, err := client.GetConfig(vo.ConfigParam{
		DataId: conf.Client.DataId,
		Group:  conf.Client.Group,
	})

	m := &config.M{}
	_, err = toml.Decode(c, &m)

	err = client.ListenConfig(vo.ConfigParam{
		DataId: conf.Client.DataId,
		Group:  conf.Client.Group,
		OnChange: func(namespace, group, dataId, data string) {
			if group != conf.Client.Group || dataId != conf.Client.Group {
				return
			}
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", content:" + data)
			_, err = toml.Decode(data, &m)
			fmt.Printf("nacos config refresh %+v", c)
		},
	})
	return m
}

// LoadConfByLocalFile 从本地文件加载配置
func LoadConfByLocalFile(confFile string) *config.M {
	content, err := os.ReadFile(confFile)
	if err != nil {
		fmt.Printf("Nacos配置文件[nacos.toml]不正确(读取文件失败):%+v", err)
		panic(err)
	}
	m := &config.M{}
	_, err = toml.Decode(string(content), &m)
	return m
}
