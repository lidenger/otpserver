package model

// AccountSecretModel 账号密钥
type AccountSecretModel struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	SecretSeed string `gorm:"column:secret_seed_cipher" json:"secret"` // 密钥种子密文 = AES(KEY, 密钥种子)
	Account    string `gorm:"column:account" json:"account"`           // 账号
	DataCheck  string `gorm:"column:data_check" json:"dataCheck"`      // 数据校验 = HMACSHA256(KEY, secret_seed_cipher + account + is_enable)
	Common
}

func (a *AccountSecretModel) TableName() string {
	return "otp_account_secret"
}
