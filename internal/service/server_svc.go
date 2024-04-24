package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/lidenger/otpserver/cmd"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/crypt"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type ServerSvc struct {
	Store                   store.ServerStore // 主存储
	StoreBackup             store.ServerStore // 备存储
	StoreMemory             store.ServerStore // 内存存储
	storeDetectionEventChan chan<- struct{}
}

// Add 增加接入的服务
func (s *ServerSvc) Add(ctx context.Context, p *param.ServerParam) error {
	exists, err := s.IsExists(ctx, p.Sign)
	if err != nil {
		return err
	}
	if exists {
		msg := fmt.Sprintf("服务%s已存在不能重复添加", p.Sign)
		return otperr.ErrRepeatAdd(errors.New(msg))
	}
	m, err := s.NewServerModel(p)
	if err != nil {
		return err
	}
	err = MultiStoreInsert[*model.ServerModel](ctx, s.storeDetectionEventChan, m, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}

func (s *ServerSvc) NewServerModel(p *param.ServerParam) (*model.ServerModel, error) {
	m := &model.ServerModel{}
	m.Name = p.Name
	m.Sign = p.Sign
	m.Remark = p.Remark
	m.IsEnable = p.IsEnable
	m.IsOperateSensitiveData = p.IsOperateSensitiveData
	m.IsEnableIPlist = p.IsEnableIPlist
	// 默认启用
	if m.IsEnable == 0 {
		m.IsEnable = 1
	}
	// 默认不启用操作敏感信息
	if m.IsOperateSensitiveData == 0 {
		m.IsOperateSensitiveData = 2
	}
	// 默认不启用IP白名单
	if m.IsEnableIPlist == 0 {
		m.IsEnableIPlist = 2
	}
	// 服务密钥
	secret := util.Generate32Str()
	var err error
	secretCipher, err := crypt.Encrypt(cmd.P.RootKey192, cmd.P.IV, []byte(secret))
	if err != nil {
		return nil, err
	}
	m.Secret = secretCipher
	// 服务密钥IV
	iv := util.Generate16Str()
	ivCipher, err := crypt.Encrypt(cmd.P.RootKey192, cmd.P.IV, []byte(iv))
	if err != nil {
		return nil, err
	}
	m.IV = ivCipher
	// 计算数据摘要
	m.DataCheck = s.CalcDataCheckSum(m)
	return m, nil
}

// CalcDataCheckSum 计算数据校验和
func (s *ServerSvc) CalcDataCheckSum(m *model.ServerModel) string {
	data := fmt.Sprintf("%d,%d,%d,%s,%s,%s", m.IsEnable, m.IsOperateSensitiveData, m.IsEnableIPlist, m.Sign, m.Secret, m.IV)
	return crypt.HmacDigest(cmd.P.RootKey192, data)
}

// IsExists 账号密钥是否存在
func (s *ServerSvc) IsExists(ctx context.Context, sign string) (bool, error) {
	secretData, err := s.GetBySign(ctx, sign, false)
	if err != nil {
		return false, err
	}
	return secretData != nil, nil
}

func (s *ServerSvc) GetBySign(ctx context.Context, sign string, isDecrypt bool) (*model.ServerModel, error) {
	var err error
	p := &param.ServerParam{Sign: sign}
	data, err := MultiStoreSelectByCondition[*param.ServerParam, *model.ServerModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return nil, err
	}
	serverModel := util.GetArrFirstItem(data)
	if serverModel == nil {
		return nil, nil
	}
	err = s.CheckModel(serverModel, isDecrypt)
	if err != nil {
		return nil, err
	}
	return serverModel, err
}

// CheckModel 校验数据
// isDecrypt 解密服务密钥密文和IV密文
func (s *ServerSvc) CheckModel(m *model.ServerModel, isDecrypt bool) error {
	check := s.CalcDataCheckSum(m)
	if m.DataCheck != check {
		msg := fmt.Sprintf("服务数据校验不通过,疑似被篡改,请关注(ID:%d,sign:%s)", m.ID, m.Sign)
		return otperr.ErrAccountSecretDataCheck(errors.New(msg))
	}
	if isDecrypt {
		// 服务密钥
		secret, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.Secret)
		if err != nil {
			return otperr.ErrDecrypt(err)
		}
		m.Secret = string(secret)
		// 服务密钥IV
		iv, err := crypt.Decrypt(cmd.P.RootKey192, cmd.P.IV, m.IV)
		if err != nil {
			return otperr.ErrDecrypt(err)
		}
		m.IV = string(iv)
	}
	return nil
}

// Paging 分页
func (s *ServerSvc) Paging(ctx context.Context, p *param.ServerPagingParam) (result []*model.ServerModel, count int64, err error) {
	result, count, err = MultiStorePaging[*param.ServerPagingParam, *model.ServerModel](ctx, s.storeDetectionEventChan, p, s.Store, s.StoreBackup)
	return
}

// SetEnable 更新启用状态和操作敏感信息状态
func (s *ServerSvc) SetEnable(ctx context.Context, p *param.ServerParam) error {
	m, err := s.GetBySign(ctx, p.Sign, false)
	if err != nil {
		return err
	}
	if m == nil {
		return otperr.ErrParamIllegal("服务不存在:" + p.Sign)
	}
	// 进一步验证参数正确性
	if m.ID != p.ID {
		return otperr.ErrParamIllegal(fmt.Sprintf("非法参数:%d", p.ID))
	}
	// 数据一致无需更新
	if m.IsEnable == p.IsEnable &&
		m.IsOperateSensitiveData == p.IsOperateSensitiveData &&
		m.IsEnableIPlist == p.IsEnableIPlist {
		return nil
	}
	m.IsEnable = p.IsEnable
	m.IsOperateSensitiveData = p.IsOperateSensitiveData
	m.IsEnableIPlist = p.IsEnableIPlist
	checkSum := s.CalcDataCheckSum(m)
	params := make(map[string]any)
	params["is_enable"] = m.IsEnable
	params["is_operate_sensitive_data"] = m.IsOperateSensitiveData
	params["is_enable_iplist"] = m.IsEnableIPlist
	params["data_check"] = checkSum
	err = MultiStoreUpdate(ctx, s.storeDetectionEventChan, m.ID, params, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}

func (s *ServerSvc) EditBase(ctx context.Context, p *param.ServerParam) error {
	m, err := s.GetBySign(ctx, p.Sign, false)
	if err != nil {
		return err
	}
	if m == nil {
		return otperr.ErrParamIllegal("服务不存在:" + p.Sign)
	}
	// 进一步验证参数正确性
	if m.ID != p.ID {
		return otperr.ErrParamIllegal(fmt.Sprintf("非法参数:%d", p.ID))
	}
	// 数据一致无需更新
	if m.Name == p.Name && m.Remark == p.Remark {
		return nil
	}
	params := make(map[string]any)
	if len(p.Name) != 0 {
		params["server_name"] = p.Name
	}
	if len(p.Remark) != 0 {
		params["server_remark"] = p.Remark
	}
	err = MultiStoreUpdate(ctx, s.storeDetectionEventChan, m.ID, params, s.Store, s.StoreBackup)
	return err
}
