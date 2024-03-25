package service

import (
	"encoding/hex"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/config/store/mysqlconf"
	"github.com/lidenger/otpserver/config/store/pgsqlconf"
	"github.com/lidenger/otpserver/internal/store/mysqlstore"
	"github.com/lidenger/otpserver/internal/store/pgsqlstore"
)

var SecretSvcIns = &SecretSvc{}
var ServerSvcIns = &ServerSvc{}

func InitSvc() {
	rootKey, err := hex.DecodeString(cmd.P.RootKey)
	if err != nil {
		panic("解析rootKey失败")
	}
	iv, err := hex.DecodeString(cmd.P.IV)
	if err != nil {
		panic("解析IV失败")
	}
	SecretSvcIns.RootKey = rootKey
	SecretSvcIns.IV = iv

	ServerSvcIns.RootKey = rootKey
	ServerSvcIns.IV = iv

	switch cmd.P.MainStore {
	case "mysql":
		SecretSvcIns.Store = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.Store = &mysqlstore.ServerStore{DB: mysqlconf.DB}
	case "pgsql":
		SecretSvcIns.Store = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	switch cmd.P.BackupStore {
	case "mysql":
		SecretSvcIns.StoreBackup = &mysqlstore.SecretStore{DB: mysqlconf.DB}
		ServerSvcIns.StoreBackup = &mysqlstore.ServerStore{DB: mysqlconf.DB}
	case "pgsql":
		SecretSvcIns.StoreBackup = &pgsqlstore.SecretStore{DB: pgsqlconf.DB}
	}
	log.Info("Service初始化完成:%s", "SecretSvc", "ServerSvc")
}
