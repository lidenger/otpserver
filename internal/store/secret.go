package store

import (
	"context"
	"github.com/lidenger/otpserver/internal/model"
)

type ISecret interface {
	Insert(ctx context.Context, m *model.AccountSecretModel)
	Update(ctx context.Context, m *model.AccountSecretModel)
	List(ctx context.Context, account string, pageNo, pageSize int)
}
