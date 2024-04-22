package localconf

import (
	"fmt"
	"github.com/lidenger/otpserver/config"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/enum"
	"github.com/lidenger/otpserver/pkg/util"
	"os"
	"path"
	"path/filepath"
	"time"
)

type LocalConf struct {
	testFilePath string
}

var Ins = &LocalConf{}

func Initialize(conf *config.M) {
	// 初始化local store目录
	dirPath := path.Join(conf.Server.RootPath, conf.Store.RootPath)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	// 初始化 local store bak 目录
	bakDirPath := path.Join(conf.Server.RootPath, conf.Store.RootPath, "/bak")
	err = os.MkdirAll(bakDirPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	Ins.testFilePath = filepath.Join(dirPath, "test_file_store")
	isNotExist, err := util.IsNotExistsFile(Ins.testFilePath)
	if err != nil {
		panic(err)
	}
	if isNotExist {
		f, err := os.Create(Ins.testFilePath)
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
	log.Info("本地存储已关闭")
}

// TestStore 检测文件存储是否正常，通过写入当前时间来做判断
func (l *LocalConf) TestStore() error {
	content := fmt.Sprintf("Local store检测时间: " + time.Now().Format(time.DateTime) + "\n")
	isExist, err := util.IsExistsFile(l.testFilePath)
	if err != nil {
		return err
	}
	if isExist {
		fi, err := os.Stat(l.testFilePath)
		if err != nil {
			return err
		}
		// 超过10MB删除
		if fi.Size() > 10*1024*1024 {
			err = os.Remove(l.testFilePath)
			if err != nil {
				return err
			}
		}
	}
	file, err := os.OpenFile(l.testFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}
