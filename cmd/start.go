package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/util"
	"os"
)

// CodeLevelProtectKey 提供代码级保护
// 注意：即使对[app.key]文件提供了保护并不意味[app.key]文件可以被泄露
// 由于代码是开源的所以是可以获取到的，因此[app.key]文件依然是高敏感的
var CodeLevelProtectKey = []byte("0d3d2206ea5711ee87ba2cf05daf3fe5")
var CodeLevelProtectIV = []byte("ad5457f8ea5711ee")

// Param 命令参数
type Param struct {
	IsInitMode  bool   // 初始化模式
	IsToolMode  bool   // 工具模式
	Encrypt     bool   // 工具模式-加密数据
	EncryptData string // 工具模式-加密的源数据
	Env         string // dev,prod
	Port        int    // 服务启动端口
	MainStore   string // 主存储 mysql,pgsql,oracle
	BackupStore string // 备存储 mysql,pgsql,oracle
	AppKeyFile  string // app key file 路径
	*Crypt
}

type Crypt struct {
	RootKey   string // 根密钥
	IV        string // CBC IV
	DataCheck string // 数据校验
}

var P *Param

func InitParam() {
	P = &Param{}
	flag.BoolVar(&P.IsInitMode, "init", false, "系统初始化")
	flag.BoolVar(&P.IsToolMode, "tool", false, "工具模式")
	flag.BoolVar(&P.Encrypt, "encrypt", false, "工具模式-加密")
	flag.StringVar(&P.EncryptData, "data", "", "工具模式-加密数据")
	flag.StringVar(&P.Env, "env", "dev", "当前环境[dev,prod]")
	flag.IntVar(&P.Port, "port", 8080, "服务启动端口")
	flag.StringVar(&P.MainStore, "mainStore", "mysql", "主存储[mysql,pgsql,oracle]")
	flag.StringVar(&P.BackupStore, "backupStore", "", "备存储[mysql,pgsql,oracle]")
	flag.StringVar(&P.AppKeyFile, "keyFile", "app.key", "系统启动KEY文件")
	flag.Parse()
}

// InitMode 系统初始化模式，生成系统启动文件（高敏感文件，需要较强的管理流程）
func InitMode() {
	rootKey := util.Generate32Str()
	iv := util.Generate16Str()
	digest := crypt.HmacDigest(CodeLevelProtectKey, rootKey+iv)
	crypto := &Crypt{
		RootKey:   rootKey,
		IV:        iv,
		DataCheck: digest,
	}
	content, err := json.Marshal(crypto)
	if err != nil {
		panic(err)
	}
	cipher, err := crypt.Encrypt(CodeLevelProtectKey, CodeLevelProtectIV, content)
	if err != nil {
		panic(err)
	}
	keyFile := "app.key"
	_, err = os.Stat(keyFile)
	if err == nil {
		panic("请注意：系统启动文件[app.key]已存在，" +
			"请确认是否需要重新生成，如果删除当前的[app.key]，历史的数据将无法正常使用！！" +
			"如果确认生成新的[app.key]，请删除当前的")
	}
	err = os.WriteFile(keyFile, []byte(cipher), 0x600)
	if err != nil {
		panic(err)
	}
	fmt.Println("系统启动文件[app.key]在当前目录已生成，该文件为敏感文件请妥善保管")
}

func AnalysisKeyFile(keyFile string) *Crypt {
	cipher, err := os.ReadFile(keyFile)
	if err != nil {
		fmt.Printf("系统启动[app.key]文件不正确(读取文件失败):%+v", err)
		panic(err)
	}
	content, err := crypt.Decrypt(CodeLevelProtectKey, CodeLevelProtectIV, string(cipher))
	if err != nil {
		fmt.Printf("系统启动[app.key]文件不正确(解密失败):%+v", err)
		panic(err)
	}
	crypto := &Crypt{}
	err = json.Unmarshal(content, crypto)
	if err != nil {
		fmt.Printf("系统启动[app.key]文件不正确(反序列化失败):%+v", err)
		panic(err)
	}
	// 对比数据摘要
	digest := crypt.HmacDigest(CodeLevelProtectKey, crypto.RootKey+crypto.IV)
	if digest != crypto.DataCheck {
		panic("系统启动[app.key]文件不正确(数据校验不通过)")
	}
	return crypto
}

// ToolMode 工具模式
func ToolMode() {
	if P.Encrypt {
		if len(P.EncryptData) == 0 {
			panic("加密模式没有提供加密数据,请使用[-data=\"xxx\"]提供加密数据")
		}
		key := []byte(P.RootKey)
		iv := []byte(P.IV)
		data := []byte(P.EncryptData)
		cipher, err := crypt.Encrypt(key, iv, data)
		if err != nil {
			panic(err)
		}
		fmt.Printf("数据:%s,加密后密文:%s", P.EncryptData, cipher)
	}
}
