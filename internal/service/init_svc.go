package service

import (
	"encoding/hex"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/serverconf"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/store/mysqlstore"
	"github.com/lidenger/otpserver/internal/store/pgsqlstore"
)

var AccountSecretSvc = &SecretSvc{}

func InitSvc() {
	rootKey, err := hex.DecodeString(serverconf.CMD.RootKey)
	if err != nil {
		panic("解析rootKey失败")
	}
	iv, err := hex.DecodeString(serverconf.CMD.IV)
	if err != nil {
		panic("解析IV失败")
	}
	AccountSecretSvc.RootKey = rootKey
	AccountSecretSvc.IV = iv
	switch serverconf.CMD.MainStore {
	case "mysql":
		AccountSecretSvc.Store = &mysqlstore.SecretStore{DB: mysqlconf.DB}
	case "pgsql":
		AccountSecretSvc.Store = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	switch serverconf.CMD.BackupStore {
	case "mysql":
		AccountSecretSvc.StoreBackup = &mysqlstore.SecretStore{DB: mysqlconf.DB}
	case "pgsql":
		AccountSecretSvc.StoreBackup = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}
