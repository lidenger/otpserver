package service

import (
	"context"
	"fmt"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type ServerIpListSvc struct {
	Store                   store.ServerIpListStore // 主存储
	StoreBackup             store.ServerIpListStore // 备存储
	StoreMemory             store.ServerIpListStore // 内存存储
	storeDetectionEventChan chan<- struct{}
}

func (s *ServerIpListSvc) GetBySign(ctx context.Context, serverSign string) ([]*model.ServerIpListModel, error) {
	var err error
	p := &param.ServerIpListParam{Sign: serverSign}
	data, err := MultiStoreSelectByCondition[*param.ServerIpListParam, *model.ServerIpListModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (s *ServerIpListSvc) Add(ctx context.Context, p *param.ServerIpListParam) error {
	ips, err := MultiStoreSelectByCondition[*param.ServerIpListParam, *model.ServerIpListModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return err
	}
	if len(ips) != 0 {
		msg := fmt.Sprintf("服务:%s,IP:%s已存在不能重复添加", p.Sign, p.IP)
		return otperr.ErrRepeatAdd(msg)
	}
	m := &model.ServerIpListModel{
		ServerSign: p.Sign,
		IP:         p.IP,
	}
	err = MultiStoreInsert[*model.ServerIpListModel](ctx, s.storeDetectionEventChan, m, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}

func (s *ServerIpListSvc) Remove(ctx context.Context, p *param.ServerIpListParam) error {
	ips, err := MultiStoreSelectByCondition[*param.ServerIpListParam, *model.ServerIpListModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		msg := fmt.Sprintf("服务:%s,IP:%s不存在无法删除", p.Sign, p.IP)
		return otperr.ErrParamIllegal(msg)
	}
	m := util.GetArrFirstItem(ips)
	err = MultiStoreDelete(ctx, s.storeDetectionEventChan, m.ID, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}
