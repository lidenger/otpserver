package localconf

import (
	"fmt"
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/lidenger/otpserver/pkg/util"
	"os"
	"time"
)

type LocalConf struct {
	testFilePath string
}

var LocalConfIns = &LocalConf{}

func Initialize(conf *config.M) {
	// 初始化local store目录
	dirPath := conf.Server.RootPath + conf.Store.RootPath
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	LocalConfIns.testFilePath = dirPath + "test_file_store.txt"
	isExist, err := util.IsExistsFile(LocalConfIns.testFilePath)
	if err != nil {
		panic(err)
	}
	if !isExist {
		f, err := os.Create(LocalConfIns.testFilePath)
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}
}

func (l *LocalConf) GetStoreType() string {
	return enum.LocalStore
}

func (l *LocalConf) CloseStore() {

}

func (l *LocalConf) TestStore() error {
	content := fmt.Sprintf("Local store检测时间: " + time.Now().Format(time.DateTime) + "\n")
	file, err := os.OpenFile(l.testFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}
