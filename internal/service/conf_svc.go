package service

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/lidenger/otpserver/internal/model"
	"github.com/lidenger/otpserver/internal/param"
	"github.com/lidenger/otpserver/internal/store"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/lidenger/otpserver/pkg/util"
)

type ConfSvc struct {
	Store                   store.ConfStore // 主存储
	StoreBackup             store.ConfStore // 备存储
	StoreMemory             store.ConfStore // 内存存储
	storeDetectionEventChan chan<- struct{}
}

func (s *ConfSvc) GetByKey(ctx context.Context, key string) (*model.SysConfModel, error) {
	var err error
	p := &param.SysConfParam{Key: key}
	data, err := MultiStoreSelectByCondition[*param.SysConfParam, *model.SysConfModel](ctx, s.storeDetectionEventChan, p, s.StoreMemory, s.Store, s.StoreBackup)
	if err != nil {
		return nil, err
	}
	conf := util.GetArrFirstItem(data)
	return conf, err
}

func (s *ConfSvc) Add(ctx context.Context, p *param.SysConfParam) error {
	conf, err := s.GetByKey(ctx, p.Key)
	if err != nil {
		return err
	}
	if conf != nil {
		return otperr.ErrRepeatAdd("系统配置%s已存在不能重复添加", p.Key)
	}
	m := &model.SysConfModel{}
	err = copier.Copy(m, p)
	if err != nil {
		return err
	}
	err = MultiStoreInsert[*model.SysConfModel](ctx, s.storeDetectionEventChan, m, s.Store, s.StoreBackup, s.StoreMemory)
	return err
}

func (s *ConfSvc) Update(ctx context.Context, p *param.SysConfParam) error {
	conf, err := s.GetByKey(ctx, p.Key)
	if err != nil {
		return err
	}
	if conf == nil {
		return otperr.ErrDataNotExists("系统配置%s不存在", p.Key)
	}
	params := make(map[string]any)
	if len(p.Val) != 0 {
		params["sys_val"] = p.Val
	}
	if len(p.Remark) != 0 {
		params["remark"] = p.Remark
	}
	err = MultiStoreUpdate(ctx, s.storeDetectionEventChan, conf.ID, params, s.Store, s.StoreBackup)
	return err
}
