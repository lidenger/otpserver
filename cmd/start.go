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
	Init         bool   // 初始化
	Env          string // dev,prod
	Port         int    // 服务启动端口
	MainStore    string // 主存储 mysql,pgsql,oracle
	BackupStore  string // 备存储 mysql,pgsql,oracle
	StartKeyFile string // start key file 路径
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
	flag.BoolVar(&P.Init, "init", false, "系统初始化")
	flag.StringVar(&P.Env, "env", "dev", "当前环境[dev,prod]")
	flag.IntVar(&P.Port, "port", 8080, "服务启动端口")
	flag.StringVar(&P.MainStore, "mainStore", "mysql", "主存储[mysql,pgsql,oracle]")
	flag.StringVar(&P.BackupStore, "backupStore", "", "备存储[mysql,pgsql,oracle]")
	flag.StringVar(&P.StartKeyFile, "keyFile", "app.key", "系统启动KEY文件")
	flag.Parse()
}

// GenKeyFile 生成系统启动文件（高敏感文件，需要较强的管理流程）
func GenKeyFile() {
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
		panic(fmt.Sprintf("系统启动[app.key]文件已存在:%s", keyFile))
	}
	err = os.WriteFile(keyFile, []byte(cipher), 0x600)
	if err != nil {
		panic(err)
	}
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
