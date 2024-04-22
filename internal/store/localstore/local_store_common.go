package localstore

import (
	"context"
	"encoding/json"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/util"
	"os"
	"path/filepath"
	"time"
)

// local store 文件处理
func manageLocalStoreFile(rootPath, fileName string) error {
	fullPath := filepath.Join(rootPath, fileName)
	isExists, err := util.IsExistsFile(fullPath)
	if err != nil {
		return err
	}
	if isExists {
		// 重命名 -> bak
		timestamp := time.Now().Format("20060102150405")
		bakPath := filepath.Join(rootPath, "/bak", fileName+"_bak_"+timestamp)
		err = os.Rename(fullPath, bakPath)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 获取所有store数据写入文件
func fetchAllStoreDataAndWriteToFile[R any](ctx context.Context, s store.SelectAllFunc[R], rootPath, fileName string) error {
	err := manageLocalStoreFile(rootPath, fileName)
	if err != nil {
		return err
	}
	fullPath := filepath.Join(rootPath, fileName)
	file, err := os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY, os.ModePerm)
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
