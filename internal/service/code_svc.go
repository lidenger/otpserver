package service

import (
	"context"
	"github.com/lidenger/otp"
	"github.com/lidenger/otpserver/config/log"
	"github.com/lidenger/otpserver/pkg/otperr"
	"github.com/patrickmn/go-cache"
	"time"
)

// 定义code缓存，用于鉴别重复验证
// 每个窗口为50秒所以定义3分钟有效期足够了
var codeCache = cache.New(3*time.Minute, 5*time.Minute)

func ValidCode(ctx context.Context, account, code string) (bool, error) {
	model, err := SecretSvcIns.GetByAccount(ctx, account)
	if err != nil {
		return false, err
	}
	if model == nil {
		return false, otperr.ErrDataNotExists("账号[%s]不存在", account)
	}
	code2, err := otp.TOTP(model.SecretSeed)
	if err != nil {
		return false, otperr.ErrGenCode(err)
	}
	if code != code2 {
		return false, nil
	}
	// 查看code是否使用过（这里需要符合otp定义：一次性密钥）
	_, exists := codeCache.Get(code)
	if exists {
		log.Warn("账号OTP重复验证:%s", account, code)
		return false, nil
	}
	// 标记code已使用
	_ = codeCache.Add(code, struct{}{}, cache.DefaultExpiration)
	return true, nil
}
