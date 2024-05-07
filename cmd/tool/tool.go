package tool

import (
	"context"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/service"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/enum"
)

// EncryptMode 工具模式
func EncryptMode() {
	if len(cmd.P.EncryptData) == 0 {
		panic("加密模式没有提供加密数据,请使用[-data=\"xxx\"]提供加密数据")
	}
	data := []byte(cmd.P.EncryptData)
	cipher, err := crypt.Encrypt(cmd.P.RootKey256, cmd.P.IV, data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("数据:%s,加密后密文:%s", cmd.P.EncryptData, cipher)
}

// InitAdminMode 初始化admin账号密码(也可用于重置admin账号密码)
func InitAdminMode() {
	if len(cmd.P.AdminPassword) == 0 {
		panic("admin密码不能为空")
	}
	ctx := context.Background()
	conf, err := service.ConfSvcIns.GetByKey(ctx, enum.AdminPasswordKey)
	if err != nil {
		panic(err)
	}
	pwdDigest := crypt.HmacDigest(cmd.P.RootKey256, cmd.P.AdminPassword)
	confParam := &param.SysConfParam{
		Key:    enum.AdminPasswordKey,
		Val:    pwdDigest,
		Remark: "Admin账号密码",
	}
	if conf == nil {
		err = service.ConfSvcIns.Add(ctx, confParam)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Admin账号密码初始化成功,账号:admin,密码:%s", cmd.P.AdminPassword)
	} else {
		err = service.ConfSvcIns.Update(ctx, confParam)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Admin账号密码更新成功,账号:admin,密码:%s", cmd.P.AdminPassword)
	}

}
