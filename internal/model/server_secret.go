package model

// ServerSecretModel 接入服务密钥，用于认证
type ServerSecretModel struct {
	ID           int64  `gorm:"column:id;primary_key" json:"id"`
	ServerSign   string `gorm:"column:server_sign" json:"server_sign"`     // 服务标识
	SecretCipher string `gorm:"column:secret_cipher" json:"secret_cipher"` // 服务密钥密文 =  = AES(KEY, 服务密钥)
	DataCheck    string `gorm:"column:data_check" json:"data_check"`       // 数据校验 = HMACSHA256(KEY, secret_seed_cipher + account + is_enable)
	Common
}

func (s *ServerSecretModel) TableName() string {
	return "otp_server_secret"
}
