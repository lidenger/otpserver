package localstore

import (
	"context"
	"encoding/json"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/util"
	"os"
	"time"
)

// local store 文件处理
func manageLocalStoreFile(filepath string) error {
	isExists, err := util.IsExistsFile(filepath)
	if err != nil {
		return err
	}
	if isExists {
		// 重命名 -> bak
		timestamp := time.Now().Format("20060102150405")
		err = os.Rename(filepath, filepath+"_bak_"+timestamp)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 获取所有store数据写入文件
func fetchAllStoreDataAndWriteToFile[R any](ctx context.Context, s store.SelectAllFunc[R], filePath string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := s.SelectAll(ctx)
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = file.Write(js)
	if err != nil {
		return err
	}
	return nil
}
