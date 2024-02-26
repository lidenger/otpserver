package model

// ServerIpWhitelistModel 接入服务IP白名单
type ServerIpWhitelistModel struct {
	ID         int64  `gorm:"column:id;primary_key" json:"id"`
	IP         string `gorm:"column:ip" json:"ip"`                   // IP白名单
	ServerSign string `gorm:"column:server_sign" json:"server_sign"` // 服务标识
	DataCheck  string `gorm:"column:data_check" json:"data_check"`   // 数据校验 = HMACSHA256(KEY, secret_seed_cipher + account + is_enable)
	Time
}

func (s *ServerIpWhitelistModel) TableName() string {
	return "otp_server_ip_whitelist"
}
